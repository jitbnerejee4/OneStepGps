package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"server/models"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/cloudinary/cloudinary-go/v2"

	"os"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/gomail.v2"

	"github.com/gin-contrib/cors"
)

var ctx = context.Background()
var jwtKey = []byte("one_step_gps")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type ApiResponse struct {
	ResultList []DeviceData `json:"result_list"`
}

type DeviceData struct {
	DeviceId          string      `json:"device_id"`
	DisplayName       string      `json:"display_name"`
	ActiveState       string      `json:"active_state"`
	LatestDevicePoint DevicePoint `json:"latest_device_point"`
	Address           string      `json:"address,omitempty"`
	IsFlipped         bool        `json:"is_flipped"`
	IsClicked         bool        `json:"is_clicked"`
	IsVisible         bool        `json:"is_visible"`
}

type DevicePoint struct {
	DTTracker string  `json:"dt_tracker"`
	Altitude  float64 `json:"altitude"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Speed     float64 `json:"speed"`
}

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func GetAddressFromLatLng(lat, lng float64) (string, error) {
	googleGeoAPIURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%v,%v&key=%s", lat, lng, os.Getenv("GOOGLE_GEOCODING_API"))
	resp, err := http.Get(fmt.Sprintf(googleGeoAPIURL))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if resultsArr, exists := result["results"].([]interface{}); exists && len(resultsArr) > 0 {
		if firstResult, ok := resultsArr[0].(map[string]interface{}); ok {
			if address, ok := firstResult["formatted_address"].(string); ok {
				return address, nil
			}
		}
	}

	return "", fmt.Errorf("failed to parse address")
}

func GetDeviceDataFromAPI() ([]DeviceData, error) {

	cacheKey := "deviceDataCache"

	// getting the cached data from Redis
	cachedData, err := rdb.Get(ctx, cacheKey).Result()

	// if the cache exists and error is nil
	if err == nil {
		var apiResponse ApiResponse
		err = json.Unmarshal([]byte(cachedData), &apiResponse)
		if err == nil {
			return apiResponse.ResultList, nil
		}
	}

	// fetch from the API when cache don't exist or there is error
	resp, err := http.Get(fmt.Sprintf("https://track.onestepgps.com/v3/api/public/device?latest_point=true&api-key=%s", os.Getenv("ONE_STEP_GPS")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	for i, device := range apiResponse.ResultList {
		address, err := GetAddressFromLatLng(device.LatestDevicePoint.Lat, device.LatestDevicePoint.Lng)
		if err != nil {
			continue
		}
		apiResponse.ResultList[i].IsFlipped = false
		apiResponse.ResultList[i].Address = address
		apiResponse.ResultList[i].IsClicked = false
		apiResponse.ResultList[i].IsVisible = true
	}

	// caching the new data in redis with a 1 hour expiration
	cachedDataBytes, err := json.Marshal(apiResponse)
	if err != nil {
		return nil, err
	}

	rdb.Set(ctx, cacheKey, string(cachedDataBytes), time.Hour)

	// returning the api response
	return apiResponse.ResultList, nil
}

func authenticateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is not Bearer token"})
		c.Abort()
		return
	}

	tokenString := parts[1]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Setting the email in the context
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token does not contain valid email"})
		c.Abort()
		return
	}

	c.Set("email", email)

	c.Next()
}
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &Claims{
		Email: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func registerUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("mydb").Collection("users")
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": newUser.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	newUser.Password = string(hashedPassword)

	devices, err := GetDeviceDataFromAPI()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching devices"})
		return
	}
	deviceIcons := make([]models.Device, len(devices))

	for i, device := range devices {
		deviceIcons[i] = models.Device{
			Name:    device.DisplayName,
			IconURL: "https://scai.kibu.ac.ke/wp-content/uploads/2021/10/NoProfile.png",
		}
	}

	newUser.Preferences = &models.Preference{
		ViewPreference:       "list",
		SortPreference:       "unsorted",
		VisibilityPreference: make(map[string]bool),
		DeviceIcons:          deviceIcons,
	}
	for _, device := range devices {
		newUser.Preferences.VisibilityPreference[device.DeviceId] = true
	}

	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}
	tokenString, err := GenerateJWT(newUser.Email)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"message": "Could not generate token"})
		return
	}
	c.JSON(200, gin.H{"token": tokenString, "email": newUser.Email, "userPreferences": newUser.Preferences})

}

func loginUser(c *gin.Context) {
	var user models.User
	var foundUser models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collection := client.Database("mydb").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	tokenString, err := GenerateJWT(user.Email)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"message": "Could not generate token"})
		return
	}
	c.JSON(200, gin.H{"token": tokenString, "email": foundUser.Email, "userPreferences": foundUser.Preferences})
}

func savePreference(c *gin.Context) {
	userEmail := c.Param("user")
	var preference models.Preference
	var foundUser models.User

	if err := c.BindJSON(&preference); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("mydb").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"email": userEmail}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user found with given email"})
		return
	}
	update := bson.M{
		"$set": bson.M{
			"preferences.viewPreference":       preference.ViewPreference,
			"preferences.sortPreference":       preference.SortPreference,
			"preferences.visibilityPreference": preference.VisibilityPreference,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"email": foundUser.Email}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"email": userEmail}).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user preferences"})
		return
	}

	c.JSON(200, gin.H{"message": "Preferences saved successfully", "userPreferences": updatedUser.Preferences})

}

func uploadImage(c *gin.Context) {
	userEmail := c.Param("user")
	var foundUser models.User

	collection := client.Database("mydb").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"email": userEmail}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user found with given email"})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get the file"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open the file"})
		return
	}
	defer file.Close()

	deviceName := c.PostForm("device")
	if deviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device information is required"})
		return
	}
	uploadResult, err := cld.Upload.Upload(c, file, uploader.UploadParams{Folder: "onestepgps"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i, device := range foundUser.Preferences.DeviceIcons {
		if device.Name == deviceName {
			foundUser.Preferences.DeviceIcons[i].IconURL = uploadResult.SecureURL
			break
		}
	}
	update := bson.M{
		"$set": bson.M{
			"preferences": foundUser.Preferences,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"email": foundUser.Email}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	var updatedUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"email": userEmail}).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user preferences"})
		return
	}
	c.JSON(200, gin.H{"message": "Icon changed successfully", "userPreferences": updatedUser.Preferences})

}

func generateToken() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Handle error
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func sendResetEmail(email string, link string) error {
	from := "jitbanerjee4@gmail.com"
	to := email
	host := "smtp-relay.brevo.com"
	port := 587

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Reset password link")
	msg.SetBody("text/html", "Please click on the following link to reset your password: <a href=\""+link+"\">Reset Password</a>")

	n := gomail.NewDialer(host, port, from, "xsmtpsib-71b62bfc06824d89dfa921e6a4ab678f145db9f112a9128427cc104e95cbd219-mWUOgPsAh0FbxqfV")
	if err := n.DialAndSend(msg); err != nil {
		log.Println("Error sending email:", err)
		return err
	}
	return nil
}
func forgetPassword(c *gin.Context) {
	userEmail := c.Param("user")
	var foundUser models.User

	token := generateToken()
	expiry := time.Now().Add(30 * time.Minute)
	update := bson.M{
		"$set": bson.M{
			"resetPasswordToken": token,
			"tokenExpiry":        expiry,
		},
	}

	collection := client.Database("mydb").Collection("users")
	err := collection.FindOneAndUpdate(context.TODO(), bson.M{"email": userEmail}, update).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user found with given email"})
		return
	}

	resetLink := "http://localhost:8080/reset-password?token=" + token
	if err := sendResetEmail(userEmail, resetLink); err != nil {
		// Handle error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reset email sent successfully."})

}

func resetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	collection := client.Database("mydb").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"resetPasswordToken": req.Token}).Decode(&user)
	if err != nil {
		c.JSON(401, gin.H{"error": "Token Expired"})
		return
	}
	if time.Now().After(user.TokenExpiry) {
		c.JSON(401, gin.H{"error": "Token expired"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}
	update := bson.M{
		"$set": bson.M{
			"password":           string(hashedPassword),
			"resetPasswordToken": "",
			"tokenExpiry":        nil,
		},
	}
	if _, err := collection.UpdateOne(ctx, bson.M{"email": user.Email}, update); err != nil {
		c.JSON(500, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

var client *mongo.Client
var cld *cloudinary.Cloudinary

func main() {
	r := gin.Default()
	godotenv.Load()
	api_key := os.Getenv("CLOUDINARY_API_KEY")
	api_secret := os.Getenv("CLOUDINARY_API_SECRET")
	cloud_name := os.Getenv("CLOUDINARY_CLOUD_NAME")
	cld, _ = cloudinary.NewFromParams(cloud_name, api_key, api_secret)
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "http://localhost:8080"
		// },
		MaxAge: 12 * time.Hour,
	}
	r.Use(cors.New(config))

	// initializing mongoDB
	var err error
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI("mongodb+srv://jitbaner:4vunF96SQJvY7DnG@cluster0.znvt1pl.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// Initializing reate limiter
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: rdb,
		Rate:        time.Second,
		Limit:       5,
	})

	// creating middleware
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	// main route to get device data. Added rate limiter to the route
	r.GET("/", mw, authenticateToken, func(c *gin.Context) {
		data, err := GetDeviceDataFromAPI()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, data)
	})
	r.POST("register", registerUser)
	r.POST("login", loginUser)
	r.POST("/preferences/:user", savePreference)
	r.POST("/upload/:user", uploadImage)
	r.GET("/forgot-password/:user", forgetPassword)
	r.POST("/reset-password", resetPassword)
	r.GET("/clearcache", func(c *gin.Context) {
		rdb.Del(ctx, "deviceDataCache")
		c.String(http.StatusOK, "Cache cleared!")
	})

	r.Run(":8000")
}

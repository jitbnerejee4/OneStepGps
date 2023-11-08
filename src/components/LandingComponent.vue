<template>
    <div class="video-container">
    <video autoplay loop muted playsinline>
      <source src="@/assets/network_-_45961 (720p).mp4" type="video/mp4">
      Your browser does not support the video tag.
    </video>
  </div>
  <div class="container">
    <div v-if="loading" class="loading-overlay">
      <img style="width: 200px; height: 200px; display: flex; justify-content: center;"  src="https://media0.giphy.com/media/QODU6spbkmhzg14hLx/giphy.gif?cid=ecf05e4742053fkspfuyavy4tabyqjg33op18hul4n9hzfxy&ep=v1_stickers_search&rid=giphy.gif&ct=s" alt="loading">
    </div>
    <div class="row">
        <div class="col col-md-12 col-sm-12 banner">
            <h1><span v-for="(char, index) in splitText" :key="index" :style="{ animationDelay: index * 0.1 + 's' }" :class="{ 'space': char === ' ' }"> {{ char === ' ' ? '\xa0' : char }}</span></h1>
        </div>
    </div>
    <div class="row">
        <div class="col col-md-12 col-sm-12">
            <div class="wrapper">
                <div class="typing-demo">
                    Secure Every Drive, Every Mile, Every Step.
                </div>
            </div>
        </div>
    </div>
    <div class="row" v-if="!isLoggedIn">
      <div class="col col-md-12 col-sm-12">
        <button class="btn btn-primary" type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasScrolling2" aria-controls="offcanvasScrolling2">Log In</button>
        <button class="btn btn-primary" type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasScrolling" aria-controls="offcanvasScrolling">Sign Up</button>
      </div>
    </div>
    <div class="row" v-if="isLoggedIn">
      <div class="col col-md-12 col-sm-12">
        <router-link class="nav-link" to="/home">
          <button class="btn btn-primary">HOME</button>
        </router-link>
      </div>
    </div>
    <div class="offcanvas offcanvas-end text-bg-success" data-bs-scroll="true" data-bs-backdrop="false" tabindex="-1" id="offcanvasScrolling" aria-labelledby="offcanvasScrollingLabel">
      <div class="offcanvas-header">
        <div style="display: flex; justify-content: center;">
          <!-- <h5 class="offcanvas-title" id="offcanvasScrollingLabel">Sign Up</h5> -->
        </div>
        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close" style="width: 65px;"></button>
      </div>
      <div class="offcanvas-body">
        <!-- <div v-if="loading" class="loading-overlay">
          <img style="widows: 500px; height: 500px; display: flex; justify-content: center;" src="https://media1.giphy.com/media/j6XZlxTBLE1taAzsGV/giphy.gif?cid=ecf05e47s0ga6mgywr84x6txtfegny2omv3a41f50n0x0jpq&ep=v1_stickers_search&rid=giphy.gif&ct=s" alt="loading">
        </div> -->
        <form @submit.prevent="handleRegister">
          <div class="row">
            <div class="col col-md-12 col-sm-12">
              <h2>Sign Up</h2>
            </div>
          </div>
          <div class="row">
            <div class="col col-md-12 col-sm-12" style="margin-top: 40px; text-align: start;">
              <label for="email" class="form-label">Email*</label>
                  <input type="email" class="form-control" id="email" name="email" placeholder="janeDoe1234@abc.xom" v-model="user.email"  required>
            </div>
          </div>
          <div class="row">
            <div class="col col-md-12 col-sm-12" style="margin-top: 40px; text-align: start;">
              <label for="password" class="form-label">Password*</label>
                <input :type="showSignupPassword ? 'text' : 'password'"  class="form-control" id="password" name="password" maxlength=20 v-model="user.password"  required>
                <span class="field-icon" style="color: black;" v-if="!showSignupPassword" @click="toggleSignupPasswordVisibility"><i class="fa-solid fa-eye-slash"></i></span>
                <span class="field-icon" style="color: black;" v-if="showSignupPassword" @click="toggleSignupPasswordVisibility"><i class="fa-solid fa-eye"></i></span>
                <p class="text-dark" v-if="!passwordValid">Password must be at least 6 characters long and include at least 1 number and 1 special character.</p>
            </div>
          </div>
          <div class="row">
            <div class="colcol-md-12 col-sm-12" style="margin-top: 40px;">
              <button type="submit" class="btn btn-danger" :disabled="!passwordValid">Submit</button>
            </div>
          </div>
        </form>
      </div>
    </div>
    <div class="offcanvas offcanvas-start text-bg-danger" data-bs-scroll="true" data-bs-backdrop="false" tabindex="-1" id="offcanvasScrolling2" aria-labelledby="offcanvasScrollingLabel">
      <div class="offcanvas-header">
        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close" style="width: 65px;"></button>
      </div>
      <div class="offcanvas-body">
        <form @submit.prevent="handleLogin">
          <div class="row">
            <div class="col col-md-12 col-sm-12">
              <h2>Log In</h2>
            </div>
          </div>
          <div class="row">
            <div class="col col-md-12 col-sm-12" style="margin-top: 40px; text-align: start;">
              <label for="loginemail" class="form-label">Email*</label>
                  <input type="email" class="form-control" id="loginemail" name="email" placeholder="janeDoe1234@abc.xom" v-model="loginUser.email"  required>
            </div>
          </div>
          <div class="row">
            <div class="col col-md-12 col-sm-12" style="margin-top: 40px; text-align: start;">
              <label for="loginpassword" class="form-label">Password*</label>
                <input :type="showPassword ? 'text' : 'password'" class="form-control" id="loginpassword" name="password" maxlength=20 v-model="loginUser.password"  required>
                <span class="field-icon" style="color: black;" v-if="!showPassword" @click="togglePasswordVisibility"><i class="fa-solid fa-eye-slash"></i></span>
                <span class="field-icon" style="color: black;" v-if="showPassword" @click="togglePasswordVisibility"><i class="fa-solid fa-eye"></i></span>
                <span>
                  <router-link to="/forgot-password">
                    <text style="text-decoration: underline; color: blue; cursor: pointer;">Forgot Password?</text>
                  </router-link>
                </span>
            </div>
          </div>
          <div class="row">
            <div class="colcol-md-12 col-sm-12" style="margin-top: 40px;">
              <button type="submit" class="btn btn-success">Submit</button>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import AuthService from '../services/AuthService.js'
import router from '@/router/router';
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';


export default {
  name: 'LandingComponent',
  data(){
    return{
      user: {
        email: '',
        password: '',
      },
      loginUser: {
        email: '',
        password: ''
      },
      loading: false,
      passwordValid: false,
      showPassword: false,
      showSignupPassword: false,
      isLoggedIn: false
    }
  },
  mounted(){
    const token = localStorage.getItem('token')
    if(token){
      this.isLoggedIn = true
    }else{
      this.isLoggedIn = false
    }
  },
  watch: {
    'user.password': function() {
      this.validatePassword();
    }
  },
  computed:{
    splitText(){
        return "One Step Gps".split('');
    }
  },
  methods:{
    validatePassword() {
      const regex = /^(?=.*[0-9])(?=.*[!@#$%^&*])[a-zA-Z0-9!@#$%^&*]{6,}$/;
      this.passwordValid = regex.test(this.user.password);
    },
    togglePasswordVisibility(){
      this.showPassword = !this.showPassword
    },
    toggleSignupPasswordVisibility(){
      this.showSignupPassword = !this.showSignupPassword
    },
    async handleRegister(){
      this.loading = true
      const response = await AuthService.register(this.user)
      if (response.status == 200) {
        this.user.email = ''
        this.user.password = ''
        setTimeout(() => {
          this.loading = false;
          this.$store.commit('setUserPreference', {preference: response.data.userPreferences})
          router.push('/home');
        }, 2000);
      }else{
        this.loading = false
        
      }
    },
    async handleLogin(){
      this.loading = true
      const response = await AuthService.login(this.loginUser)
      if (response.status == 200) {
        this.loginUser.email = ''
        this.loginUser.password = ''
        setTimeout(() => {
          this.loading = false;
          this.$store.commit('setUserPreference', {preference: response.data.userPreferences})
          router.push('/home');
        }, 2000); 
      }else{
        this.loading = false
        this.loginUser.password = ''
        toast.error("Invalid Username or Password!", {
          position: toast.POSITION.TOP_RIGHT,
        });
      }
    }
  }
}
</script>
<style>

.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(255, 255, 255, 0.7);
  z-index: 1000;
}
.field-icon{
    float: right;
    margin-left: -25px;
    margin-top: -30px;
    margin-right: 20px;
    position: relative;
    z-index: 2;
}

.video-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: -1; 
}

video {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  min-width: 100%;
  min-height: 100%;
  width: auto;
  height: auto;
  object-fit:cover;
}
.container{
    margin-top: 10%;
}
.container .banner{
    text-transform: uppercase;
    letter-spacing: 0.5em;
    font-size: 2vw; 
}
.container button{
    height: 50px;
    letter-spacing: 5px;
    text-transform: uppercase;
    padding: 10px;
    margin: 10px;
    /* width: 80%; */
}
.space {
  width: 0.25em; 
  display: inline-block;
}

h1 span {
  display: inline-block;
  opacity: 0;
  transform: translateY(1em);
  animation: fadeInUp 0.5s forwards;
}

@keyframes fadeInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.wrapper {
    display: grid;
    place-items: center;
}
  
.typing-demo {
    width: 43ch;
    animation: typing 2s steps(22), blink .5s step-end infinite alternate;
    white-space: nowrap;
    overflow: hidden;
    border-right: 3px solid;
    font-family: monospace;
    font-size: 2vw;
}
  
@keyframes typing {
    from {
      width: 0
    }
}
      
@keyframes blink {
    50% {
      border-color: transparent
    }
}

@media (max-width: 768px) {
    .container .banner {
        letter-spacing: 0.5em; 
    }
    .container button {
        width: 80%;
        margin: 10px auto;
    }
    .wrapper {
        padding: 20px;
    }
}
</style>
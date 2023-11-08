package models

import (
	"time"
)

type User struct {
	Email              string      `json:"email"`
	Password           string      `json:"password"`
	Preferences        *Preference `json:"preferences" bson:"preferences"`
	ResetPasswordToken string      `json:"resetPasswordToken,omitempty" bson:"resetPasswordToken,omitempty"`
	TokenExpiry        time.Time   `json:"tokenExpiry,omitempty" bson:"tokenExpiry,omitempty"`
}

type Preference struct {
	ViewPreference       string          `json:"viewPreference" bson:"viewPreference"`
	SortPreference       string          `json:"sortPreference" bson:"sortPreference"`
	VisibilityPreference map[string]bool `json:"visibilityPreference" bson:"visibilityPreference"`
	DeviceIcons          []Device        `json:"deviceIcons" bson:"deviceIcons"`
}

type Device struct {
	Name    string `json:"name"`
	IconURL string `json:"iconUrl"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

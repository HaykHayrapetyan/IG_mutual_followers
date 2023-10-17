package model

import (
	"diary_api/database"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"size:255;not null;unique" json:"username"`
	Password  string `gorm:"size:255;not null;" json:"-"`
	Entries   []Entry
	Followers Followers `gorm:"type:json"` //TODO find a way to change this to []Follower
	InstaId   string
}

func (user *User) Save() (*User, error) {
	err := database.Database.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) BeforeSave(*gorm.DB) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByUsername(username string) (User, error) {
	var user User
	err := database.Database.Where("username=?", username).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func FindUserById(id uint) (User, error) {
	var user User
	err := database.Database.Preload("Entries").Where("ID=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (user *User) SaveFollowers(followers []Follower) error {

	// Update the Followers field with the provided followers
	user.Followers = followers
	// Save the user with the updated followers
	if err := database.Database.Save(user).Error; err != nil {
		return err
	}
	return nil
}

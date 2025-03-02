package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"jwt-psql/utils/token"
	"log"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

func (u *User) SaveUser() (*User, error) {
	var err error
	err = db.Create(&u).Error

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func GetUserByID(id uint) (*User, error) {
	var u User
	err := db.First(&u, id).Error
	if err != nil {
		return &User{}, err
	}

	u.Password = ""
	return &u, nil
}

func LoginAndGenerateToken(username, password string) (string, error) {
	var err error

	u := User{}

	err = db.Model(u).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", err
	}

	generateToken, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return generateToken, nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func UserExistsByUsername(username string) (bool, error) {
	var exists bool
	var model User

	err := db.Model(model).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&exists).
		Error

	if err != nil {
		log.Println(err)
		return false, err
	}

	return exists, err
}

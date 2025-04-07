package flags

import (
	"fast-gin/global"
	"fast-gin/models"
	"fast-gin/utils/pwd"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

type User struct {
}

// Create Unnamed receiver acts like static method
func (User) Create() {
	var user models.UserModel

	// Role
	fmt.Println("Please select a role for user (1 (admin) 2 (normal)): ")
	_, err := fmt.Scanln(&user.RoleID)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}
	if user.RoleID != 1 && user.RoleID != 2 {
		fmt.Println("Role err:", err)
		return
	}

	// Username
	for {
		fmt.Println("Please input username: ")
		_, err = fmt.Scanln(&user.Username)
		if err != nil {
			fmt.Println("Input error:", err)
			return
		}
		var u models.UserModel
		err = global.DB.Take(&u, "username = ?", user.Username).Error
		if err == nil {
			fmt.Println("User already exists")
			continue
		}
		break
	}

	// Password
	fmt.Println("Please input password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to read password:", err)
		return
	}
	fmt.Println("Please input password again: ")
	rePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to read password:", err)
		return
	}
	if string(password) != string(rePassword) {
		fmt.Println("Password mismatched")
		return
	}

	// Persist
	encryptedPassword, err := pwd.Encrypt(string(password))
	if err != nil {
		fmt.Println("Failed to encrypt password:", err)
	}
	err = global.DB.Create(&models.UserModel{
		Username: user.Username,
		Password: encryptedPassword,
		RoleID:   user.RoleID,
	}).Error
	if err != nil {
		logrus.Errorf("Failed to create user: %s", err)
		return
	}
	logrus.Infof("Create user [%s] successfully", user.Username)
}

// List Unnamed receiver acts like static method
func (User) List() {
	var userList []models.UserModel
	global.DB.Order("created_at desc").Limit(10).Find(&userList)
	for _, model := range userList {
		fmt.Printf("UserID: %d  Username: %s Nickname: %s Role: %d CreatedAt: %s\n",
			model.ID,
			model.Username,
			model.Nickname,
			model.RoleID,
			model.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}
}

// Remove Unnamed receiver acts like static method
func (User) Remove() {
	var username string

	// Username
	for {
		fmt.Println("Please input username of user to be deleted: ")
		_, err := fmt.Scanln(&username)
		if err != nil {
			fmt.Println("Input error:", err)
			return
		}
		var u models.UserModel
		err = global.DB.Take(&u, "username = ?", username).Error
		if err != nil {
			fmt.Println("User does not exist")
			continue
		}
		break
	}

	err := global.DB.
		Where("username = ?", username).
		Delete(&models.UserModel{}).Error
	if err != nil {
		logrus.Errorf("Failed to delete user: %s", err)
		return
	}
	logrus.Infof("Delete user [%s] successfully", username)
}

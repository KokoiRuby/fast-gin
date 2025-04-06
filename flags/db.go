package flags

import (
	"fast-gin/global"
	"fast-gin/models"

	"github.com/sirupsen/logrus"
)

func MigrateDB() {
	err := global.DB.AutoMigrate(&models.UserModel{})
	if err != nil {
		logrus.Errorf("Failed to migrate database: %s", err)
		return
	}
	logrus.Infof("Migrate database successfully")
}

package image

import (
	"fast-gin/global"
	"fast-gin/utils/md5"
	"fast-gin/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var whiteList = map[string]struct{}{
	".jpg": {},
	".png": {},
	".gif": {},
}

func (API) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.FailWithMsg(c, "Please select an image")
		return
	}

	// Size
	if fileHeader.Size > global.Config.Upload.Size*1024*1024 {
		response.FailWithMsg(c, "File size too big (>2MB)")
		return
	}

	// Suffix whitelist
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if _, ok := whiteList[ext]; !ok {
		response.FailWithMsg(c, "Image extension not supported (.jpg, .gif, .png)")
		return
	}

	// De-duplicate
	dir := path.Join("uploads", global.Config.Upload.Dir, fileHeader.Filename)

	_, err = os.Stat(dir)
	if !os.IsNotExist(err) {
		// vs.
		toUpload, err := fileHeader.Open()
		if err != nil {
			response.FailWithMsg(c, err.Error())
		}
		toUploadHash := md5.GetMD5(toUpload)
		logrus.Debugf("tUpload hash: %s", toUploadHash)

		existed, err := os.Open(dir)
		if err != nil {
			response.FailWithMsg(c, err.Error())
		}
		existedHash := md5.GetMD5(existed)
		logrus.Debugf("existed hash: %s", existedHash)

		if existedHash == toUploadHash {
			response.OK(c, gin.H{}, "Upload successfully")
			return
		}

		// TODO: random string
		newFilename := fmt.Sprintf("%s_%s.%s", strings.TrimSuffix(fileHeader.Filename, ext), "random_string", ext)
		dir = path.Join("uploads", global.Config.Upload.Dir, newFilename)
	}

	err = c.SaveUploadedFile(fileHeader, dir)
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}

	response.OK(c, gin.H{}, "Upload successfully")
}

package md5

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"io"
)

func GetMD5(file io.Reader) string {
	m := md5.New()
	_, err := io.Copy(m, file)
	if err != nil {
		logrus.Errorf("Failed to copy file to hash: %v", err)
		return ""
	}
	sum := m.Sum(nil)
	return hex.EncodeToString(sum)
}

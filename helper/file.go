package helper

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func SaveFile(file *multipart.FileHeader, fileName string, dirPath string) error {
	myFile, err := file.Open()
	if err != nil {
		logrus.Errorln("open file err -->", err)
		return err
	}
	defer myFile.Close()

	dst, err := os.Create(filepath.Join(dirPath, filepath.Base(fileName)))
	if err != nil {
		logrus.Errorln("create file err -->", err)
		return err

	}
	defer dst.Close()
	if _, err = io.Copy(dst, myFile); err != nil {
		logrus.Errorln("copy file err -->", err)
		return err
	}

	return nil
}
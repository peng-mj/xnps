package sysTool

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os"
)

func DirExisted(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func FileExisted(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CreateFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func CreateFileWithContent(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	return nil
}
func CreateFolder(dir string) {
	if !DirExisted(dir) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
	return
}

func SaveImageFromBase64(filename, path, base64Data string) (length int, err error) {
	var tempData []byte
	if tempData, err = base64.StdEncoding.DecodeString(base64Data); err != nil {
		return
	}
	var outputFile *os.File
	if outputFile, err = os.Create(path + filename); err != nil {
		return
	} else {
		defer outputFile.Close()
	}
	length, err = outputFile.Write(tempData)
	return
}
func SaveBase64Img(filePath, base64Data string) (length int, err error) {
	var tempData []byte
	if tempData, err = base64.StdEncoding.DecodeString(base64Data); err != nil {
		return
	}
	var outputFile *os.File
	if outputFile, err = os.Create(filePath); err != nil {
		return
	} else {
		defer outputFile.Close()
	}
	length, err = outputFile.Write(tempData)
	return
}
func EncodeBase64File(imagePath string) (string, error) {
	if data, err := os.ReadFile(imagePath); err == nil {
		encodedData := base64.StdEncoding.EncodeToString(data)
		return encodedData, nil
	} else {
		return "", err
	}
}

// CheckPortOccupied to check net port
func CheckPortOccupied(port int, protocol string) error {
	switch protocol {
	case "tcp":
	case "udp":
	default:
		return errors.New("not support " + protocol)
	}
	_, err := net.Listen(protocol, fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return nil
}

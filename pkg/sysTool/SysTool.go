package sysTool

import (
	"encoding/base64"
	"os"
	"runtime"
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

func GetBasePath() string {
	base, err := os.Getwd()
	if err != nil || base == "" {
		base = "."
	}
	if runtime.GOOS == "windows" {
		base = base + "\\.tunpx"
	} else {
		base = base + "/.tunpx"
	}
	CreateFolder(base)
	return base

}

func SaveBase64File(filePath, base64Data string) (length int, err error) {
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

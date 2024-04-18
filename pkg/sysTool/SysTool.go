package sysTool

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"xnps/lib/crypt"
)

// 判断文件夹是否存在
func DirExisted(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// 判断文件是否存在
func FileExisted(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 创建新文件
func CreateFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// 创建新文件
func CreateAndWriteFile(path, content string) error {

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(content))
	defer file.Close()
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
func EncodeImgFileToBase64(imagePath string) (string, error) {
	if data, err := os.ReadFile(imagePath); err == nil {
		encodedData := base64.StdEncoding.EncodeToString(data)
		return encodedData, nil
	} else {
		return "", err
	}
}
func GetMd5FromFile(imagePath string) string {
	if !FileExisted(imagePath) {
		return "@"
	}
	if data, err := os.ReadFile(imagePath); err == nil {
		if len(data) < 3000 {
			return "#"
		}
		return crypt.Md5(data)
	} else {
		return "#"
	}
}

func BackupFile(src, dst string, replace bool) error {

	if !FileExisted(dst) && replace {
		return fmt.Errorf("目标文件已存在: %s", src)
	}
	// 检查源文件是否存在
	if !FileExisted(src) {
		return fmt.Errorf("源文件不存在: %s", src)
	}
	CreateFolder(dst)

	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %v", err)
	}
	defer srcFile.Close()

	// 创建目标文件名，添加当前日期作为前缀
	filename := time.Now().Format("2006-01-02.") + filepath.Base(src)
	dstPath := filepath.Join(dst, filename)

	if FileExisted(dstPath) {
		os.Remove(dstPath)
	}
	// 创建目标文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dstFile.Close()
	// 复制源文件到目标文件
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件失败: %v", err)
	}

	return nil
}

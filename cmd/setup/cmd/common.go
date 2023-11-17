package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/lightec-xyz/daemon/common"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileMd5 struct {
	File string `json:"file"`
	Md5  string `json:"md5"`
}

func ComputeMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)[:16]
	hashString := hex.EncodeToString(hashInBytes)
	return hashString, nil
}

func isDesiredFileType(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".ccs" || ext == ".pk" || ext == ".vk" || ext == ".sol"
}

func GetFilesMd5(dirPath string) ([]FileMd5, error) {
	var fileMd5s []FileMd5
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !isDesiredFileType(info.Name()) {
			return nil
		}
		md5Value, err := ComputeMd5(path)
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		fileMd5s = append(fileMd5s, FileMd5{File: relativePath, Md5: md5Value})

		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileMd5s, nil
}

func SaveMd5sToJson(fileMd5s []FileMd5, outputPath string) error {
	data, err := json.MarshalIndent(fileMd5s, "", "  ")
	if err != nil {
		return err
	}
	file, err := common.OverwriteFile(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

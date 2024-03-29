package common

import (
	"encoding/json"
	"github.com/google/uuid"
)

func Uuid() (string, error) {
	newV7, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return newV7.String(), nil
}

func objToJson(obj interface{}) string {
	ojbBytes, err := json.Marshal(obj)
	if err != nil {
		return "error obj to josn"
	}
	return string(ojbBytes)

}

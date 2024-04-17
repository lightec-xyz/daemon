package common

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"reflect"
)

func Uuid() (string, error) {
	newV7, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return newV7.String(), nil
}
func MustUUID() string {
	newV7, err := uuid.NewV7()
	if err != nil {
		panic("gen uuid error,should never happen")
	}
	return newV7.String()
}

func objToJson(obj interface{}) string {
	ojbBytes, err := json.Marshal(obj)
	if err != nil {
		return "error obj to josn"
	}
	return string(ojbBytes)
}

func ParseObj(src, dst interface{}) error {
	if reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}
	srcBytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(srcBytes, dst)
	if err != nil {
		return err
	}
	return nil
}

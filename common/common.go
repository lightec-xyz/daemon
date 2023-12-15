package common

import "github.com/google/uuid"

func Uuid() (string, error) {
	newV7, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return newV7.String(), nil
}

package utils

import "strings"

func StandardError(err error) string {
	errorMessage := strings.Split(err.Error(), ": ")
	if len(errorMessage) > 1 {
		return errorMessage[len(errorMessage)-1]
	}
	return err.Error()
}

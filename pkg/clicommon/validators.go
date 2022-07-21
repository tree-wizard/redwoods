package clicommon

import (
	"errors"
	"strconv"
	"strings"
)

//ValidateEmptyInput
func ValidateEmptyInput(input string) error {
	if len(strings.TrimSpace(input)) < 1 {
		return errors.New("sorry! this input must not be empty")
	}
	return nil
}

//ValidateIntegerNumberInput
func ValidateIntegerNumberInput(input string) error {
	_, err := strconv.ParseInt(input, 0, 64)
	if err != nil {
		return errors.New("invalid number")
	}
	return nil
}

package dto

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type PublicKeyEncrypDto struct {
	Data string `json:"data" binding:"required,gte=1"`
}

func CommonValidateName(name string, minLen, maxLen int) error {
	if name == "" {
		return errors.New("name is empty")
	}
	if strings.HasPrefix(name, " ") || strings.HasSuffix(name, " ") {
		return errors.New("name cannot head and tail with space")
	}
	nameLen := utf8.RuneCountInString(name)
	if nameLen < minLen {
		return errors.New("name is too short")
	}
	if nameLen > maxLen {
		return errors.New("name is too long")
	}
	return nil
}

func CommonValidateNameList(names []string, minLen, maxLen int) error {
	seen := make(map[string]struct{})
	for _, name := range names {
		if _, ok := seen[name]; ok {
			return errors.New("name contains duplicate elements")
		}
		seen[name] = struct{}{}
		if err := CommonValidateName(name, minLen, maxLen); err != nil {
			return err
		}
	}
	return nil
}

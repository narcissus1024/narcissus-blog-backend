package dto

import (
	"errors"
	"strings"
)

const (
	TAG_MIN_LEN = 2
	TAG_MAX_LEN = 20
)

type TagDto struct {
	NameList []string `json:"nameList" binding:"required"`
}

// ValidateAndDefault 校验并默认值
func (r *TagDto) ValidateAndDefault() error {
	if len(r.NameList) == 0 {
		return errors.New("nameList is empty")
	}
	for i := range r.NameList {
		if err := CommonValidateName(r.NameList[i], TAG_MIN_LEN, TAG_MAX_LEN); err != nil {
			return err
		}
	}
	return nil
}

type TagQueryDto struct {
	ID   *int64  `json:"id" form:"id" binding:"omitempty,gte=1"`
	Name *string `json:"name" form:"name" binding:"omitempty,no_spacing,gte=2,lt=20"`
}

func (r *TagQueryDto) ValidateAndDefault() error {
	if r.ID == nil && r.Name == nil {
		return errors.New("id and name are empty")
	}
	return nil
}

type TagListDto struct {
	Pageinate
	NameList       string   `json:"nameList" form:"nameList"`
	NameListFormat []string `json:"-" form:"-"`
}

func (r *TagListDto) ValidateAndDefault() error {
	if len(r.NameList) > 0 {
		r.NameListFormat = strings.Split(r.NameList, ",")
	}
	for i := range r.NameListFormat {
		if err := CommonValidateName(r.NameListFormat[i], TAG_MIN_LEN, TAG_MAX_LEN); err != nil {
			return err
		}
	}
	return nil
}

type TagUpdateDto struct {
	ID      int64  `json:"id" binding:"required,gte=1"`
	NewName string `json:"newName" binding:"required,no_spacing,gte=2,lte=20"`
}

func (r *TagUpdateDto) ValidateAndDefault() error {
	return nil
}

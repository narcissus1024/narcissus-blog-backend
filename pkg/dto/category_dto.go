package dto

import (
	"errors"
	"strings"
)

const (
	CATEGORY_MIN_LEN = 2
	CATEGORY_MAX_LEN = 20
)

type CategoryDto struct {
	NameList []string `json:"nameList" binding:"required"`
}

func (r *CategoryDto) ValidateAndDefault() error {
	if len(r.NameList) == 0 {
		return errors.New("nameList is empty")
	}
	for i := range r.NameList {
		if err := CommonValidateName(r.NameList[i], CATEGORY_MIN_LEN, CATEGORY_MAX_LEN); err != nil {
			return err
		}
	}
	return nil
}

type CategoryQueryDto struct {
	ID   *int64  `json:"id" form:"id" binding:"omitempty,gte=1"`
	Name *string `json:"name" form:"name" binding:"omitempty,no_spacing,gte=2,lt=20"`
}

func (r *CategoryQueryDto) ValidateAndDefault() error {
	if r.ID == nil && r.Name == nil {
		return errors.New("id and name are empty")
	}
	return nil
}

type CategoryListDto struct {
	Pageinate
	NameList       string   `json:"nameList" form:"nameList"`
	NameListFormat []string `json:"-" form:"-"`
}

func (r *CategoryListDto) ValidateAndDefault() error {
	if len(r.NameList) > 0 {
		r.NameListFormat = strings.Split(r.NameList, ",")
	}
	for i := range r.NameListFormat {
		if err := CommonValidateName(r.NameListFormat[i], CATEGORY_MIN_LEN, CATEGORY_MAX_LEN); err != nil {
			return err
		}
	}
	return nil
}

type CategoryUpdateDto struct {
	ID      int64  `json:"id" binding:"required,gte=1"`
	NewName string `json:"newName" binding:"required,no_spacing,gte=2,lte=20"`
}

func (r *CategoryUpdateDto) ValidateAndDefault() error {
	return nil
}

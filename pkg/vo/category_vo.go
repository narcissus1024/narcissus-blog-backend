package vo

import "github.com/narcissus1949/narcissus-blog/pkg/dto"

type CategoryVo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	UpdatedTime string `json:"updatedTime"`
	CreatedTime string `json:"createdTime"`
}

type CategoryListVo struct {
	CategoryList []CategoryVo  `json:"categoryList"`
	Pageinate    dto.Pageinate `json:"pageinate"`
}

package vo

import "github.com/narcissus1949/narcissus-blog/pkg/dto"

type TagVo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	CreatedTime string `json:"createdTime"`
	UpdatedTime string `json:"updatedTime"`
}

type TagListVo struct {
	TagList   []TagVo       `json:"tagList"`
	Pageinate dto.Pageinate `json:"pageinate"`
}

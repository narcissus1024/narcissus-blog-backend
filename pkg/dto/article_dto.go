package dto

import (
	"errors"

	"github.com/mcuadros/go-defaults"
	"github.com/narcissus1949/narcissus-blog/internal/utils"
)

// 创建文章参数
type ArticleDto struct {
	ID                  *int64   `json:"id"`
	Title               string   `json:"title" binding:"required,no_lt_spacing,gte=1"`         // 文章标题
	Summary             string   `json:"summary"`                                              // 摘要
	Content             string   `json:"content" binding:"required,no_lt_spacing,gte=1"`       // 文章内容
	Type                uint8    `json:"type" binding:"gte=0,oneof=0 1 2"`                     // 文章类型，0表示博客文章，1表示随笔，2表示关于
	Category            string   `json:"category" binding:"no_spacing"`                        // 文章分类，每篇文章最多1个分类，可以为空
	Tags                []string `json:"tags"`                                                 // 文章标签列表，每篇文章最多有5个标签
	Author              string   `json:"author" binding:"required,no_lt_spacing,gte=2,lte=50"` // 作者姓名或标识
	AllowComment        bool     `json:"allow_comment"`                                        // 是否允许评论，0表示不允许评论，1表示允许评论
	Weight              int      `json:"weight"`                                               // 文章权重，默认初始值为0
	IsSticky            bool     `json:"is_sticky"`                                            // 是否置顶。0表示不置顶，1表示置顶。默认初始值为 0
	IsOriginal          bool     `json:"is_original"`                                          // 原创/转载标识。0表示非原创，1表示原创。默认初始值为1，表示原创
	OriginalArticleLink string   `json:"original_article_link"`                                // 转载文章的原始文章链接，可为空
	Status              uint8    `json:"status" binding:"oneof=0 1"`                           // 状态，0表示offline，1表示online
}

func (req *ArticleDto) VlidateAndDefault() error {
	if req.Category != "" {
		if err := CommonValidateName(req.Category, CATEGORY_MIN_LEN, CATEGORY_MAX_LEN); err != nil {
			return err
		}
	}
	for i := range req.Tags {
		if err := CommonValidateName(req.Tags[i], TAG_MIN_LEN, TAG_MAX_LEN); err != nil {
			return err
		}
	}
	return nil
}

// 查询文章列表参数
type ArticleListDto struct {
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Type       []int    `json:"type"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
	IsOriginal *bool    `json:"is_original"`
	Status     *bool    `json:"status"`
	StartTime  int64    `json:"start_time"`
	EndTime    int64    `json:"end_time"`
	Pageinate
}

func (req *ArticleListDto) VlidateAndSetDefault() error {
	defaults.SetDefaults(req)
	for i := range req.Type {
		if req.Type[i] != utils.ARTICLE_TYPE_ABOUT &&
			req.Type[i] != utils.ARTICLE_TYPE_ESSAY &&
			req.Type[i] != utils.ARTICLE_TYPE_POST {
			return errors.New("article type invalide")
		}
	}
	for i := range req.Tags {
		if err := CommonValidateName(req.Tags[i], TAG_MIN_LEN, TAG_MAX_LEN); err != nil {
			return err
		}
	}
	// 文章类型默认查询博文和随笔
	if len(req.Type) == 0 {
		req.Type = append(req.Type, utils.ARTICLE_TYPE_POST, utils.ARTICLE_TYPE_ESSAY)
	}
	return nil
}

type ArticleDeleteDto struct {
	IDs []int64 `json:"ids" binding:"required"`
}

func (req *ArticleDeleteDto) VlidateAndDefault() error {
	if len(req.IDs) <= 0 {
		return errors.New("article id is empty")
	}
	for i := range req.IDs {
		if req.IDs[i] <= 0 {
			return errors.New("article id is invalid")
		}
	}
	return nil
}

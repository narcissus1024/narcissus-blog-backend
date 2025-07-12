package vo

import "github.com/narcissus1949/narcissus-blog/pkg/dto"

type ArticleMeta struct {
	ID                  int64  `json:"id"`                  // 文章ID
	Title               string `json:"title"`               // 文章标题
	Summary             string `json:"summary"`             // 摘要
	CategoryID          *int   `json:"categoryID"`          // 文章分类ID，每篇文章最多1个分类，可以为空
	Type                uint8  `json:"type"`                // 文章类型，0表示原创，1表示转载
	TagsID              string `json:"tagsID"`              // 文章标签ID列表，用,隔开，每篇文章最多有5个标签
	Author              string `json:"author"`              // 作者姓名或标识
	AllowComment        bool   `json:"allowComment"`        // 是否允许评论，0表示不允许评论，1表示允许评论
	Weight              int    `json:"weight"`              // 文章权重，默认初始值为0
	IsSticky            bool   `json:"isSticky"`            // 是否置顶。0表示不置顶，1表示置顶。默认初始值为 0
	IsOriginal          bool   `json:"isOriginal"`          // 原创/转载标识。0表示非原创，1表示原创。默认初始值为1，表示原创
	OriginalArticleLink string `json:"originalArticleLink"` // 转载文章的原始文章链接，可为空
	Status              uint8  `json:"status"`              // 状态，0表示offline，1表示online
	CreatedTime         int64  `json:"createdTime"`         // 创建时间
	UpdatedTime         int64  `json:"updatedTime"`         // 更新时间
}

type ArticleDetailVo struct {
	ArticleMeta
	CategoryName string   `json:"categoryName"`
	TagNameList  []string `json:"tagNameList"`
	Content      string   `json:"content"`
}

// 查询文章列表响应内容
type ArticleListVo struct {
	ArticleList []ArticleDetailVo `json:"articleList"`
	Pageinate   dto.Pageinate     `json:"pageinate"`
}

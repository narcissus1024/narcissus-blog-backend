package service

import (
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/cmd/blog/app/config"
	"github.com/narcissus1949/narcissus-blog/internal/database/cache"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	cerr "github.com/narcissus1949/narcissus-blog/internal/error"
	"github.com/narcissus1949/narcissus-blog/internal/logger"
	"github.com/narcissus1949/narcissus-blog/internal/model"
	"github.com/narcissus1949/narcissus-blog/internal/utils"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"github.com/narcissus1949/narcissus-blog/pkg/server/dao"
	"github.com/narcissus1949/narcissus-blog/pkg/vo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ArticleService = new(articleService)

type articleService struct {
}

func (s *articleService) ListArticleAdmin(ctx *gin.Context, articleListRequest dto.ArticleListDto) (*vo.ArticleListVo, error) {
	l := logger.FromContext(ctx.Request.Context())
	var articleListResponse vo.ArticleListVo
	var articleList []model.ArticleDetail
	var totalArticle int64
	txErr := mysql.RunDBTransaction(ctx, func() error {
		// 获取文章列表
		var listArticleErr error
		articleList, listArticleErr = dao.ArticleDao.ListArticle(ctx, articleListRequest)
		if listArticleErr != nil {
			l.Error("Failed to list article", zap.Error(listArticleErr))
			return listArticleErr
		}

		// 获取文章总数
		var countArticleErr error
		totalArticle, countArticleErr = dao.ArticleDao.CountArticle(ctx, articleListRequest)
		if countArticleErr != nil {
			l.Error("Failed to count article", zap.Error(countArticleErr))
			return countArticleErr
		}
		return nil
	})
	if txErr != nil {
		l.Error("Failed to run db transaction", zap.Error(txErr))
		return nil, txErr
	}

	// 响应参数封装
	list := make([]vo.ArticleDetailVo, 0, len(articleList))
	for i := range articleList {
		tagList := []string{}
		if len(articleList[i].TagNameList) > 0 {
			tagList = strings.Split(articleList[i].TagNameList, ",")
			sort.Strings(tagList)
		}
		// 获取文章浏览量缓存
		viewsCache := 0
		if intCmd := cache.Client.SCard(ctx, utils.GetArticlePageViewKey(articleList[i].ID)); intCmd.Err() != nil {
			l.Error("Failed to query article view from cache", zap.Error(intCmd.Err()), zap.Int64("article id", articleList[i].ID))
		} else {
			viewsCache = int(intCmd.Val())
		}
		list = append(list, vo.ArticleDetailVo{
			ArticleMeta: vo.ArticleMeta{
				ID:                  articleList[i].ID,
				Title:               articleList[i].Title,
				Summary:             articleList[i].Summary,
				CategoryID:          articleList[i].CategoryID,
				Type:                articleList[i].Type,
				Author:              articleList[i].Author,
				AllowComment:        articleList[i].AllowComment,
				Views:               articleList[i].Views + viewsCache,
				Weight:              articleList[i].Weight,
				IsSticky:            articleList[i].IsSticky,
				IsOriginal:          articleList[i].IsOriginal,
				OriginalArticleLink: articleList[i].OriginalArticleLink,
				Status:              articleList[i].Status,
				CreatedTime:         articleList[i].CreatedTime.UnixMilli(),
				UpdatedTime:         articleList[i].UpdatedTime.UnixMilli(),
			},
			CategoryName: articleList[i].CategoryName,
			TagNameList:  tagList,
		})
	}
	articleListResponse.ArticleList = list

	// set pageinate
	pageCount := totalArticle / int64(articleListRequest.Pageinate.PageSize)
	if totalArticle%int64(articleListRequest.Pageinate.PageSize) != 0 {
		pageCount += 1
	}
	articleListResponse.Pageinate = dto.Pageinate{
		PageSize:  articleListRequest.Pageinate.PageSize,
		PageNum:   articleListRequest.Pageinate.PageNum,
		Total:     totalArticle,
		PageCount: int(pageCount),
	}

	return &articleListResponse, nil
}

func (s *articleService) CreateArticle(c *gin.Context, articleDto dto.ArticleDto) error {
	l := logger.FromContext(c.Request.Context())

	now := time.Now()
	articleModel := &model.Article{
		Title:               articleDto.Title,
		Type:                articleDto.Type,
		Summary:             articleDto.Summary,
		Author:              articleDto.Author,
		AllowComment:        articleDto.AllowComment,
		Views:               0,
		Weight:              articleDto.Weight,
		IsSticky:            articleDto.IsSticky,
		IsOriginal:          articleDto.IsOriginal,
		OriginalArticleLink: articleDto.OriginalArticleLink,
		Status:              articleDto.Status,
		CreatedTime:         now,
		UpdatedTime:         now,
	}

	var tagIdList []int64
	// 如果文章类型为博文，则获取对应的分类、标签信息
	if articleDto.Type == utils.ARTICLE_TYPE_POST {
		if len(strings.TrimSpace(articleDto.Category)) > 0 {
			categoryId, err := CategoryService.GetCategoryIDByName(c, articleDto.Category)
			if err != nil {
				l.Error("Failed to get category id", zap.Error(err), zap.String("category", articleDto.Category))
				return err
			}
			if categoryId < 0 {
				l.Error("Category not exist, category name: ", zap.String("category", articleDto.Category))
				return cerr.New(cerr.ERROR_ARTICLE_CATEGORY_NOT_EXIST)
			}
			articleModel.CategoryID = &categoryId
		}

		if len(articleDto.Tags) > 0 {
			var listTagErr error
			tagIdList, listTagErr = TagService.ListTagIdByNameArr(c, dto.TagDto{NameList: articleDto.Tags})
			if listTagErr != nil {
				l.Error("Failed to list tag id by name arr", zap.Error(listTagErr), zap.Strings("tags", articleDto.Tags))
				return listTagErr
			}
			if len(tagIdList) == 0 {
				l.Error("Tag not exist, tags: ", zap.Strings("tags", articleDto.Tags))
				return cerr.New(cerr.ERROR_ARTICLE_TAG_NOT_EXIST)
			}
		}
	}

	txErr := mysql.RunDBTransaction(c, func() error {
		// 插入文章
		if err := dao.ArticleDao.InsertArticle(c, articleModel); err != nil {
			l.Error("Failed to insert article", zap.Error(err))
			return err
		}

		// 插入文章内容
		if err := dao.ArticleContentDao.InsertContent(c, &model.ArticleContent{
			ArticleID: articleModel.ID,
			Content:   articleDto.Content,
		}); err != nil {
			l.Error("Failed to insert article content", zap.Error(err))
			return err
		}

		// 插入文章标签关联
		if len(tagIdList) > 0 && articleModel.Type == utils.ARTICLE_TYPE_POST {
			var relations []*model.ArticleTagRelation
			for _, tagID := range tagIdList {
				relation := &model.ArticleTagRelation{
					ArticleID: articleModel.ID,
					TagID:     tagID,
				}
				relations = append(relations, relation)
			}
			if err := dao.ArticleTagRelationDao.InsertArticleTagRelations(c, relations); err != nil {
				l.Error("Failed to insert article tag relation", zap.Error(err))
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		l.Error("Failed to execute transaction", zap.Error(txErr))
		return txErr
	}

	return nil
}

func (s *articleService) UpdateArticle(c *gin.Context, articleDto dto.ArticleDto) error {
	l := logger.FromContext(c.Request.Context())
	if articleDto.ID == nil || *articleDto.ID <= 0 {
		return cerr.NewParamError("文章id无效")
	}

	txErr := mysql.RunDBTransaction(c, func() error {
		// 0.检查文章存在
		articleDetail, err := s.GetArticleDetail(c, *articleDto.ID)
		if err != nil {
			l.Error("Failed to get article detail", zap.Error(err), zap.Int64("articleID", *articleDto.ID))
			return err
		}
		if articleDetail == nil {
			return cerr.New(cerr.ERROR_ARTICLE_NOT_EXIST)
		}
		if articleDetail.Type != articleDto.Type {
			return cerr.NewParamError("不支持更改文章类型")
		}
		// 1.更新文章元数据
		now := time.Now()
		articleModel := &model.Article{
			ID:      *articleDto.ID,
			Title:   articleDto.Title,
			Summary: articleDto.Summary,
			Type:    articleDetail.Type, // 不支持更改文章类型
			// CategoryID:          articleDto.CategoryID,
			Author:              articleDto.Author,
			AllowComment:        articleDto.AllowComment,
			Weight:              articleDto.Weight,
			IsSticky:            articleDto.IsSticky,
			IsOriginal:          articleDto.IsOriginal,
			OriginalArticleLink: articleDto.OriginalArticleLink,
			Status:              articleDto.Status,
			CreatedTime:         time.Unix(articleDetail.CreatedTime, 0),
			UpdatedTime:         now,
		}
		// 处理文章分类
		if articleModel.Type != utils.ARTICLE_TYPE_ABOUT {
			if articleDetail.CategoryName == articleDto.Category {
				articleModel.CategoryID = articleDetail.CategoryID
			} else if len(strings.TrimSpace(articleDto.Category)) > 0 {
				categoryId, err := CategoryService.GetCategoryIDByName(c, articleDto.Category)
				if err != nil {
					l.Error("Failed to get category id", zap.Error(err), zap.String("category", articleDto.Category))
					return err
				}
				articleModel.CategoryID = &categoryId
			}
		}

		// 更新文章
		if _, updateArticleErr := dao.ArticleDao.UpdateArticle(c, articleModel); updateArticleErr != nil {
			l.Error("Failed to update article", zap.Error(updateArticleErr), zap.Int64("articleID", articleModel.ID))
			return updateArticleErr
		}
		// 2.更新文章内容
		if articleDetail.Content != articleDto.Content {
			articleContentModel := &model.ArticleContent{
				ArticleID: articleModel.ID,
				Content:   articleDto.Content,
			}
			if _, updateContentErr := dao.ArticleContentDao.UpdateContentByArticleID(c, articleContentModel); updateContentErr != nil {
				l.Error("Failed to update article content", zap.Error(updateContentErr), zap.Int64("articleID", articleModel.ID))
				return updateContentErr
			}
		}
		// 3.更新文章标签关联
		addTags, deleteTags := getNewTagRelation(articleDto.Tags, articleDetail.TagNameList)
		if len(deleteTags) > 0 {
			// 获取新增标签id
			tagIdList, listTagErr := dao.TagDao.ListTagIdByNameArr(c, deleteTags)
			if listTagErr != nil {
				l.Error("Failed to list tag id by name arr", zap.Error(listTagErr), zap.Strings("tags", deleteTags))
				return listTagErr
			}
			var relations []*model.ArticleTagRelation
			for _, tagID := range tagIdList {
				relation := &model.ArticleTagRelation{
					ArticleID: articleModel.ID,
					TagID:     tagID,
				}
				relations = append(relations, relation)
			}
			// 删除旧的标签关联
			if err := dao.ArticleTagRelationDao.DeleteArticleTagRelations(c, relations); err != nil {
				l.Error("Failed to delete article tag relation", zap.Error(err), zap.Int64("articleID", articleModel.ID))
				return err
			}
		}
		if len(addTags) > 0 {
			// 获取新增标签id
			tagIdList, listTagErr := dao.TagDao.ListTagIdByNameArr(c, addTags)
			if listTagErr != nil {
				l.Error("Failed to list tag id by name arr", zap.Error(listTagErr), zap.Strings("tags", addTags))
				return listTagErr
			}
			var relations []*model.ArticleTagRelation
			for _, tagID := range tagIdList {
				relation := &model.ArticleTagRelation{
					ArticleID: articleModel.ID,
					TagID:     tagID,
				}
				relations = append(relations, relation)
			}
			// 插入文章-标签关系
			if err := dao.ArticleTagRelationDao.InsertArticleTagRelations(c, relations); err != nil {
				l.Error("Failed to insert article tag relation", zap.Error(err), zap.Int64("articleID", articleModel.ID))
				return err
			}
		}
		return nil
	})
	if txErr != nil {
		l.Error("Failed to update article", zap.Error(txErr))
		return txErr
	}
	return nil
}

func (s *articleService) GetArticleDetail(c *gin.Context, id int64) (*vo.ArticleDetailVo, error) {
	l := logger.FromContext(c.Request.Context())
	if id <= 0 {
		return nil, cerr.NewParamError("id无效")
	}
	// 查询文章
	articleDetail, err := dao.ArticleDao.QueryArticleDetail(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.Error("Article not exist", zap.Int64("article id", id))
			return nil, cerr.New(cerr.ERROR_ARTICLE_NOT_EXIST)
		}
		l.Error("Failed to query article detail", zap.Error(err), zap.Int64("article id", id))
		return nil, err
	}
	if articleDetail == nil {
		l.Error("Article not exist", zap.Int64("article id", id))
		return nil, cerr.New(cerr.ERROR_ARTICLE_NOT_EXIST)
	}
	// 查询文章浏览量缓存
	if intCmd := cache.Client.SCard(c, utils.GetArticlePageViewKey(id)); intCmd.Err() != nil {
		l.Error("Failed to query article view from cache", zap.Error(intCmd.Err()), zap.Int64("article id", id))
	} else {
		articleDetail.Views += int(intCmd.Val())
	}

	tagList := []string{}
	if len(articleDetail.TagNameList) > 0 {
		tagList = strings.Split(articleDetail.TagNameList, ",")
	}

	resp := vo.ArticleDetailVo{
		ArticleMeta: vo.ArticleMeta{
			ID:                  articleDetail.ID,
			Title:               articleDetail.Title,
			Summary:             articleDetail.Summary,
			CategoryID:          articleDetail.CategoryID,
			Author:              articleDetail.Author,
			AllowComment:        articleDetail.AllowComment,
			Views:               articleDetail.Views,
			Weight:              articleDetail.Weight,
			IsSticky:            articleDetail.IsSticky,
			IsOriginal:          articleDetail.IsOriginal,
			OriginalArticleLink: articleDetail.OriginalArticleLink,
			Status:              articleDetail.Status,
			CreatedTime:         articleDetail.CreatedTime.UnixMilli(),
			UpdatedTime:         articleDetail.UpdatedTime.UnixMilli(),
		},
		CategoryName: articleDetail.CategoryName,
		TagNameList:  tagList,
		Content:      articleDetail.Content,
	}

	return &resp, nil
}

func (s *articleService) DeleteArticleList(c *gin.Context, deleteDto dto.ArticleDeleteDto) error {
	l := logger.FromContext(c.Request.Context())
	txErr := mysql.RunDBTransaction(c, func() error {
		// 删除文章
		if err := dao.ArticleDao.DeleteArticleByIDs(c, deleteDto.IDs); err != nil {
			l.Error("Failed to delete article by ids", zap.Error(err), zap.Int64s("ids", deleteDto.IDs))
			return err
		}
		// 删除文章内容
		if err := dao.ArticleContentDao.DeleteContentByArticleIDs(c, deleteDto.IDs); err != nil {
			l.Error("Failed to delete article content by ids", zap.Error(err), zap.Int64s("ids", deleteDto.IDs))
			return err
		}
		// 删除文章标签关联
		if err := dao.ArticleTagRelationDao.DeleteArticleTagRelationsByArticleIDs(c, deleteDto.IDs); err != nil {
			l.Error("Failed to delete article tag relation by ids", zap.Error(err), zap.Int64s("ids", deleteDto.IDs))
			return err
		}
		return nil
	})
	if txErr != nil {
		l.Error("Failed to delete article list", zap.Error(txErr))
		return txErr
	}
	return nil
}

func (s *articleService) AddPageView(c *gin.Context, pageViewDto dto.ArticlePageViewDto) error {
	l := logger.FromContext(c.Request.Context())
	_, err := s.GetArticleDetail(c, pageViewDto.ArticleID)
	if err != nil {
		l.Error("Failed to get article detail", zap.Error(err), zap.Int64("article id", pageViewDto.ArticleID))
		return err
	}

	cookie, err := c.Request.Cookie(utils.COOKIE_TEMP_USER_ID)
	hasCookie := true
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			hasCookie = false
		} else {
			l.Error("Failed to get cookie", zap.String("cookie", utils.COOKIE_TEMP_USER_ID), zap.Error(err))
			return err
		}
	}
	if cookie == nil {
		hasCookie = false
	}

	var tempUserID string
	if hasCookie {
		tempUserID = cookie.Value
	} else {
		tempUserID = utils.GenerateTempUserID(c.ClientIP(), c.Request.UserAgent())
		// 设置cookie
		c.SetCookie(
			utils.COOKIE_TEMP_USER_ID,
			tempUserID,
			int(48*time.Hour),
			"/",
			"."+config.Config.App.Domain,
			false,
			true)
	}
	key := utils.GetArticlePageViewKey(pageViewDto.ArticleID)
	if status := cache.Client.SAdd(c, key, tempUserID); status.Err() != nil {
		l.Error("Failed to set page view", zap.Error(status.Err()), zap.String("key", key))
		return status.Err()
	}

	return nil
}

// 根据新标签列表和旧标签列表，找出需要删除的标签列表和需要新增的标签列表
func getNewTagRelation(newTagList, oldTagList []string) ([]string, []string) {
	// 1. 存入map，方便查找
	oldTagSet := map[string]bool{}
	for _, tag := range oldTagList {
		oldTagSet[tag] = true
	}
	newTagSet := map[string]bool{}
	for _, tag := range newTagList {
		newTagSet[tag] = true
	}
	// 2. 找出需要删除的标签列表
	deleteTagList := []string{}
	for tag := range oldTagSet {
		if !newTagSet[tag] {
			deleteTagList = append(deleteTagList, tag)
		}
	}
	// 3. 找出需要新增的标签列表
	addTagList := []string{}
	for tag := range newTagSet {
		if !oldTagSet[tag] {
			addTagList = append(addTagList, tag)
		}
	}
	return addTagList, deleteTagList
}

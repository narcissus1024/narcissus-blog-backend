package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/narcissus1949/narcissus-blog/internal/database/cache"
	"github.com/narcissus1949/narcissus-blog/internal/logger"
	"github.com/narcissus1949/narcissus-blog/internal/utils"
	"github.com/narcissus1949/narcissus-blog/pkg/server/dao"
	"go.uber.org/zap"
)

func RunPageViewProcessor(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			// 计算到下一个凌晨3点的时间差
			now := time.Now()
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
			if now.After(nextRun) {
				nextRun = nextRun.Add(24 * time.Hour)
			}
			duration := nextRun.Sub(now)

			select {
			case <-time.After(duration):
				startTime := time.Now()
				logger.FromContext(ctx).Info("Start to sync article view count")
				if err := utils.Retry(3, 1000, func() error {
					if errList := UpdateArticleViewCount(ctx); len(errList) > 0 {
						return errList[0]
					}
					return nil
				}); err != nil {
					logger.FromContext(ctx).Error("Page view processor failed after 3 times retry")
				}
				logger.FromContext(ctx).Info("Sync article view count done", zap.Duration("cost", time.Since(startTime)))
			case <-ctx.Done():
				logger.FromContext(ctx).Info("Article page view processor stopped")
				return
			}
		}
	}(ctx)
}

func UpdateArticleViewCount(ctx context.Context) []error {
	// 获取所有文章的浏览量的key
	keysCmd := cache.Client.Keys(ctx, fmt.Sprintf(utils.ARTICLE_PAGE_VIEW_KEY_TEMPLATE, "*"))
	if keysCmd.Err() != nil {
		logger.FromContext(ctx).Error("Failed to get page view keys", zap.Error(keysCmd.Err()))
		return []error{keysCmd.Err()}
	}
	var errList []error
	logger.FromContext(ctx).Info("Get article page view keys", zap.Int("total", len(keysCmd.Val())))
	// 遍历所有的key
	for _, key := range keysCmd.Val() {
		// 获取缓存的文章的浏览量
		countCmd := cache.Client.SCard(ctx, key)
		if countCmd.Err() != nil {
			logger.FromContext(ctx).Error("Failed to get page view count", zap.Error(countCmd.Err()), zap.String("key", key))
			errList = append(errList, countCmd.Err())
			continue
		}
		// 从key中解析文章的ID
		articleID, idConvErr := utils.GetArticleIDFromPageViewKey(key)
		if idConvErr != nil {
			logger.FromContext(ctx).Error("Failed to get article id from page view key", zap.Error(idConvErr), zap.String("key", key))
			errList = append(errList, idConvErr)
			continue
		}
		// 更新文章的浏览量
		RowsAffected, err := dao.ArticleDao.IncreaseArticleViews(ctx, articleID, int(countCmd.Val()))
		if err != nil {
			logger.FromContext(ctx).Error("Failed to increase article views", zap.Error(err), zap.String("key", key))
			errList = append(errList, err)
			continue
		}
		if RowsAffected == 0 {
			logger.FromContext(ctx).Warn("Increase article views failed, RowsAffected is 0, will remove cache",
				zap.String("key", key),
				zap.Int64("viewCount", countCmd.Val()))
		}
		// 删除缓存的key
		if status := cache.Client.Del(ctx, key); status.Err() != nil {
			logger.FromContext(ctx).Error("Failed to delete page view key", zap.Error(status.Err()), zap.String("key", key))
			errList = append(errList, status.Err())
			continue
		}
	}
	return errList
}

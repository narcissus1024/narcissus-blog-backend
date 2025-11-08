package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	ARTICLE_TYPE_POST = iota
	ARTICLE_TYPE_ESSAY
	ARTICLE_TYPE_ABOUT
)

const (
	CONTEXT_USER_ID                = "UserID"
	ACCESS_TOKEN_BLACKLIST         = "access_token_blacklist:"
	REFRESH_TOKEN_BLACKLIST        = "refresh_token_blacklist:"
	ARTICLE_PAGE_VIEW_KEY_TEMPLATE = "article_page_view:%s" // article_id

	COOKIE_TEMP_USER_ID = "temp_user_id"

	X_REQUEST_ID = "X-Request-Id" // 请求ID，用于日志跟踪
)

func GetArticleIDFromPageViewKey(key string) (int64, error) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return 0, errors.New("invalid key")
	}
	articleIDStr := parts[1]
	articleID, idConvErr := strconv.ParseInt(articleIDStr, 10, 64)
	if idConvErr != nil {
		return 0, idConvErr
	}
	return articleID, nil
}

func GetArticlePageViewKey(articleID int64) string {
	return fmt.Sprintf(ARTICLE_PAGE_VIEW_KEY_TEMPLATE, strconv.FormatInt(articleID, 10))
}

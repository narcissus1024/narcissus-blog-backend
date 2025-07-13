package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/cmd/blog/app/config"
	"github.com/narcissus1949/narcissus-blog/internal/database/cache"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	"github.com/narcissus1949/narcissus-blog/internal/logger"
	"github.com/narcissus1949/narcissus-blog/internal/middleware"
	"github.com/narcissus1949/narcissus-blog/internal/validator"
	"github.com/narcissus1949/narcissus-blog/pkg/route"
)

func Run(configPath string) {
	ctx := context.Background()
	MustInit(ctx, configPath)
	if err := StartServer(ctx); err != nil {
		panic(err.Error())
	}
}

func MustInit(ctx context.Context, configPath string) {
	config.MustInit(configPath)
	logger.MustInit(config.Config.Logger)
	cache.MustInit(config.Config.Redis)
	mysql.MustInit(config.Config.Mysql)
	validator.MustRegistValidator()

}

func StartServer(ctx context.Context) error {
	g := gin.Default()
	g.Use(middleware.GinLogger(), middleware.GinRecovery(true), middleware.Cors())
	route.Setup(g)
	return g.Run(fmt.Sprintf(":%d", config.Config.App.Port))
}

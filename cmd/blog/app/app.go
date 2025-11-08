package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/cmd/blog/app/config"
	"github.com/narcissus1949/narcissus-blog/internal/database/cache"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	"github.com/narcissus1949/narcissus-blog/internal/logger"
	"github.com/narcissus1949/narcissus-blog/internal/middleware"
	"github.com/narcissus1949/narcissus-blog/internal/validator"
	"github.com/narcissus1949/narcissus-blog/pkg/route"
	"github.com/narcissus1949/narcissus-blog/pkg/server/processor"
	"go.uber.org/zap"
)

func Run(configPath string) {
	ctx, cancel := context.WithCancel(context.Background())
	RegisterShutdown(ctx, cancel)
	MustInit(ctx, configPath)
	if err := StartServer(ctx); err != nil {
		panic(err.Error())
	}
}

func RegisterShutdown(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		// 监听信号
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
		// 阻塞等待信号
		<-sigChan
		zap.L().Info("Shutdown Server ...")
		cancel()
	}()
}

func MustInit(ctx context.Context, configPath string) {
	config.MustInit(configPath)
	logger.MustInit(config.Config.Logger)
	cache.MustInit(config.Config.Redis)
	mysql.MustInit(config.Config.Mysql)
	validator.MustRegistValidator()

	// 更新文章浏览量
	processor.RunPageViewProcessor(ctx)
}

func StartServer(ctx context.Context) error {
	g := gin.Default()
	g.Use(middleware.GinLogger(), middleware.GinRecovery(true), middleware.Cors())
	route.Setup(g)
	return g.Run(fmt.Sprintf(":%d", config.Config.App.Port))
}

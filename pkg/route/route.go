package route

import (
	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/middleware"
	"github.com/narcissus1949/narcissus-blog/pkg/server/handler"
)

func Setup(gin *gin.Engine) {

	g := gin.Group("")

	// 用户路由
	userRoute := g.Group("/user")
	{
		userRoute.POST("/signin", handler.UserHandler.SignIn)
		userRoute.POST("/login", handler.UserHandler.Login)
	}

	// 文章路由
	articleRoute := g.Group("/article")
	{
		articleRoute.POST("/views", handler.ArticleHandler.IncreasePageView)

		articleRoute.POST("/list", handler.ArticleHandler.ListArticle)
		articleRoute.GET("/detail", handler.ArticleHandler.GetArticleeDetail)

		// 文章分类路由
		articleRoute.GET("/category/listAll", handler.CategoryHandler.ListAllCategory)
		articleRoute.GET("/category/list", handler.CategoryHandler.ListCategory)
		articleRoute.GET("/category/get", handler.CategoryHandler.GetCategoryDetail)

		// 文章标签路由
		articleRoute.GET("/tag/listAll", handler.TagHandler.ListAllTag)
		articleRoute.GET("/tag/list", handler.TagHandler.ListTag)
		articleRoute.GET("/tag/get", handler.TagHandler.GetTagDetail)
	}

	// 通用路由
	commonRoute := g.Group("/common")
	commonRoute.GET("/ssl", handler.CommonHandler.GetRASPublicKey)
	commonRoute.POST("/ssl/encrypt", handler.CommonHandler.PublicKeyEncrypt)

	// 需要权限路由
	userAuthRoute := g.Group("/user", middleware.JWTAuth())
	{
		userAuthRoute.POST("/logout", handler.UserHandler.Logout)
		userAuthRoute.POST("/token/refresh", handler.UserHandler.RefreshToken)
	}

	// 文章
	articleAuthRoute := g.Group("/article", middleware.JWTAuth())
	{
		// 文章
		articleAuthRoute.POST("/admin/list", handler.ArticleHandler.ListArticleAdmin)
		articleAuthRoute.POST("/save", handler.ArticleHandler.SaveArticle)
		articleAuthRoute.POST("/delete", handler.ArticleHandler.DeleteArticleList)

		// 分类
		articleAuthRoute.POST("/category/create", handler.CategoryHandler.CreateCategoryList)
		articleAuthRoute.POST("/category/update", handler.CategoryHandler.UpdateCategory)
		articleAuthRoute.POST("/category/delete", handler.CategoryHandler.DeleteCategoryList)
		// 标签
		articleAuthRoute.POST("/tag/create", handler.TagHandler.CreateTagList)
		articleAuthRoute.POST("/tag/update", handler.TagHandler.UpdateTag)
		articleAuthRoute.POST("/tag/delete", handler.TagHandler.DeleteTagList)
	}

	// 通用
	commonAuthRoute := g.Group("/common", middleware.JWTAuth())
	commonAuthRoute.POST("/upload/image", handler.CommonHandler.UploadImage)

}

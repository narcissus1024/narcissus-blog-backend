package route

import (
	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/middleware"
	"github.com/narcissus1949/narcissus-blog/pkg/server/controller"
)

func Setup(gin *gin.Engine) {

	g := gin.Group("")

	// 用户路由
	userRoute := g.Group("/user")
	{
		userRoute.POST("/signin", controller.UserController.SignIn)
		userRoute.POST("/login", controller.UserController.Login)
	}

	// 文章路由
	articleRoute := g.Group("/article")
	{
		articleRoute.POST("/list", controller.ArticleController.ListArticle)
		articleRoute.GET("/detail", controller.ArticleController.GetArticleeDetail)

		// 文章分类路由
		articleRoute.GET("/category/listAll", controller.CategoryController.ListAllCategory)
		articleRoute.GET("/category/list", controller.CategoryController.ListCategory)
		articleRoute.GET("/category/get", controller.CategoryController.GetCategoryDetail)

		// 文章标签路由
		articleRoute.GET("/tag/listAll", controller.TagController.ListAllTag)
		articleRoute.GET("/tag/list", controller.TagController.ListTag)
		articleRoute.GET("/tag/get", controller.TagController.GetTagDetail)
	}

	// 通用路由
	commonRoute := g.Group("/common")
	commonRoute.GET("/ssl", controller.CommonController.GetRASPublicKey)
	commonRoute.POST("/ssl/encrypt", controller.CommonController.PublicKeyEncrypt)

	// 需要权限路由
	userAuthRoute := g.Group("/user", middleware.JWTAuth())
	{
		userAuthRoute.POST("/logout", controller.UserController.Logout)
		userAuthRoute.POST("/token/refresh", controller.UserController.RefreshToken)
	}

	// 文章
	articleAuthRoute := g.Group("/article", middleware.JWTAuth())
	{
		// 文章
		articleAuthRoute.POST("/admin/list", controller.ArticleController.ListArticleAdmin)
		articleAuthRoute.POST("/save", controller.ArticleController.SaveArticle)
		articleAuthRoute.POST("/delete", controller.ArticleController.DeleteArticleList)

		// 分类
		articleAuthRoute.POST("/category/create", controller.CategoryController.CreateCategoryList)
		articleAuthRoute.POST("/category/update", controller.CategoryController.UpdateCategory)
		articleAuthRoute.POST("/category/delete", controller.CategoryController.DeleteCategoryList)
		// 标签
		articleAuthRoute.POST("/tag/create", controller.TagController.CreateTagList)
		articleAuthRoute.POST("/tag/update", controller.TagController.UpdateTag)
		articleAuthRoute.POST("/tag/delete", controller.TagController.DeleteTagList)
	}

	// 通用
	commonAuthRoute := g.Group("/common", middleware.JWTAuth())
	commonAuthRoute.POST("/upload/image", controller.CommonController.UploadImage)

}

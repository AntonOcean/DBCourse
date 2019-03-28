package main

import (
	"DBCourse/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	user := r.Group("api/user")
	{
		user.POST("/:nickname/create", handlers.UserCreateHandler)
		user.GET("/:nickname/profile", handlers.UserProfileHandler)
		user.POST("/:nickname/profile", handlers.UserProfileHandler)
	}

	forum := r.Group("api/forum")
	{
		forum.POST("/:path1", GetHandler)
		forum.GET("/:slug/details", handlers.ForumDetailHandler)
		forum.POST("/:path1/:path2", GetHandler)
		forum.GET("/:slug/threads", handlers.ForumThreadList)
		forum.GET("/:slug/users", handlers.ForumUserList)
	}

	thread := r.Group("api/thread")
	{
		thread.POST("/:identifier/create", handlers.PostsCreateHandler)
		thread.POST("/:identifier/vote", handlers.VoteCreateHandler)
		thread.GET("/:identifier/details", handlers.DetailsHandler)
		thread.POST("/:identifier/details", handlers.DetailsHandler)

		thread.GET("/:identifier/posts", handlers.GetPostsHandler)
	}

	post := r.Group("api/post")
	{
		post.GET("/:id/details", handlers.PostDetailsHandler)
		post.POST("/:id/details", handlers.PostDetailsHandler)
	}

	service := r.Group("api/service")
	{
		service.GET("/status", handlers.ServiceStatusHandler)
		service.POST("/clear", handlers.ServiceClearHandler)
	}

	//r.Use(DataBaseConnectionMiddleware())

	_ = r.Run(":5000") // listen and serve on 0.0.0.0:8080
}

//router.GET("/v1/images/:path1", GetHandler)           //      /v1/images/detail
//    router.GET("/v1/images/:path1/:path2", GetHandler)    //      /v1/images/<id>/history
//
func GetHandler(c *gin.Context) {
	path1 := c.Param("path1")
	path2 := c.Param("path2")

	if path1 == "create" && path2 == "" {
		handlers.ForumCreateHandler(c)
	} else if path1 != "" && path2 == "create" {
		slug := path1
		handlers.ThreadCreateHandler(c, slug)
	}
}

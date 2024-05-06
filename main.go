package main

import (
	"golang-blog-application/handlers"
	"golang-blog-application/helpers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(helpers.CORSMiddleware())

	r.POST("/api/create", handlers.CreatePost)

	r.GET("/api/getAllBlogs", handlers.GetAllBlogs)

	r.GET("/api/getBlog/:id", handlers.GetBlogById)

	r.DELETE("/api/deleteBlog/:id", handlers.DeleteBlog)

	r.PUT("/api/updateBlog/:id", handlers.UpdatePost)

	r.POST("/api/createComment/:id", handlers.CreateComment)

	r.DELETE("/api/deleteComment/:id", handlers.DeleteComment)

	r.POST("/api/likeBlog/:id", handlers.LikeBlog)

	r.Run()
}

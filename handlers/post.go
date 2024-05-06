package handlers

import (
	"context"
	"golang-blog-application/helpers"
	"golang-blog-application/prisma/db"

	"github.com/gin-gonic/gin"
)

type Create struct {
	AuthorId string `json:"authorId" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type Delete struct {
	UserId string `json:"userId" binding:"required"`
}

type Update struct {
	UserId  string `json:"userId" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func UpdatePost(c *gin.Context) {
	client := helpers.CreateClient()
	var userInput Update
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	ctx := context.Background()

	_, err := client.Post.FindMany(
		db.Post.ID.Equals(id),
		db.Post.AuthorID.Equals(userInput.UserId),
	).Update(
		db.Post.Title.Set(userInput.Title),
		db.Post.Content.Set(userInput.Content),
	).Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	helpers.CloseClient(client)

	c.JSON(200, gin.H{
		"message": "Post updated successfully",
	})
}

func CreatePost(c *gin.Context) {
	client := helpers.CreateClient()
	var userInput Create
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	createdPost, err := client.Post.CreateOne(
		db.Post.Title.Set(userInput.Title),
		db.Post.Content.Set(userInput.Content),
		db.Post.AuthorID.Set(userInput.AuthorId),
	).Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	helpers.CloseClient(client)

	c.JSON(200, gin.H{
		"message": "Post created successfully",
		"data":    createdPost.ID,
	})
}

func GetAllBlogs(c *gin.Context) {
	client := helpers.CreateClient()
	ctx := context.Background()

	posts, err := client.Post.FindMany().With(
		db.Post.Comments.Fetch(),
		db.Post.Likes.Fetch(),
	).OrderBy(
		db.Post.CreatedAt.Order(db.SortOrderDesc),
	).Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	helpers.CloseClient(client)

	c.JSON(200, gin.H{
		"data": posts,
	})
}

func GetBlogById(c *gin.Context) {
	client := helpers.CreateClient()
	id := c.Param("id")
	ctx := context.Background()

	post, err := client.Post.FindUnique(
		db.Post.ID.Equals(id),
	).With(
		db.Post.Comments.Fetch().OrderBy(
			db.Comment.CreatedAt.Order(db.SortOrderDesc),
		),
		db.Post.Likes.Fetch(),
	).Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	helpers.CloseClient(client)

	c.JSON(200, gin.H{
		"data": post,
	})
}

func DeleteBlog(c *gin.Context) {
	client := helpers.CreateClient()
	var userInput Delete
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	ctx := context.Background()

	_, err := client.Post.FindMany(
		db.Post.ID.Equals(id),
		db.Post.AuthorID.Equals(userInput.UserId),
	).Delete().Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	helpers.CloseClient(client)

	c.JSON(200, gin.H{
		"message": "Post deleted successfully",
	})
}

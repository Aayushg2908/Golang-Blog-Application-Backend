package handlers

import (
	"context"
	"golang-blog-application/helpers"
	"golang-blog-application/prisma/db"

	"github.com/gin-gonic/gin"
)

type createComment struct {
	AuthorId string `json:"authorId" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type deleteComment struct {
	UserId string `json:"userId" binding:"required"`
}

func CreateComment(c *gin.Context) {
	client := helpers.CreateClient()
	var create createComment
	if err := c.ShouldBindJSON(&create); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	postId := c.Param("id")

	ctx := context.Background()

	_, err := client.Comment.CreateOne(
		db.Comment.Content.Set(create.Content),
		db.Comment.AuthorID.Set(create.AuthorId),
		db.Comment.Post.Link(
			db.Post.ID.Equals(postId),
		),
	).Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	helpers.CloseClient(client)

	c.JSON(200, gin.H{
		"message": "Comment created successfully",
	})
}

func DeleteComment(c *gin.Context) {
	client := helpers.CreateClient()
	var delete deleteComment
	if err := c.ShouldBindJSON(&delete); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	ctx := context.Background()

	_, err := client.Comment.FindMany(
		db.Comment.ID.Equals(id),
		db.Comment.AuthorID.Equals(delete.UserId),
	).Delete().Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	helpers.CloseClient(client)

	c.JSON(200, gin.H{
		"message": "Comment deleted successfully",
	})
}

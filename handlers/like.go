package handlers

import (
	"context"
	"golang-blog-application/helpers"
	"golang-blog-application/prisma/db"

	"github.com/gin-gonic/gin"
)

type like struct {
	UserId string `json:"userId" binding:"required"`
}

func LikeBlog(c *gin.Context) {
	client := helpers.CreateClient()
	var json like
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	postId := c.Param("id")
	ctx := context.Background()

	_, err := client.Like.FindFirst(
		db.Like.PostID.Equals(postId),
		db.Like.AuthorID.Equals(json.UserId),
	).Exec(ctx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			_, err = client.Like.CreateOne(
				db.Like.AuthorID.Set(json.UserId),
				db.Like.Post.Link(
					db.Post.ID.Equals(postId),
				),
			).Exec(ctx)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			helpers.CloseClient(client)

			c.JSON(200, gin.H{
				"message": "Post liked successfully",
			})
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err = client.Like.FindMany(
		db.Like.PostID.Equals(postId),
		db.Like.AuthorID.Equals(json.UserId),
	).Delete().Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Like removed successfully",
	})
}

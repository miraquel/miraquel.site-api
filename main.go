package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"

	"miraquel.site/api/model"
	"miraquel.site/api/repository"
)

func main() {
	r := gin.New()

	db, err := sql.Open("pgx", "postgres://postgres:Babakan13!@localhost:5432/miraquel")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r.POST("/users", func(c *gin.Context) {
		var user model.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// save user to database
		userRepo := repository.NewRepositoryUserPsql(db)

		user.RegisteredAt = time.Now()
		newUser, err := userRepo.Create(c, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, newUser)
	})

	r.GET("/users", func(c *gin.Context) {
		userRepo := repository.NewRepositoryUserPsql(db)
		users, err := userRepo.All(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	r.GET("/users/:id/posts", func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("id"), 0, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		userRepo := repository.NewRepositoryUserPsql(db)
		users, err := userRepo.GetByIdWithPosts(c, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	r.POST("/posts", func(c *gin.Context) {
		var post model.Post

		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// save user to database
		postRepo := repository.NewRepositoryPostPsql(db)

		post.Published = 0
		post.CreatedAt = time.Now()
		post.UpdatedAt = time.Now()

		newPost, err := postRepo.Create(c, post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, newPost)
	})

	r.GET("/posts", func(c *gin.Context) {
		userRepo := repository.NewRepositoryPostPsql(db)
		users, err := userRepo.All(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	r.Run()
}

func ParseInt(s string) {
	panic("unimplemented")
}

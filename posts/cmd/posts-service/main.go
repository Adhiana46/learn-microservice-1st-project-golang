package main

import (
	"net/http"

	"github.com/adhiana46/posts-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var posts map[string]Post = make(map[string]Post)

func main() {
	r := echo.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.GET("/posts", handleGetPosts)
	r.POST("/posts/create", handleCreatePosts)

	r.Logger.Fatal(r.Start(":4000"))
}

func handleGetPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, posts)
}

func handleCreatePosts(c echo.Context) error {
	id, _ := utils.RandomHex(4)
	post := Post{}

	if err := c.Bind(&post); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	post.ID = id

	// store in-memory
	posts[id] = post

	// TODO: sendEvent PostCreated

	return c.JSON(http.StatusCreated, post)
}

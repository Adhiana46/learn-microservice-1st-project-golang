package main

import (
	"log"
	"net/http"

	"github.com/adhiana46/comments-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Comment struct {
	ID      string `json:"id"`
	PostID  string `json:"postId"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:data"`
}

var commentsByPostId map[string][]Comment = make(map[string][]Comment)

func main() {
	r := echo.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.GET("/posts/:id/comments", handleGetComments)
	r.POST("/posts/:id/comments", handleCreateComment)

	r.POST("/events", handleEvents)

	r.Logger.Fatal(r.Start(":4001"))
}

func handleGetComments(c echo.Context) error {
	comments := commentsByPostId[c.Param("id")]

	if comments == nil {
		return c.JSON(http.StatusOK, []Comment{})
	}

	return c.JSON(http.StatusOK, comments)
}

func handleCreateComment(c echo.Context) error {
	postId := c.Param("id")
	commentId, _ := utils.RandomHex(4)

	comment := Comment{}
	if err := c.Bind(&comment); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	comment.ID = commentId
	comment.PostID = postId
	comment.Status = "pending"

	comments := commentsByPostId[postId]
	if comments == nil {
		comments = []Comment{}
	}

	comments = append(comments, comment)

	commentsByPostId[postId] = comments

	// TODO: send event CommentCreated

	return c.JSON(http.StatusCreated, comments)
}

func handleEvents(c echo.Context) error {
	event := Event{}

	if err := c.Bind(&event); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	log.Println("Reveived Event", event.Type)

	return c.NoContent(http.StatusOK)
}

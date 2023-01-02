package main

import (
	"bytes"
	"encoding/json"
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
	Data interface{} `json:"data"`
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

func sendEvent(url string, event Event) error {
	log.Println("sendEvent", url, event)

	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
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
	event := Event{
		Type: "CommentCreated",
		Data: comment,
	}

	sendEvent("http://event-bus-srv:4005/events", event)

	return c.JSON(http.StatusCreated, comments)
}

func handleEvents(c echo.Context) error {
	event := Event{}

	if err := c.Bind(&event); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	log.Println("Reveived Event", event.Type)

	if event.Type == "CommentModerated" {
		payload := event.Data.(Comment)

		comments := commentsByPostId[payload.PostID]

		for i := 0; i < len(comments); i++ {
			if comments[i].ID == payload.ID {
				comments[i].Status = payload.Status
			}
		}

		sendEvent("http://event-bus-srv:4005/events", Event{
			Type: "CommentUpdated",
			Data: payload,
		})
	}

	return c.NoContent(http.StatusOK)
}

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/adhiana46/posts-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var posts map[string]Post = make(map[string]Post)

func main() {
	r := echo.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.GET("/posts", handleGetPosts)
	r.POST("/posts/create", handleCreatePosts)

	r.POST("/events", handleEvents)

	r.Logger.Fatal(r.Start(":4000"))
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

func handleGetPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, posts)
}

func handleCreatePosts(c echo.Context) error {
	id, _ := utils.RandomHex(4)
	post := Post{}

	if err := c.Bind(&post); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	post.ID = id

	// store in-memory
	posts[id] = post

	// TODO: sendEvent PostCreated
	sendEvent("http://event-bus-srv:4005/events", Event{
		Type: "PostCreated",
		Data: post,
	})

	return c.JSON(http.StatusCreated, post)
}

func handleEvents(c echo.Context) error {
	event := Event{}

	if err := c.Bind(&event); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	log.Println("Reveived Event", event.Type)

	return c.NoContent(http.StatusOK)
}

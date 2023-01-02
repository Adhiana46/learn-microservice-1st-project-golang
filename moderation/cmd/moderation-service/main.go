package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

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

func main() {
	r := echo.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.POST("/events", handleEvents)

	r.Logger.Fatal(r.Start(":4003"))
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

func handleEvent(event Event) error {
	if event.Type == "CommentCreated" {
		comment := Comment{}

		jsonData, err := json.Marshal(event.Data)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(jsonData, &comment); err != nil {
			return err
		}

		if strings.Contains(comment.Content, "orange") || strings.Contains(comment.Content, "apple") {
			comment.Status = "rejected"
		} else {
			comment.Status = "approved"
		}

		time.Sleep(5 * time.Second)

		// send event
		event := Event{
			Type: "CommentModerated",
			Data: comment,
		}

		sendEvent("http://event-bus-srv:4005/events", event)
	}

	return nil
}

func handleEvents(c echo.Context) error {
	event := Event{}

	if err := c.Bind(&event); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	log.Println("Reveived Event", event.Type)

	handleEvent(event)

	return c.NoContent(http.StatusOK)
}

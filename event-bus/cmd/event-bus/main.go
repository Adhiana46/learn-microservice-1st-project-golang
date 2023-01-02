package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var events []Event = []Event{}

func main() {
	r := echo.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.GET("/events", handleGetEvent)
	r.POST("/events", handlePostEvent)

	r.Logger.Fatal(r.Start(":4005"))
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

func handleGetEvent(c echo.Context) error {
	return c.JSON(http.StatusOK, events)
}

func handlePostEvent(c echo.Context) error {
	event := Event{}

	if err := c.Bind(&event); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	log.Println("Received Event: ", event.Type)

	events = append(events, event)

	sendEvent("http://posts-clusterip-srv:4000/events", event)
	sendEvent("http://comments-clusterip-srv:4001/events", event)
	sendEvent("http://query-clusterip-srv:4002/events", event)
	sendEvent("http://moderation-clusterip-srv:4003/events", event)

	return c.NoContent(http.StatusOK)
}

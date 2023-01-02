package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Post struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Comments []Comment `json:"comments"`
}

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

var posts map[string]Post = make(map[string]Post)

func main() {
	r := echo.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.GET("/posts", handleGetPosts)
	r.POST("/events", handleEvents)

	go func() {
		events, err := getEvents()
		log.Println("Events", events)
		if err != nil {
			log.Println("Error processing events: ", err)
		}

		for _, event := range events {
			log.Println("Processing Event: ", event.Type)
			handleEvent(event)
		}
	}()

	r.Logger.Fatal(r.Start(":4002"))
}

func getEvents() ([]Event, error) {
	request, err := http.NewRequest("GET", "http://event-bus-srv:4005/events", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var events []Event
	err = json.NewDecoder(response.Body).Decode(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func handleEvent(event Event) error {
	if event.Type == "PostCreated" {
		post := Post{}

		jsonData, err := json.Marshal(event.Data)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(jsonData, &post); err != nil {
			return err
		}

		post.Comments = []Comment{}

		posts[post.ID] = post
	}

	if event.Type == "CommentCreated" {
		comment := Comment{}

		jsonData, err := json.Marshal(event.Data)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(jsonData, &comment); err != nil {
			return err
		}

		post, isExists := posts[comment.PostID]

		if isExists {
			comments := post.Comments
			comments = append(comments, comment)

			post.Comments = comments

			posts[comment.PostID] = post
		}
	}

	if event.Type == "CommentUpdated" {
		comment := Comment{}

		jsonData, err := json.Marshal(event.Data)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(jsonData, &comment); err != nil {
			return err
		}

		_, isPostExists := posts[comment.PostID]
		if isPostExists {
			for i := 0; i < len(posts[comment.PostID].Comments); i++ {
				if comment.ID == posts[comment.PostID].Comments[i].ID {
					posts[comment.PostID].Comments[i].Status = comment.Status
					posts[comment.PostID].Comments[i].Content = comment.Content
				}
			}
		}
	}

	return nil
}

func handleGetPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, posts)
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

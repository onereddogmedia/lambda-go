package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// APIKey the TMDB key
	APIKey = os.Getenv("API_KEY")
	// ErrorBackend general error
	ErrorBackend = errors.New("Something went wrong")
)

// Request client request
type Request struct {
	ID int `json:"id"`
}

// MovieDBResponse response data
type MovieDBResponse struct {
	Movies []Movie `json:"results"`
}

// Movie individual movie
type Movie struct {
	Title       string `json:"title"`
	Description string `json:"overview"`
	Cover       string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
}

// Handler client request handler
func Handler(rctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/discover/movie?api_key=%s", APIKey)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{}, ErrorBackend
	}

	var parameters Request
	log.Print("Body=", request.Body)
	if err := json.Unmarshal([]byte(request.Body), &parameters); err == nil {
		log.Print("id=", parameters.ID)
		if parameters.ID > 0 {
			q := req.URL.Query()
			q.Add("with_genres", strconv.Itoa(parameters.ID))
			req.URL.RawQuery = q.Encode()
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return events.APIGatewayProxyResponse{}, ErrorBackend
	}
	defer resp.Body.Close()

	var data MovieDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return events.APIGatewayProxyResponse{}, ErrorBackend
	}

	body, err := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}

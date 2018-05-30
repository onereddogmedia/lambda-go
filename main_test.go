package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()
	event := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/",
		Body:       `{ "id": 878 }`,
	}

	movies, err := Handler(ctx, event)

	assert.IsType(t, nil, err)
	assert.NotEqual(t, 0, len(movies.Body))
	fmt.Print(movies.Body)
}

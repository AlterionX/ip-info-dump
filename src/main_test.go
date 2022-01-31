package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func Test_runIPLookup_success(t *testing.T) {
	req := events.APIGatewayProxyRequest {}
	req.QueryStringParameters = map[string]string {
		"q": "example.com",
	}

	response, err := runIPLookup(context.Background(), req)
	assert.Nil(t, err, "error from lambda handler is always nil")
	assert.Equal(t, http.StatusOK, response.StatusCode, "response to be okay")
	assert.NotEmpty(t, response.Body, "some data to be returned")
}

func Test_runIPLookup_badArgs(t *testing.T) {
	req := events.APIGatewayProxyRequest {}
	req.QueryStringParameters = map[string]string {
		"q": "thisisnotawebsite",
	}

	response, err := runIPLookup(context.Background(), req)
	assert.Nil(t, err, "error from lambda handler is always nil")
	assert.Equal(t, http.StatusBadRequest, response.StatusCode, "response to not be okay")
	assert.NotEmpty(t, response.Body, "some data to be returned")
}


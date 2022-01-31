package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_responses(t *testing.T) {
	client := clientErrorResponse(http.StatusBadRequest)
	assert.Equal(t, http.StatusBadRequest, client.StatusCode, "client response to be fine")

	server := serverErrorResponse()
	assert.Equal(t, http.StatusInternalServerError, server.StatusCode, "client response to be fine")
}

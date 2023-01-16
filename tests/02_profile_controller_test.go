package tests

import (
	"encoding/json"
	"goapi/request"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	res := makeRequest("GET", "/api/v1/profile", body, true)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)
}

func TestUpdateProfile(t *testing.T) {
	data := request.ProfileInput{
		FirstName: "Otto Eliu",
		LastName:  "Florez",
		Phone:     "55555551",
	}

	res := makeRequest("POST", "/api/v1/profile", data, true)
	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)
}

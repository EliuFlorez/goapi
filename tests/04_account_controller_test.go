package tests

import (
	"encoding/json"
	"goapi/request"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAccount(t *testing.T) {
	data := request.AccountInput{
		Name: "iDev",
	}

	res := makeRequest("POST", "/api/v1/accounts", data, true)
	printHResponse(res)

	assert.Equal(t, http.StatusCreated, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)
}

func TestGetAccounts(t *testing.T) {
	res := makeRequest("GET", "/api/v1/accounts", body, true)
	printHResponse(res)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)
}

func TestGetAccount(t *testing.T) {
	res := makeRequest("GET", "/api/v1/accounts/1", body, true)
	printHResponse(res)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)
}

func TestUpdateAccount(t *testing.T) {
	data := request.AccountInput{
		Name: "iDev",
	}

	res := makeRequest("PUT", "/api/v1/accounts/1", data, true)
	printHResponse(res)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)
}

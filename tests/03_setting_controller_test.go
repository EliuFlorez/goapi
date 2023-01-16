package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserPassword(t *testing.T) {
	res := makeRequest("POST", "/api/v1/setting/password", body, true)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)
}

func TestUserTwoFa(t *testing.T) {
	res := makeRequest("POST", "/api/v1/setting/twofa", body, true)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)
}

func TestEmailReset(t *testing.T) {
	res := makeRequest("POST", "/api/v1/setting/email", body, true)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)
}

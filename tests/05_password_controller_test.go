package tests

import (
	"encoding/json"
	"goapi/request"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordForgot(t *testing.T) {
	data := request.EmailInput{
		Email: "demo@demo.com",
	}

	res := makeRequest("POST", "/auth/password/forgot", data, false)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)
}

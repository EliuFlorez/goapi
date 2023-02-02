package tests

import (
	"encoding/json"
	"goapi/model"
	"goapi/request"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	count := model.TotalUsers()

	if count <= 0 {
		data := request.SignUpInput{
			CompanyName:          "iDev",
			FirstName:            "Eliu",
			LastName:             "Florez",
			Email:                "demo@demo.com",
			Password:             "123456789",
			PasswordConfirmation: "123456789",
		}

		res := makeRequest("POST", "/auth/signup", data, false)
		printHResponse(res)

		assert.Equal(t, http.StatusOK, res.Code)

		var response map[string]string
		json.Unmarshal(res.Body.Bytes(), &response)

		printResponse(response)
	}
}

func TestSignin(t *testing.T) {
	data := request.SignInInput{
		Email:    "demo@demo.com",
		Password: "123456789",
	}

	res := makeRequest("POST", "/auth/signin", data, false)
	printHResponse(res)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)

	printResponse(response)

	_, exists := response["token"]

	assert.Equal(t, true, exists)

	CreateSession(response["token"])
}

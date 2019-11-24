package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hrshadhin/license-server/controllers"
)

var token = ""

func TestHealthCheck(t *testing.T) {

	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Welcome)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got => '%v' want => '%v'",
			status, http.StatusBadRequest)
	}

}

func TestLogin(t *testing.T) {

	var jsonStr = []byte(`{
		"email": "dev@hrshadhin.me",
		"password": "dev#321"
	}`)

	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Authenticate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got => '%v' want => '%v'",
			status, http.StatusBadRequest)
	}
	var response map[string]interface{}
	json.Unmarshal([]byte(rr.Body.String()), &response)

	expected := `Logged In`
	getmessage := response["message"].(string)
	if getmessage != expected {
		t.Errorf("Login fail!: got => '%v' want => '%v'", getmessage, expected)
	} else {
		loginUser := response["user"].(map[string]interface{})
		token = loginUser["token"].(string)
	}

}

func TestUserCreate(t *testing.T) {

	var jsonStr = []byte(`{
		"name": "H.R.S",
		"email": "hrs@hrshadhin.me",
		"Password": "demo123"
	}`)

	req, err := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	var authToken = fmt.Sprintf("Barer %v", token)
	req.Header.Set("Authorization", authToken)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got => '%v' want => '%v'",
			status, http.StatusBadRequest)
	}

	var response map[string]interface{}
	json.Unmarshal([]byte(rr.Body.String()), &response)

	expected := `User has been created`
	getmessage := response["message"].(string)
	if getmessage != expected {
		t.Errorf("User create failed!: got => '%v' want => '%v'", getmessage, expected)
	}
}

func TestUserExists(t *testing.T) {

	var jsonStr = []byte(`{
		"name": "H.R.S",
		"email": "hrs@hrshadhin.me",
		"Password": "demo123"
	}`)

	req, err := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	var authToken = fmt.Sprintf("Barer %v", token)
	req.Header.Set("Authorization", authToken)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got => '%v' want => '%v'",
			status, http.StatusBadRequest)
	}

	var response map[string]interface{}
	json.Unmarshal([]byte(rr.Body.String()), &response)

	expected := `Email address already in use by another user.`
	getmessage := response["message"].(string)
	if getmessage != expected {
		t.Errorf("User create failed!: got => '%v' want => '%v'", getmessage, expected)
	}
}

func TestGetUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.UserList)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got => '%v' want => '%v'",
			status, http.StatusBadRequest)
	}

	var response map[string]interface{}
	json.Unmarshal([]byte(rr.Body.String()), &response)

	expected := `User List`
	getmessage := response["message"].(string)
	if getmessage != expected {
		t.Errorf("User list: got => '%v' want => '%v'", getmessage, expected)
	} else {
		users := response["users"].([]interface{})
		if len(users) != 2 {
			t.Errorf("User count not match: got => '%v' want => '%v'", len(users), 2)
		}
	}

}

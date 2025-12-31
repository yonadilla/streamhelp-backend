package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"streamhelper-backend/internal/entity"
	"streamhelper-backend/internal/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T){
	ClearAll()
	requestBody := model.RegisterUserRequest{
		ID: "Mousetri",
		Password: "hayolo",
		Name: "Mousetri janedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost,"/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	
	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBody.ID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestRegisterError(t *testing.T){
	ClearAll()
	requestBody := model.RegisterUserRequest{
		ID: "",
		Password: "",
		Name: "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestRegisterDuplicate(t *testing.T){
	ClearAll()
	TestRegister(t)

	requestBody := model.RegisterUserRequest{
		ID: "Mousetri",
		Password: "hayolo",
		Name: "Mousetri janedy",
	}

	
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost,"/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	
	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict , response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestLogin(t *testing.T){
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		ID: "Mousetri",
		Password: "hayolo",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.Token)

	user := new(entity.User)
	err = DB.Where("id = ? ", requestBody.ID).First(user).Error
	assert.Nil(t, err)
	assert.Equal(t, user.Token, responseBody.Data.Token)
}

func TestLoginWrongUsername(t *testing.T){
	ClearAll()
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		ID: "Mouse",
		Password: "hayolo",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestLoginWrongPassword(t *testing.T){
	ClearAll()
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		ID: "Mousetri",
		Password: "hayo",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestLogout(t *testing.T){
	ClearAll()
	TestLogin(t)

	user := new(entity.User)
	err := DB.Where("id = ?", "Mousetri").First(user).Error
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response , err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestLogoutWrongAuthorization(t *testing.T){
	ClearAll()
	TestLogin(t)

	request := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong")

	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}


func TestGetCurrentUser(t *testing.T){
	ClearAll()
	TestLogin(t)

	user := new(entity.User)
	err := DB.Where("id = ?", "Mousetri").First(user).Error
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t , err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, user.Name, responseBody.Data.Name)
	assert.Equal(t, user.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, user.UpdatedAt, responseBody.Data.UpdatedAt)

}

func TestGetCurrentUserFailed(t *testing.T){
	ClearAll()
	TestLogin(t)

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong")

	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t , err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestUpdateUserName(t *testing.T){
	ClearAll()
	TestLogin(t)

	user := new(entity.User)
	err := DB.Where("id = ?", "Mousetri").First(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateUserRequest{
		Name : "Mouse",
	}

	bodyJson , err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t , err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateUserPassword(t *testing.T){
	ClearAll()
	TestLogin(t)

	user := new(entity.User)
	err := DB.Where("id = ?", "Mousetri").First(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateUserRequest{
		Password: "rahasia",
	}

	bodyJson , err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t , err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	user = new(entity.User)
	err = DB.Where("id = ? ", "Mousetri").First(user).Error
	assert.Nil(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	assert.Nil(t, err)
}


func TestUpdateFailed(t *testing.T) {
	ClearAll()
	TestLogin(t)

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong")

	response, err := App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}


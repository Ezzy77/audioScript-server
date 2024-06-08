package main

import (
	"encoding/json"
	"fmt"
	"github.com/Ezzy77/audioScript-server/model"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"net/http"
)

func (app *application) loginHandler(ctx *gin.Context) {
	var authRequest model.AuthRequest
	err := ctx.ShouldBindJSON(&authRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user with Supabase
	client, err := supabase.NewClient(app.supabaseConfig.ApiUrl, app.supabaseConfig.ApiKey, nil)
	if err != nil {
		app.logger.Println("Failed to create Supabase client")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Supabase client"})
		return
	}

	resp, err := client.Auth.SignInWithEmailPassword(authRequest.Email, authRequest.Password)
	if err != nil {
		app.logger.Println("Failed to authenticate user with Supabase")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	userJson, err := json.Marshal(resp.User)
	if err != nil {
		// handle error
		app.logger.Println("Failed to marshal user JSON")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal user JSON"})
		return
	}
	var authResponse model.AuthResponse

	if err := json.Unmarshal(userJson, &authResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}
	// Set the access token in the response header
	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", authResponse.AccessToken))
	ctx.JSON(http.StatusOK, resp.User)
}

func (app *application) logoutHandler(ctx *gin.Context) {
	ctx.Header("Authorization", "")
	ctx.SetCookie("access_token", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (app *application) registerHandler(ctx *gin.Context) {
	var registerRequest model.RegisterRequest
	err := ctx.ShouldBindJSON(&registerRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := supabase.NewClient(app.supabaseConfig.ApiUrl, app.supabaseConfig.ApiKey, nil)
	if err != nil {
		app.logger.Println("Failed to create Supabase client")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Supabase client"})
		return
	}
	req := types.SignupRequest{
		Email:    registerRequest.Email,
		Password: registerRequest.Password,
	}
	resp, err := client.Auth.Signup(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userJson, err := json.Marshal(resp.User)
	if err != nil {
		// handle error
		app.logger.Println("Failed to marshal user JSON")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal user JSON"})
		return
	}
	registerResponse := model.RegisterResponse{}
	if err := json.Unmarshal(userJson, &registerResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}
	client.From("Profiles").Insert(registerResponse, true, "", "", "")

	ctx.JSON(http.StatusOK, registerResponse)
}

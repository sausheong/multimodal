package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// call the Vertex AI Imagen API with an image
func imagen(prompt string, imageBase64 string) ([]string, error) {
	requestURL := "https://us-central1-aiplatform.googleapis.com/v1/projects/" +
		os.Getenv("GOOGLE_PROJECT_ID") +
		"/locations/us-central1/publishers/google/models/imagetext:predict"

	// create the request
	var request ImagenRequest
	if prompt == "" {
		request = ImagenRequest{
			Instances: []Instance{
				{
					Image: Image{
						BytesBase64Encoded: imageBase64,
					},
				},
			},
			Parameters: Parameters{
				SampleCount: 1,
			},
		}

	} else {
		request = ImagenRequest{
			Instances: []Instance{
				{
					Prompt: prompt,
					Image: Image{
						BytesBase64Encoded: imageBase64,
					},
				},
			},
			Parameters: Parameters{
				SampleCount: 1,
			},
		}
	}
	// marshal the request into JSON
	requestJson, err := json.Marshal(request)
	if err != nil {
		return []string{}, err
	}
	// create a HTTP request with the JSON
	req, err := http.NewRequest(http.MethodPost, requestURL,
		bytes.NewReader(requestJson))
	if err != nil {
		return []string{}, err
	}
	// get the refresh token
	refresh, err := getRefreshToken("./application_default_credentials.json")
	if err != nil {
		return []string{}, err
	}
	// get the access token using the refresh token
	access, err := getAccessToken(refresh)
	if err != nil {
		return []string{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	// set the authorization using the access token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access.AccessToken))
	// create a client to send the HTTP request
	client := http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}
	// read the JSON results from the HTTP response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []string{}, err
	}
	// unmarshal the JSON results into structs
	response := ImagenResponse{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return []string{}, err
	}
	// just return the predictions
	return response.Predictions, nil
}

// get an access token from Google OAuth2 to access APIs
func getAccessToken(refreshToken RefreshToken) (AccessToken, error) {
	requestURL := "https://oauth2.googleapis.com/token"
	// marshal the refresh token into JSON
	requestJson, err := json.Marshal(refreshToken)
	if err != nil {
		return AccessToken{}, err
	}
	// create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, requestURL,
		bytes.NewReader(requestJson))
	if err != nil {
		return AccessToken{}, err
	}
	//set the content type to JSON
	req.Header.Set("Content-Type", "application/json")
	// create the client to send the HTTP request
	client := http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return AccessToken{}, err
	}
	// read the results from the HTTP response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return AccessToken{}, err
	}
	// unmarshal the returned access token
	accessToken := AccessToken{}
	err = json.Unmarshal(resBody, &accessToken)
	if err != nil {
		return AccessToken{}, err
	}
	return accessToken, nil
}

// read the refresh token data from the
// application_default_credentials.json file
func getRefreshToken(path string) (RefreshToken, error) {
	// open the file
	file, err := os.Open(path)
	if err != nil {
		return RefreshToken{}, err
	}
	// read the file
	data, err := io.ReadAll(file)
	if err != nil {
		return RefreshToken{}, err
	}
	// unmarshal the contents of the file to structs
	t := RefreshToken{}
	err = json.Unmarshal(data, &t)
	if err != nil {
		return RefreshToken{}, err
	}
	t.GrantType = "refresh_token"
	return t, nil
}

// structs

type RefreshToken struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
	Type         string `json:"type"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token"`
}

// structs for sending request to API
type ImagenRequest struct {
	Instances  []Instance `json:"instances"`
	Parameters Parameters `json:"parameters"`
}

type Instance struct {
	Prompt string `json:"prompt,omitempty"`
	Image  Image  `json:"image"`
}
type Image struct {
	BytesBase64Encoded string `json:"bytesBase64Encoded"`
}
type Parameters struct {
	SampleCount int `json:"sampleCount"`
}

// struct for parsing response
type ImagenResponse struct {
	Predictions      []string `json:"predictions"`
	DeployedModelID  string   `json:"deployedModelId"`
	Model            string   `json:"model"`
	ModelDisplayName string   `json:"modelDisplayName"`
	ModelVersionID   string   `json:"modelVersionId"`
}

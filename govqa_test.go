package main

import (
	"fmt"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	token, err := getRefreshToken("./gcloud/application_default_credentials.json")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%+v\n", token)
	}
}

func TestAccessToken(t *testing.T) {
	refresh, err := getRefreshToken("./gcloud/application_default_credentials.json")
	if err != nil {
		t.Error(err)
	}
	access, err := getAccessToken(refresh)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%+v\n", access)
	}
}

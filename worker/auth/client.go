package auth

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pubudu2003060/go-proxy-prototype/worker/models"
)

type AuthClient struct {
	captainURL string
}

func NewAuthClient(captainURL string) *AuthClient {
	return &AuthClient{
		captainURL: captainURL,
	}
}

func (c *AuthClient) Authenticate(username, password string) (*models.AuthResponse, error) {
	reqBody := models.AuthRequest{
		Username: username,
		Password: password,
	}

	jsonData,err := json.Marshal(reqBody)
	if err != nil {
		return nil,err
	}

	resp,err := http.Post(c.captainURL+"api/v1/auth","application/json",bytes.NewBuffer(jsonData))
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp);err != nil {
		return nil,err
	}

	return &authResp,nil
}
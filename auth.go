package goprsc

import "net/http"

const authURL = "auth"
const loginURL = authURL + "/signin"
const logoutURL = authURL + "/signout"

// AuthService handles communication with the authentication APIs in the Postfix Rest Server.
type AuthService service

// LoginRequest represents a request for logging in
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LogoutRequest representes a request for logging out
type LogoutRequest struct {
	Login        string `json:"login"`
	RefreshToken string `json:"refreshToken"`
}

// AuthResponse represents a successfull authentication response
type AuthResponse struct {
	AuthToken    string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// Login makes a post request to the API for logging in
func (s *AuthService) Login(login, password string) (*AuthResponse, error) {
	req := &LoginRequest{
		Login:    login,
		Password: password,
	}
	request, err := s.client.NewRequest(http.MethodPost, loginURL, req)
	if err != nil {
		return nil, err
	}
	res := &AuthResponse{}
	_, err = s.client.Do(request, res)
	return res, err
}

// Logout makes a post request to the API for logging out
func (s *AuthService) Logout(login, refreshToken string) error {
	req := &LogoutRequest{
		Login:        login,
		RefreshToken: refreshToken,
	}
	request, err := s.client.NewRequest(http.MethodPost, logoutURL, req)
	if err != nil {
		return err
	}
	_, err = s.client.Do(request, nil)
	return err
}

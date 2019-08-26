package client

import "net/http"

// AuthenticateUserRequest represents the body of a request to POST /auth
type AuthenticateUserRequest struct {
	Username string
	Password string
}

// AuthenticateUserResponse represents the body of a response for a request to POST /auth
type AuthenticateUserResponse struct {
	Jwt string
}

func (n *portainerClientImp) AuthenticateUser() (token string, err error) {
	reqBody := AuthenticateUserRequest{
		Username: n.user,
		Password: n.password,
	}

	respBody := AuthenticateUserResponse{}

	err = n.doJSON("auth", http.MethodPost, http.Header{}, &reqBody, &respBody)
	if err != nil {
		return
	}

	token = respBody.Jwt

	return
}

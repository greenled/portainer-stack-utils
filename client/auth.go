package client

import "net/http"

// AuthenticateUserOptions represents options passed to PortainerClient.AuthenticateUser()
type AuthenticateUserOptions struct {
	Username string
	Password string
}

// AuthenticateUserRequest represents the body of a request to POST /auth
type AuthenticateUserRequest struct {
	Username string
	Password string
}

// AuthenticateUserResponse represents the body of a response for a request to POST /auth
type AuthenticateUserResponse struct {
	Jwt string
}

func (n *portainerClientImp) AuthenticateUser(options AuthenticateUserOptions) (token string, err error) {
	reqBody := AuthenticateUserRequest{
		Username: options.Username,
		Password: options.Password,
	}

	respBody := AuthenticateUserResponse{}

	err = n.doJSON("auth", http.MethodPost, http.Header{}, &reqBody, &respBody)
	if err != nil {
		return
	}

	token = respBody.Jwt

	return
}

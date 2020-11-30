package main

type SessionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionRequests []SessionRequest

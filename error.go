package main

import "errors"

var (
	BadAuthError = errors.New("Bad authentication")
	RateLimitError = errors.New("Please wait before doing this again")
)

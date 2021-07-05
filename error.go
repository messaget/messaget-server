package main

import "errors"

var (
	BadAuthError     = errors.New("Bad authentication")
	RateLimitError   = errors.New("Please wait before doing this again")
	NoNamespaceError = errors.New("You must provide a namespace")
	NamespaceTooLong = errors.New("That namespace is too long")
	FailedIntent     = errors.New("Failed to parse intent")
)

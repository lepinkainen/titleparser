package handler

import "github.com/lepinkainen/titleparser/common"

// Common configuration for all handlers
// Re-export from common package for backward compatibility

var (
	// UserAgent string to use when connecting to servers
	UserAgent = common.UserAgent
	// AcceptLanguage header
	AcceptLanguage = common.AcceptLanguage
	// Accept header
	Accept = common.Accept
)

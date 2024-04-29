package json

import (
	"net/http"
)

type ErrorCode int

const (
	// Code Error
	ErrInvalidRequest         ErrorCode = 1000 // E1000
	ErrUnauthorizedClient     ErrorCode = 1001 // E1001
	ErrAccessDenied           ErrorCode = 1002 // E1002
	ErrServerError            ErrorCode = 1003 // E1003
	ErrTemporarilyUnavailable ErrorCode = 1004 // E1004
	ErrInvalidClient          ErrorCode = 1005 // E1005
	ErrInvalidGrant           ErrorCode = 1006 // E1006
	ErrUnsupportedGrantType   ErrorCode = 1007 // E1007
)

var ErrorCodeDescriptions = map[ErrorCode]string{
	ErrInvalidRequest:         "The request is invalid",
	ErrUnauthorizedClient:     "The client is not authorized to access the requested resource",
	ErrAccessDenied:           "Access to the requested resource is denied",
	ErrServerError:            "The server encountered an internal error while processing the request",
	ErrTemporarilyUnavailable: "The requested resource is temporarily unavailable",
	ErrInvalidClient:          "The client is invalid or not recognized",
	ErrInvalidGrant:           "The grant or authorization code is invalid or expired",
	ErrUnsupportedGrantType:   "The grant type requested is not supported by the server",
}

type StatusCode int

const (
	// Status Code
	StatusOK                  StatusCode = http.StatusOK                  // 200 OK
	StatusCreated             StatusCode = http.StatusCreated             // 201 Created
	StatusAccepted            StatusCode = http.StatusAccepted            // 202 Accepted
	StatusNoContent           StatusCode = http.StatusNoContent           // 204 No Content
	StatusMovedPermanently    StatusCode = http.StatusMovedPermanently    // 301 Moved Permanently
	StatusFound               StatusCode = http.StatusFound               // 302 Found
	StatusSeeOther            StatusCode = http.StatusSeeOther            // 303 See Other
	StatusNotModified         StatusCode = http.StatusNotModified         // 304 Not Modified
	StatusTemporaryRedirect   StatusCode = http.StatusTemporaryRedirect   // 307 Temporary Redirect
	StatusBadRequest          StatusCode = http.StatusBadRequest          // 400 Bad Request
	StatusUnauthorized        StatusCode = http.StatusUnauthorized        // 401 Unauthorized
	StatusForbidden           StatusCode = http.StatusForbidden           // 403 Forbidden
	StatusNotFound            StatusCode = http.StatusNotFound            // 404 Not Found
	StatusMethodNotAllowed    StatusCode = http.StatusMethodNotAllowed    // 405 Method Not Allowed
	StatusInternalServerError StatusCode = http.StatusInternalServerError // 500 Internal Server Error
	StatusServiceUnavailable  StatusCode = http.StatusServiceUnavailable  // 503 Service Unavailable
)

// HTTPCodeDescriptions maps HTTP status codes to brief descriptions.
var HTTPCodeDescriptions = map[StatusCode]string{
	StatusOK:                  "Request successful",
	StatusCreated:             "Resource created",
	StatusAccepted:            "Request accepted",
	StatusNoContent:           "No content to return",
	StatusMovedPermanently:    "Resource moved permanently",
	StatusFound:               "Resource temporarily moved",
	StatusSeeOther:            "Response at a different URL",
	StatusNotModified:         "Resource not modified",
	StatusTemporaryRedirect:   "Temporarily redirected",
	StatusBadRequest:          "Invalid request",
	StatusUnauthorized:        "Authentication required",
	StatusForbidden:           "Request forbidden",
	StatusNotFound:            "Resource not found",
	StatusMethodNotAllowed:    "Method not allowed",
	StatusInternalServerError: "Server error",
	StatusServiceUnavailable:  "Service unavailable",
}

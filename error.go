package jsonapi

import "github.com/nats-io/nuid"

// Error represents the error object that is specified on the
// http://jsonapi.org
type Error struct {
	ID     string       `json:"id"`
	Links  *ErrorLinks  `json:"links,omitempty"`
	Status string       `json:"status,omitempty"`
	Code   string       `json:"code,omitempty"`
	Title  string       `json:"title,omitempty"`
	Detail string       `json:"detail,omitempty"`
	Source *ErrorSource `json:"source,omitempty"`
	Meta   interface{}  `json:"meta,omitempty"`
}

// ErrorLinks is used to provide an About URL that leads to
// further details about the particular occurrence of the problem.
// for more information see http://jsonapi.org/format/#error-objects
type ErrorLinks struct {
	About string `json:"about,omitempty"`
}

// ErrorSource is used to provide references to the source of an error.
// The Pointer is a JSON Pointer to the associated entity in the request
// document.
// The Paramter is a string indicating which query parameter caused the error.
// for more information see http://jsonapi.org/format/#error-objects
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// GetID satisfies the ResourceGetIdentifier interface. If the given Error.ID is
// set it returns the value for the Error.ID. If it is not set does set it
// returns nuid.Next() which is a randomly generated string.
func (e Error) GetID() string {
	if e.ID != "" {
		return e.ID
	}
	return nuid.Next()
}

// GetName ...
func (e Error) GetName() string {
	return "error"
}

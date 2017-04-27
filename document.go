package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Document represents the top level JSON Object for every JSON API request
// and response that contains the data as specified on the
// http://jsonapi.org/format/#document-top-level
type Document struct {
	Links    *Links                 `json:"links,omitempty"`
	Included []Data                 `json:"included,omitempty"`
	Data     *DataContainer         `json:"data,omitempty"`
	Errors   []Error                `json:"errors,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

// DataContainer is a custom type that is used to represent the data in the
// Document type. The main use for this type is to hold the slice or non-slice
// values of the data when unmarshaling or marshaling.
type DataContainer struct {
	DataObject *Data
	DataArray  []Data
}

// UnmarshalJSON satisfies the Unmarshaler interface so that we can implement
// a custom UnmarshalJSON function that checks if the object being unmarshaled
// is a single object or a collection.
func (c *DataContainer) UnmarshalJSON(payload []byte) error {
	// payload is an object
	if bytes.HasPrefix(payload, []byte("{")) {
		return json.Unmarshal(payload, &c.DataObject)
	}

	// payload is an array
	if bytes.HasPrefix(payload, []byte("[")) {
		return json.Unmarshal(payload, &c.DataArray)
	}

	return errors.New("Invalid json for data array/object")
}

// MarshalJSON satisfies the Marshaler interface so that we can implement a
// custom MarshalJSON function that checks if the object being marshaled is
// a single object or a collection.
func (c *DataContainer) MarshalJSON() ([]byte, error) {
	if c.DataArray != nil {
		return json.Marshal(&c.DataArray)
	}
	return json.Marshal(&c.DataObject)
}

// Links is general links struct for top level and relationships
type Links struct {
	Self     string `json:"self,omitempty"`
	Related  string `json:"related,omitempty"`
	First    string `json:"first,omitempty"`
	Previous string `json:"prev,omitempty"`
	Next     string `json:"next,omitempty"`
	Last     string `json:"last,omitempty"`
}

// Data for top level and included data
type Data struct {
	Type          string                  `json:"type"`
	ID            string                  `json:"id,omitempty"`
	Attributes    json.RawMessage         `json:"attributes"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
	Links         *Links                  `json:"links,omitempty"`
}

// Relationship contains reference IDs to the related structs
type Relationship struct {
	Links *Links                     `json:"links,omitempty"`
	Data  *RelationshipDataContainer `json:"data,omitempty"`
	Meta  map[string]interface{}     `json:"meta,omitempty"`
}

// RelationshipDataContainer is needed to either keep relationship "data"
// contents as array or object.
type RelationshipDataContainer struct {
	DataObject *RelationshipData
	DataArray  []RelationshipData
}

// UnmarshalJSON implements Unmarshaler and also detects array/object type
func (c *RelationshipDataContainer) UnmarshalJSON(payload []byte) error {
	if bytes.HasPrefix(payload, []byte("{")) {
		// payload is an object
		return json.Unmarshal(payload, &c.DataObject)
	}

	if bytes.HasPrefix(payload, []byte("[")) {
		// payload is an array
		return json.Unmarshal(payload, &c.DataArray)
	}

	return errors.New("Invalid JSON for relationship data array/object")
}

// MarshalJSON either Marshals an array or object of relationship data
func (c *RelationshipDataContainer) MarshalJSON() ([]byte, error) {
	if c.DataArray != nil {
		return json.Marshal(c.DataArray)
	}
	return json.Marshal(c.DataObject)
}

// RelationshipData represents one specific reference ID
type RelationshipData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

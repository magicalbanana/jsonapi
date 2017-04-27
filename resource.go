package jsonapi

// ResourceGetIdentifier ...
type ResourceGetIdentifier interface {
	GetID() string
}

// ResourceSetIdentifier ...
type ResourceSetIdentifier interface {
	SetID(string) error
}

// ResourceTyper ...
type ResourceTyper interface {
	GetName() string
}

// UnmarshalIdentifier ...
type UnmarshalIdentifier interface {
	SetID(string) error
}

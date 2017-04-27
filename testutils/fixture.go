// Package testutils contains the utilities that are used for testing
package testutils

import "time"

// FisherBoy does not satisfy the MarshalIndetifer interfaces
type FisherBoy struct {
	ID   string `json:"-"`
	Fish string `json:"fish"`
}

// Manbearpig ...
type Manbearpig struct {
	ID      string    `json:"-"`
	Moo     string    `json:"moo"`
	Zoo     string    `json:"zoo"`
	FooTime time.Time `json:"foo_time"`
	ZooTime time.Time `json:"zoo_time"`
}

// GetID ...
func (m Manbearpig) GetID() string {
	return m.ID
}

// SetID ...
func (m *Manbearpig) SetID(ID string) error {
	m.ID = ID

	return nil
}

// GetName ...
func (m Manbearpig) GetName() string {
	return "manbearpigs"
}

// NestedManbearpigs ...
type NestedManbearpigs struct {
	ID   string       `json:"-"`
	MBPS []Manbearpig `json:"mbps"`
}

// GetID ...
func (n NestedManbearpigs) GetID() string {
	return n.ID
}

// SetID ...
func (n *NestedManbearpigs) SetID(ID string) error {
	n.ID = ID
	return nil
}

// GetName ...
func (n NestedManbearpigs) GetName() string {
	return "nestedManbearpigs"
}

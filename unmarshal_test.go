package jsonapi

import (
	"testing"
	"time"

	th "git.enova.com/kgan/jsonapi/testutils"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	ti, _ := time.Parse(time.RFC3339, "2014-11-10T16:30:48.823Z")
	mbp1 := th.Manbearpig{ID: "1", Moo: "Manbearpig", Zoo: "Pigglywiggly", FooTime: ti}
	mbp2 := th.Manbearpig{ID: "2", Moo: "Man bear pig", Zoo: "Piggly wiggly", FooTime: ti, ZooTime: ti}
	single := []byte(`
		{
			"data": {
				"id": "1",
				"type": "manbearpigs",
				"attributes": {
					"moo": "Manbearpig",
					"zoo": "Pigglywiggly",
					"foo_time": "2014-11-10T16:30:48.823Z"
				}
			}
		}
		`)

	collection := []byte(`
				{
					"data": [
					{
						"id": "1",
						"type": "manbearpigs",
						"attributes": {
							"moo": "Manbearpig",
							"zoo": "Pigglywiggly",
							"foo_time": "2014-11-10T16:30:48.823Z"
						}
					},
					{
						"id": "2",
						"type": "manbearpigs",
						"attributes": {
							"moo": "Man bear pig",
							"zoo": "Piggly wiggly",
							"foo_time": "2014-11-10T16:30:48.823Z",
							"zoo_time": "2014-11-10T16:30:48.823Z"
						}
					}
					]
				}
	`)

	nestedCollection := []byte(`
			{
				"data": {
					"type": "nestedManbearpigs",
					"id": "1",
					"attributes": {
						"mbps": [
						    {
							    "moo": "Manbearpig",
							    "zoo": "Pigglywiggly"
						    }
						]
					}
				}
			}
	`)

	tests := []struct {
		desc      string
		assertion func(*testing.T, string)
	}{
		{
			desc: "unmarshals a single payload to a struct",
			assertion: func(t *testing.T, desc string) {
				mbp := th.Manbearpig{}
				err := Unmarshal(single, &mbp)
				assert.NoError(t, err, desc)
				assert.Equal(t, mbp, mbp1)
			},
		},
		{
			desc: "unmarshals a collection payload into a slice of the given type",
			assertion: func(t *testing.T, desc string) {
				var mpbs []th.Manbearpig
				err := Unmarshal(collection, &mpbs)
				assert.NoError(t, err, desc)
				assert.Equal(t, mpbs, []th.Manbearpig{mbp1, mbp2})
			},
		},
		{
			desc: "unmarshalls an array value for an attribute",
			assertion: func(t *testing.T, desc string) {
				var nmpbs th.NestedManbearpigs
				expected := th.NestedManbearpigs{ID: "1", MBPS: []th.Manbearpig{th.Manbearpig{Moo: "Manbearpig", Zoo: "Pigglywiggly"}}}
				err := Unmarshal(nestedCollection, &nmpbs)
				assert.NoError(t, err, desc)
				assert.Equal(t, nmpbs, expected, desc)
			},
		},
		{
			desc: "returns error if struct does not satisfy MarshalIdentifer interface",
			assertion: func(t *testing.T, desc string) {
				loo := struct {
					Foo string `json:"foo"`
				}{}
				err := Unmarshal(single, &loo)
				assert.Error(t, err, desc)
			},
		},
		{
			desc: "returns error if struct is nil",
			assertion: func(t *testing.T, desc string) {
				err := Unmarshal(single, nil)
				assert.Error(t, err, desc)
			},
		},
		{
			desc: "returns error if passed a non pointer struct",
			assertion: func(t *testing.T, desc string) {
				mbp := &th.Manbearpig{}
				err := Unmarshal(single, *mbp)
				assert.Error(t, err, desc)
			},
		},
		{
			desc: "returns error if JSON payload does not contain 'attributes'",
			assertion: func(t *testing.T, desc string) {
				invalid := []byte(`{"life": {"type":"manbearpigs"}}`)
				err := Unmarshal(invalid, &th.Manbearpig{})
				assert.Error(t, err, desc)
				assert.Equal(t, err.Error(), `Source JSON is empty and does not satisfy the JSONAPI specification!`, desc)
			},
		},
		{
			desc: "returns error if JSON payload is malformed",
			assertion: func(t *testing.T, desc string) {
				// it's malformed because the closing brace
				// for the top level document is missing.
				invalid := []byte(`{"life": {"type":"manbearpig"}`)
				err := Unmarshal(invalid, &th.Manbearpig{})
				assert.Error(t, err, desc)
				// the error that's matched CAN change
				// depending on the malformd JSON we are
				// unmarshalling.
				assert.Equal(t, err.Error(), `unexpected end of JSON input`, desc)
			},
		},
		{
			desc: "returns error if JSON payload is invalid",
			assertion: func(t *testing.T, desc string) {
				// this is invalida because the data attribute
				// is not a nested object but rather an int
				invalid := []byte(`{"data": 1`)
				err := Unmarshal(invalid, &th.Manbearpig{})
				assert.Error(t, err, desc)
				// the error that's matched CAN change
				// depending on the malformd JSON we are
				// unmarshalling.
				assert.Equal(t, err.Error(), `unexpected end of JSON input`, desc)
			},
		},
		{
			desc: "returns error if JSON payload is collection but target is not",
			assertion: func(t *testing.T, desc string) {
				err := Unmarshal(collection, &th.Manbearpig{})
				assert.Error(t, err, desc)
				assert.Equal(t, err.Error(), `Cannot unmarshal array to struct target testutils.Manbearpig`, desc)
			},
		},
		{
			desc: "returns error if JSON payload is collection but target is not",
			assertion: func(t *testing.T, desc string) {
				err := Unmarshal(collection, &th.FisherBoy{})
				assert.Error(t, err, desc)
				assert.Equal(t, err.Error(), `Cannot unmarshal array to struct target testutils.FisherBoy`, desc)
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}

}

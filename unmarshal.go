package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrInterface occurs if the target struct did not implement the
	// UnmarshalIdentifier interface
	ErrInterface = errors.New("target must implement UnmarshalIdentifier interface")
	// ErrNoType occurs if the JSON payload did not have a type field
	ErrNoType = errors.New("invalid record, no type was specified")
)

// Unmarshal reads a jsonapi compatible JSON as []byte
// target must at least implement the `UnmarshalIdentifier` interface.
func Unmarshal(data []byte, target interface{}) error {
	if target == nil {
		return errors.New("target must not be nil")
	}

	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		return errors.New("target must be a ptr")
	}

	doc := &Document{}
	err := json.Unmarshal(data, doc)
	if err != nil {
		return err
	}

	if doc.Data == nil {
		return errors.New(`Source JSON is empty and does not satisfy the JSONAPI specification!`)
	}

	if doc.Data.DataObject != nil {
		return setDataIntoTarget(doc.Data.DataObject, target)
	}

	if doc.Data.DataArray != nil {
		targetSlice := reflect.TypeOf(target).Elem()
		if targetSlice.Kind() != reflect.Slice {
			return fmt.Errorf("Cannot unmarshal array to struct target %s", targetSlice)
		}
		targetType := targetSlice.Elem()
		targetPointer := reflect.ValueOf(target)
		targetValue := targetPointer.Elem()

		for _, record := range doc.Data.DataArray {
			// check if there already is an entry with the same ID
			// in target slice, otherwise create a new target and
			// append
			var targetRecord, emptyValue reflect.Value
			for i := 0; i < targetValue.Len(); i++ {
				// type assert the targetValue so that we can
				// leverage the GetID() func that MUST be
				// implemented by the targetValue.
				targetTyped, ok := targetValue.Index(i).Interface().(ResourceGetIdentifier)
				if !ok {
					// TODO: Return better error messages
					return errors.New("existing structs must satisfy interface ResourceGetIdentifier")
				}
				if record.ID == targetTyped.GetID() {
					targetRecord = targetValue.Index(i).Addr()
					break
				}
			}

			if targetRecord == emptyValue || targetRecord.IsNil() {
				targetRecord = reflect.New(targetType)
				err := setDataIntoTarget(&record, targetRecord.Interface())
				if err != nil {
					return err
				}
				targetValue = reflect.Append(targetValue, targetRecord.Elem())
			} else {
				err := setDataIntoTarget(&record, targetRecord.Interface())
				if err != nil {
					return err
				}
			}
		}

		targetPointer.Elem().Set(targetValue)
	}

	return nil
}

func setDataIntoTarget(data *Data, target interface{}) error {
	castedTarget, ok := target.(UnmarshalIdentifier)
	if !ok {
		return ErrInterface
	}

	if data.Type == "" {
		return ErrNoType
	}

	err := checkType(data.Type, castedTarget)
	if err != nil {
		return err
	}

	if data.Attributes != nil {
		err = json.Unmarshal(data.Attributes, castedTarget)
		if err != nil {
			return err
		}
	}
	return castedTarget.SetID(data.ID)
}

func checkType(incomingType string, target UnmarshalIdentifier) error {
	actualType, err := getStructType(target)
	if err != nil {
		return err
	}

	if incomingType != actualType {
		return fmt.Errorf("Type %s in JSON does not match target struct type %s", incomingType, actualType)
	}

	return nil
}

package gslutils

import (
	"reflect"

	"github.com/pkg/errors"
)

var ErrNotConvertible = errors.New("types are not convertible")

// CompareInterfaceValues will compares |a| and |b| as type T.
// It returns true if |a| and |b| are of type T and are equal in value,
// or if |a| and |b| can be converted to T and are equal in value.
// If |a| and |b| are not equal but can be converted to T, then it returns false.
// CompareInterfaceValues *only* returns non-nil error if any of the values cannot be converted into T.
func CompareInterfaceValues[T comparable](a, b interface{}) (bool, error) {
	var t T
	typeTarget := reflect.TypeOf(t)

	assertedA, err := InterfaceTo[T](a)
	if err != nil {
		return false, errors.Wrapf(err, "cannot convert a to type %s", typeTarget.String())
	}

	assertedB, err := InterfaceTo[T](b)
	if err != nil {
		return false, errors.Wrapf(err, "cannot convert b to type %s", typeTarget.String())
	}

	return assertedA == assertedB, nil
}

// InterfaceTo converts v from interface{} to T.
// It returns a zeroed T and an error if not convertible.
func InterfaceTo[T any](v interface{}) (T, error) {
	// If we can directly assert the type, then return the asserted value
	t, ok := v.(T)
	if ok {
		return t, nil
	}

	// Otherwise use reflect package to convert v to T
	valueType := reflect.TypeOf(v)
	targetType := reflect.TypeOf(t)

	if valueType.ConvertibleTo(targetType) {
		valueInTypeT, ok := reflect.ValueOf(v).Convert(targetType).Interface().(T)
		if ok {
			return valueInTypeT, nil
		}
	}

	return t, errors.Wrapf(ErrNotConvertible, "cannot convert %s to %s", valueType.String(), targetType.String())
}

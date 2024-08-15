package assert

import (
	"errors"
	"reflect"
	"slices"
	"testing"
)

func Nil(t *testing.T, value any) bool {
	if !isNil(value) {
		t.Errorf("Expected nil. Got: %+v", value)
		return false
	}

	return true
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()

	isNilableKind := slices.Contains([]reflect.Kind{
		reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer}, kind)
	if isNilableKind && value.IsNil() {
		return true
	}

	return false
}

func Equal(t *testing.T, expected, got any) bool {
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected: %+v. Got: %+v", expected, got)
		return false
	}

	return true
}

func NotEqual(t *testing.T, expected, got any) bool {
	return !Equal(t, expected, got)
}

func NoError(t *testing.T, err error) bool {
	if err != nil {
		t.Errorf("Received unexpected error:\n%+v", err)
		return false
	}

	return true
}

func Error(t *testing.T, err error) bool {
	if err == nil {
		t.Error("An error was expected")
		return false
	}

	return true
}

func ErrorIs(t *testing.T, err error, target error) bool {
	if err == nil {
		t.Error("An error was expected")
		return false
	}
	if !errors.Is(err, target) {
		t.Errorf("Expected error: %+v. Got: %+v", err, target)
		return false
	}

	return true
}

func NotEmpty(t *testing.T, value any) bool {
	if isEmpty(value) {
		t.Error("Should not be empty")
		return false
	}
	return true
}

// isEmpty gets whether the specified object is considered empty or not.
func isEmpty(object interface{}) bool {

	// get nil case out of the way
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	// pointers are empty if nil or if the value they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return isEmpty(deref)
	// for all other types, compare against the zero value
	// array types are empty when they match their zero-initialized state
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

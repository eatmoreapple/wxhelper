package structcopy

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrorInvalidStruct is returned when the src or dst is not a struct.
	ErrorInvalidStruct = errors.New("struct required")

	// ErrorNotPointer is returned when the dst is not a pointer.
	ErrorNotPointer = errors.New("pointer required")
)

type StructCopier[T any] struct{}

// Copy copies the src to dst.
func (c *StructCopier[T]) Copy(from any) (result T, err error) {
	src := reflect.ValueOf(from)
	if src.Kind() == reflect.Ptr {
		src = reflect.Indirect(src)
	}
	if src.Kind() != reflect.Struct {
		return result, ErrorInvalidStruct
	}

	dst := reflect.ValueOf(result)

	if dst.Kind() == reflect.Ptr {
		result = reflect.New(dst.Type().Elem()).Interface().(T)
		dst = reflect.ValueOf(result)
	} else {
		dst = reflect.ValueOf(&result)
	}

	if dst.Kind() != reflect.Ptr {
		return result, ErrorNotPointer
	}
	dst = reflect.Indirect(dst)
	if dst.Kind() != reflect.Struct {
		return result, ErrorInvalidStruct
	}
	return result, c.copyField(src, dst)
}

// copyField copies the fields from src to dst.
func (c *StructCopier[T]) copyField(src reflect.Value, dst reflect.Value) error {
	dstType := dst.Type()
	for i := 0; i < dstType.NumField(); i++ {
		field := dst.Field(i)
		if !field.CanSet() {
			continue
		}
		value := src.FieldByName(dstType.Field(i).Name)
		if !value.IsValid() {
			continue
		}
		if value.Type() != field.Type() {
			if !value.Type().ConvertibleTo(field.Type()) {
				return fmt.Errorf("cannot convert %s to %s", value.Type(), field.Type())
			}
			value = value.Convert(field.Type())
		}
		field.Set(value)
	}
	return nil
}

// Copy copies the src to dst.
func Copy[T any](src any) (T, error) {
	copier := StructCopier[T]{}
	return copier.Copy(src)
}

// CopySlice copies the src to dst.
func CopySlice[T, E any](from []E) ([]T, error) {
	length := len(from)
	result := make([]T, length)
	copier := StructCopier[T]{}
	for i := 0; i < length; i++ {
		item, err := copier.Copy(from[i])
		if err != nil {
			return nil, err
		}
		result[i] = item
	}
	return result, nil
}

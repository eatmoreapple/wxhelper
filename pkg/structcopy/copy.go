package structcopy

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrorInvalidStruct is returned when the src or dst is not a struct.
	ErrorInvalidStruct = errors.New("struct required")

	// ErrorNotSlice is returned when the src is not a slice.
	ErrorNotSlice = errors.New("slice required")

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

func Copy[T any](src any) (T, error) {
	copier := StructCopier[T]{}
	return copier.Copy(src)
}

func CopySlice[T any](from any) ([]T, error) {
	srv := reflect.ValueOf(from)
	if srv.Kind() != reflect.Slice {
		return nil, ErrorNotSlice
	}
	dst := make([]T, srv.Len())
	for i := 0; i < srv.Len(); i++ {
		item := srv.Index(i)
		dstItem, err := Copy[T](item.Interface())
		if err != nil {
			return nil, err
		}
		dst[i] = dstItem
	}
	return dst, nil
}

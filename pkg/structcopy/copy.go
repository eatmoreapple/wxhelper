package structcopy

import (
	"errors"
	"reflect"
)

var (
	ErrNotStruct  = errors.New("struct required")
	ErrorNotSlice = errors.New("slice required")
)

type StructCopier[T any] struct{}

func (c *StructCopier[T]) Copy(from any) (result T, err error) {
	src := reflect.ValueOf(from)
	if src.Kind() == reflect.Ptr {
		src = reflect.Indirect(src)
	}
	if src.Kind() != reflect.Struct {
		return result, ErrNotStruct
	}

	dst := reflect.ValueOf(result)

	if dst.Kind() == reflect.Ptr {
		result = reflect.New(dst.Type().Elem()).Interface().(T)
		dst = reflect.ValueOf(result)
	} else {
		dst = reflect.ValueOf(&result)
	}

	if dst.Kind() != reflect.Ptr {
		return result, errors.New("dst must be a pointer")
	}
	dst = reflect.Indirect(dst)
	if dst.Kind() != reflect.Struct {
		return result, ErrNotStruct
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
				return errors.New("type mismatch")
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

func CopySlice[T any](src any) ([]T, error) {
	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() != reflect.Slice {
		return nil, ErrorNotSlice
	}
	dst := make([]T, srcValue.Len())
	for i := 0; i < srcValue.Len(); i++ {
		item := srcValue.Index(i)
		dstItem, err := Copy[T](item.Interface())
		if err != nil {
			return nil, err
		}
		dst[i] = dstItem
	}
	return dst, nil
}

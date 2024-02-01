package pkg

type StructCopier[T any] struct{}

func (c *StructCopier[T]) Copy(from any) (T, error) {
	panic("implement me")
}

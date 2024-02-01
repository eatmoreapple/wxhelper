package structcopy

import "testing"

func TestCopy(t *testing.T) {
	type A struct {
		Name string
		Age  int
	}
	type B struct {
		Name string
		Age  int64
	}
	a := A{
		Name: "test",
		Age:  10,
	}
	b, err := Copy[*B](a)
	if err != nil {
		t.Fatal(err)
	}
	if b.Name != a.Name {
		t.Fatalf("expected %s, got %s", a.Name, b.Name)
	}
	if b.Age != int64(a.Age) {
		t.Fatalf("expected %d, got %d", a.Age, b.Age)
	}
}

func TestCopySlice(t *testing.T) {
	type A struct {
		Name string
		Age  int
	}
	type B struct {
		Name string
		Age  int64
	}
	a := []A{
		{
			Name: "test",
			Age:  10,
		},
		{
			Name: "test2",
			Age:  20,
		},
	}
	b, err := CopySlice[*B](a)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != len(a) {
		t.Fatalf("expected %d, got %d", len(a), len(b))
	}
	for i := range a {
		if b[i].Name != a[i].Name {
			t.Fatalf("expected %s, got %s", a[i].Name, b[i].Name)
		}
		if b[i].Age != int64(a[i].Age) {
			t.Fatalf("expected %d, got %d", a[i].Age, b[i].Age)
		}
	}
}

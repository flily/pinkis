package meta

import (
	"testing"
)

func TestValueEqualSimpleValue(t *testing.T) {
	var a, b interface{}

	a = 1
	b = 1
	if !Equal(a, b) {
		t.Errorf("unexpected result: %v", a)
	}

	a = 1
	b = 2
	if Equal(a, b) {
		t.Errorf("unexpected result: %v", a)
	}

	a = nil
	b = nil
	if !Equal(a, b) {
		t.Errorf("unexpected result: %v", a)
	}

	a = nil
	b = 1
	if Equal(a, b) {
		t.Errorf("unexpected result: %v", a)
	}

	a = 1
	b = nil
	if Equal(a, b) {
		t.Errorf("unexpected result: %v", a)
	}

	a = 1
	b = "a"
	if Equal(a, b) {
		t.Errorf("unexpected result: %v", a)
	}

	a = "a"
	b = "a"
	if !Equal(a, b) {
		t.Errorf("unexpected result: %v", a)
	}
}

func TestValueEqualOnSimpleArray(t *testing.T) {
	aa := []int{1, 3, 5, 7, 9}
	ab := []int{1, 3, 5, 7, 9}

	if !Equal(aa, ab) {
		t.Errorf("unexpected result: %v != %v", aa, ab)
	}

	ac := []int{1, 3, 5, 7, 9, 11}
	if Equal(aa, ac) {
		t.Errorf("unexpected result: %v == %v", aa, ac)
	}

	ad := []int32{1, 3, 5, 7, 9}
	if Equal(aa, ad) {
		t.Errorf("unexpected result: %v == %v", aa, ad)
	}

	ae := []int{1, 3, 6, 7, 9}
	if Equal(aa, ae) {
		t.Errorf("unexpected result: %v == %v", aa, ae)
	}
}

func TestValueEqualOnStruct(t *testing.T) {
	type test struct {
		Name string
		Age  int
	}

	a := test{"a", 1}
	b := test{"b", 2}
	c := test{"a", 1}

	if Equal(a, b) {
		t.Errorf("unexpected result: %+v == %+v", a, b)
	}

	if !Equal(a, c) {
		t.Errorf("unexpected result: %+v != %+v", a, c)
	}
}

func TestValueEqualOnStructWithPrivateField(t *testing.T) {
	type test struct {
		Name string
		age  int
	}

	a := test{"a", 1}
	b := test{"b", 2}
	c := test{"a", 1}

	if Equal(a, b) {
		t.Errorf("unexpected result: %+v == %+v", a, b)
	}

	if !Equal(a, c) {
		t.Errorf("unexpected result: %+v != %+v", a, c)
	}
}

package meta

import (
	"testing"

	"strings"
)

func TestCompareNil(t *testing.T) {
	type testType struct {
		Name string
	}

	var a *testType = nil
	b := &testType{"aaa"}
	c := testType{"aaa"}

	if InstanceEqual(a, b) {
		t.Errorf("unexpected result: nil <-> pointer")
	}

	if !InstanceEqual(b, c) {
		t.Errorf("unexpected result: pointer <-> instance")
	}

	if !InstanceEqual(a, nil) {
		t.Errorf("unexpected result: nil pointer <-> nil literal")
	}
}

func TestInstanceEqual(t *testing.T) {
	data := 3
	ptr := &data

	if !InstanceEqual(data, ptr) {
		t.Errorf("data does not equal")
	}
}

func TestInstanceEqualOnStruct(t *testing.T) {
	type testType struct {
		S string
		I int
	}

	data := testType{
		S: "the quick brown fox jumps over the lazy dog",
		I: 42,
	}

	ptr := &testType{
		S: "the quick brown fox jumps over the lazy dog",
		I: 42,
	}

	if !InstanceEqual(data, ptr) {
		t.Errorf("data does not equal")
	}
}

func TestArrayEqualInfoNonArray(t *testing.T) {
	m := []int{1, 2, 3, 4, 5, 6}
	n := 7

	{
		eq, err := ArrayEqualInfo(n, m)
		if eq {
			t.Errorf("unexpected result %v", eq)
		}

		if !strings.Contains(err.Error(), "a is not array") {
			t.Errorf("unexpected error message: %s", err)
		}
	}

	{
		eq, err := ArrayEqualInfo(m, n)
		if eq {
			t.Errorf("unexpected result %v", eq)
		}

		if !strings.Contains(err.Error(), "b is not array") {
			t.Errorf("unexpected error message: %s", err)
		}
	}
}

func TestArrayEqualInfoDiffLength(t *testing.T) {
	m := []int{1, 2, 3, 4, 5, 6}
	n := []int{1, 2, 3, 4}

	eq, err := ArrayEqualInfo(n, m)
	if eq {
		t.Errorf("unexpected result %v", eq)
	}

	if !strings.Contains(err.Error(), "Len") {
		t.Errorf("unexpected error message: %s", err)
	}
}

func TestArrayEqualInfoNotEqual(t *testing.T) {
	m := []int{1, 2, 3, 4, 5, 6}
	n := []int{1, 2, 3, 6, 6, 7}

	eq, err := ArrayEqualInfo(m, n)
	if eq {
		t.Errorf("unexpected result %v", eq)
	}

	if err.Error() != "a[3] != b[3]: 4 (int) <=> 6 (int)" {
		t.Errorf("unexpected error message: %s", err)
	}
}

func TestArrayEqualInfo(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 2, 3, 4, 5}
	c := []*int{intPtr(1), intPtr(2), intPtr(3), intPtr(4), intPtr(5)}

	{
		eq, err := ArrayEqualInfo(a, b)
		if !eq {
			t.Errorf("unexpected result %v", eq)
		}

		if err != nil {
			t.Errorf("unexpected error message: %s", err)
		}
	}

	{
		eq, err := ArrayEqualInfo(a, c)
		if eq {
			t.Errorf("unexpected result %v", eq)
		}

		if !strings.Contains(err.Error(), "a[0] != b[0]") {
			t.Errorf("unexpected error message: %s", err)
		}
	}

	{
		eq, err := ArrayInstanceEqualInfo(a, b)
		if !eq {
			t.Errorf("unexpected result %v", eq)
		}

		if err != nil {
			t.Errorf("unexpected error message: %s", err)
		}
	}

	{
		eq, err := ArrayInstanceEqualInfo(a, c)
		if !eq {
			t.Errorf("unexpected result %v", eq)
		}

		if err != nil {
			t.Errorf("unexpected error message: %s", err)
		}
	}
}

func TestArrayEqual(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 2, 3, 4, 5}
	c := []*int{intPtr(1), intPtr(2), intPtr(3), intPtr(4), intPtr(5)}

	if !ArrayEqual(a, b) {
		t.Errorf("unexpected result")
	}

	if ArrayEqual(a, c) {
		t.Errorf("unexpected result")
	}

	if !ArrayInstanceEqual(a, b) {
		t.Errorf("unexpected result")
	}

	if !ArrayInstanceEqual(a, c) {
		t.Errorf("unexpected result")
	}
}

func TestArrayItemEqualInfoNonArray(t *testing.T) {
	m := []int{1, 2, 3, 4, 5, 6}
	n := 7

	{
		eq, err := ArrayItemEqualInfo(n, m)
		if eq {
			t.Errorf("unexpected result %v", eq)
		}

		if !strings.Contains(err.Error(), "a is not array") {
			t.Errorf("unexpected error message: %s", err)
		}
	}

	{
		eq, err := ArrayItemEqualInfo(m, n)
		if eq {
			t.Errorf("unexpected result %v", eq)
		}

		if !strings.Contains(err.Error(), "b is not array") {
			t.Errorf("unexpected error message: %s", err)
		}
	}
}

func TestArrayItemEqualInfoDiffLength(t *testing.T) {
	m := []int{1, 2, 3, 4, 5, 6}
	n := []int{1, 2, 3, 4}

	eq, err := ArrayItemEqualInfo(n, m)
	if eq {
		t.Errorf("unexpected result %v", eq)
	}

	if !strings.Contains(err.Error(), "Len") {
		t.Errorf("unexpected error message: %s", err)
	}
}

func TestArrayItemEqualInfoNotContain(t *testing.T) {
	a := []int{1, 2, 3, 5, 7}
	b := []int{1, 2, 3, 5, 8}

	{
		eq, err := ArrayItemEqualInfo(a, b)
		if eq {
			t.Errorf("unexpected result %v", eq)
		}

		if err.Error() != "a[4] is not in b" {
			t.Errorf("unexpected error message: %s", err)
		}
	}

	{
		eq, err := ArrayItemEqualInfo(a, b)
		if eq {
			t.Errorf("unexpected result %v", eq)
		}

		if err.Error() != "a[4] is not in b" {
			t.Errorf("unexpected error message: %s", err)
		}
	}
}

func TestArrayItemEqual(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{5, 4, 3, 2, 1}
	c := []*int{intPtr(2), intPtr(1), intPtr(5), intPtr(3), intPtr(4)}

	if !ArrayItemEqual(a, b) {
		t.Errorf("unexpected result")
	}

	if ArrayItemEqual(a, c) {
		t.Errorf("unexpected result")
	}

	if !ArrayItemInstanceEqual(a, b) {
		t.Errorf("unexpected result")
	}

	if !ArrayItemInstanceEqual(a, c) {
		t.Errorf("unexpected result")
	}
}

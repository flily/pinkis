package meta

import (
	"reflect"
	"testing"
)

type testWandType struct {
	Length float32
	Core   string
	Wood   string
}

type testWizardType struct {
	Name  string
	Born  int
	Blood string
	Wand  *testWandType
}

func TestGetField(t *testing.T) {
	hermione := testWizardType{
		Name:  "Hermione Granger",
		Born:  1979,
		Blood: "muggle-born",
		Wand: &testWandType{
			Length: 10.75,
			Core:   "dragon heartstring",
			Wood:   "vine",
		},
	}

	{
		data, err := GetField(hermione, "born")
		if data != nil {
			t.Errorf("unexpected data: %v <=> %v", data, nil)
		}

		if err.Error() != "no field 'born'" {
			t.Errorf("unexpected error: %s", err)
		}
	}

	{
		exp := 1979
		data, err := GetField(hermione, "Born")
		if data != exp {
			t.Errorf("unexpected data: %v <=> %v", data, exp)
		}

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}

	ptr := &hermione
	{
		data, err := GetField(ptr, "born")
		if data != nil {
			t.Errorf("unexpected data: %v <=> %v", data, nil)
		}

		if err.Error() != "no field 'born'" {
			t.Errorf("unexpected error: %s", err)
		}
	}

	{
		exp := 1979
		data, err := GetField(ptr, "Born")
		if data != exp {
			t.Errorf("unexpected data: %v <=> %v", data, exp)
		}

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}
}

func TestSetField(t *testing.T) {
	hermione := testWizardType{
		Name:  "Hermione Granger",
		Born:  1979,
		Blood: "muggle-born",
		Wand: &testWandType{
			Length: 10.75,
			Core:   "dragon heartstring",
			Wood:   "vine",
		},
	}

	{
		data, err := SetField(hermione, "blood", "half-blood")
		if data != nil {
			t.Errorf("unexpected data: %v <=> %v", data, nil)
		}

		if err.Error() != "no field 'blood'" {
			t.Errorf("unexpected error: %s", err)
		}
	}

	{
		// Can not update field of a copyed data
		data, err := SetField(hermione, "Blood", "pure-blood")
		if data != nil {
			t.Errorf("unexpected data: %v <=> %v", data, nil)
		}

		if err.Error() != "field 'Blood' can not be set" {
			t.Errorf("unexpected error: %s", err)
		}
	}

	ptr := &hermione
	{
		data, err := SetField(ptr, "blood", "half-blood")
		if data != nil {
			t.Errorf("unexpected data: %v <=> %v", data, nil)
		}

		if err.Error() != "no field 'blood'" {
			t.Errorf("unexpected error: %s", err)
		}
	}

	{
		data, err := SetField(ptr, "Blood", "pure-blood")
		if data != "pure-blood" {
			t.Errorf("unexpected data: %v <=> %v", data, nil)
		}

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}
}

func TestSetFieldWithConvert(t *testing.T) {
	hermione := &testWizardType{
		Name:  "Hermione Granger",
		Born:  1979,
		Blood: "muggle-born",
		Wand: &testWandType{
			Length: 10.75,
			Core:   "dragon heartstring",
			Wood:   "vine",
		},
	}

	{
		data, err := SetField(hermione, "Born", int64(1999))
		if data != 1999 {
			t.Errorf("unexpected data: %v <=> %v", data, nil)
		}

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}

	{
		data, err := SetField(hermione, "Born", "pure-blood")
		if data != nil {
			t.Errorf("unexpected data: %v <=> %v", data, nil)
		}

		if err.Error() != "field 'Born' requires type int, but string" {
			t.Errorf("unexpected error: %s", err)
		}
	}

	expected := &testWizardType{
		Name:  "Hermione Granger",
		Born:  1999,
		Blood: "muggle-born",
		Wand: &testWandType{
			Length: 10.75,
			Core:   "dragon heartstring",
			Wood:   "vine",
		},
	}

	if hermione == expected {
		t.Errorf("pointer must not be equal")
	}

	if !InstanceEqual(hermione, expected) {
		t.Errorf("unexpected data: %v <=> %v", hermione, expected)
	}

	if !reflect.DeepEqual(hermione, expected) {
		t.Errorf("unexpected data: %v <=> %v", hermione, expected)
	}
}

package meta

import (
	"reflect"
	"testing"
)

func TestIsPointer(t *testing.T) {
	v := 42
	if IsPointer(v) {
		t.Error("IsPointer(v) should be false")
	}

	if !IsPointer(&v) {
		t.Error("IsPointer(&v) should be true")
	}
}

func TestIsPointerOfNil(t *testing.T) {
	var v *int
	if !IsPointer(v) {
		t.Error("IsPointer(v) should be true")
	}

	if IsPointer(nil) {
		t.Error("IsPointer(nil) should be false")
	}
}

func TestIsNil(t *testing.T) {
	var i *int
	if u, n := IsNil(i); u || !n {
		t.Errorf("IsNil(i) should be (false, true), got (%v, %v)", u, n)
	}

	if u, n := IsNil(nil); !u || n {
		t.Errorf("IsNil(nil) should be (true, false), got (%v, %v)", u, n)
	}
}

type ttIsStruct1 struct {
}

func (ttIsStruct1) One()    {}
func (ttIsStruct1) Two()    {}
func (*ttIsStruct1) Three() {}
func (*ttIsStruct1) Four()  {}

type ttIsStruct1Interface1 interface {
	One()
	Two()
}

type ttIsStruct1Interface2 interface {
	Three()
	Four()
}

func TestStructBasis(t *testing.T) {
	data := ttIsStruct1{}

	{
		var i ttIsStruct1Interface1 = data
		v := reflect.ValueOf(i)

		if v.Kind() != reflect.Struct {
			t.Errorf("%s is not struct: %s", v, v.Kind())
		}
	}

	{
		var i ttIsStruct1Interface2 = &data

		v := reflect.ValueOf(i)
		if v.Kind() != reflect.Ptr {
			t.Errorf("%s is not ptr: %s", v, v.Kind())
		}
	}

	{
		var i ttIsStruct1Interface2 = &data

		v := reflect.ValueOf(&i).Type().Elem()
		if v.Kind() != reflect.Interface {
			t.Errorf("%s is not interface: %s", v, v.Kind())
		}
	}

	{
		var i ttIsStruct1Interface2 = &data
		x := &i

		v := reflect.ValueOf(*x)
		if v.Kind() != reflect.Ptr {
			t.Errorf("%s is not ptr: %s", v, v.Kind())
		}
	}
}

func TestIsStruct(t *testing.T) {
	{
		v := ttIsStruct1{}
		if !IsStruct(v) {
			t.Error("IsStruct(v) should be true")
		}

		if !IsStruct(&v) {
			t.Error("IsStruct(&v) should be true")
		}
	}

	{
		var v ttIsStruct1Interface1 = ttIsStruct1{}
		if !IsStruct(v) {
			t.Error("IsStruct(v) should be true")
		}

		if !IsStruct(&v) {
			t.Error("IsStruct(&v) should be true")
		}
	}
}

func TestIsExportedName(t *testing.T) {
	caseList := []struct {
		Name       string
		IsExported bool
	}{
		{"", false},
		{"a", false},
		{"A", true},
		{"_x9", false},
		{"ThisVariableIsExported", true},
		{"english", false},
		{"English", true},
		{"ελληνικά", false},
		{"Ελληνικά", true},
		{"русские", false},
		{"Русские", true},
		{"中文", false},
		{"日本語", false},
		{"にほんご", false},
		{"ニホンゴ", false},
	}

	for _, kase := range caseList {
		got := IsExportedName(kase.Name)
		if got != kase.IsExported {
			t.Errorf("IsExportedName(%q) = %v, expect %v", kase.Name, got, kase.IsExported)
		}
	}
}

func intPtr(n int) *int {
	return &n
}

func TestInterfaceEqual(t *testing.T) {
	data1 := 3
	value1 := reflect.ValueOf(data1)

	data2 := int64(3)
	value2 := reflect.ValueOf(data2)

	if value1.Interface() == value2.Interface() {
		t.Errorf("unexpected equal")
	}

	data3 := 1 + 2
	value3 := reflect.ValueOf(data3)
	if value1.Interface() != value3.Interface() {
		t.Errorf("data does not equal")
	}
}

func TestInstanceOf(t *testing.T) {
	data := 3

	dataPtr := &data
	value := reflect.ValueOf(data)
	{
		got := InstanceOf(dataPtr)
		if got.Interface() != value.Interface() {
			t.Errorf("got error value: %v <=> %v", got, value)
		}
	}

	{
		got := reflect.Indirect(reflect.ValueOf(dataPtr))
		if got.Interface() != value.Interface() {
			t.Errorf("got error value: %v <=> %v", got, value)
		}
	}

	dataPtrPtr := &dataPtr
	{
		got := InstanceOf(dataPtrPtr)
		if got.Interface() != value.Interface() {
			t.Errorf("got error value: %v <=> %v", got, value)
		}
	}

	{
		// NOTES: reflect.Indirect returns the first data that points to.
		got := reflect.Indirect(reflect.ValueOf(dataPtrPtr))
		if got.Interface() == value.Interface() {
			t.Errorf("got error value: %v <=> %v", got, value)
		}
	}

	{
		var nilPtr *int = nil
		got := InstanceOf(nilPtr)
		if !got.IsValid() {
			t.Errorf("got error value: %v, is valid %v", got, got.IsValid())
		}

		if !got.IsNil() {
			t.Errorf("got error value: %v <=> %v", got, value)
		}
	}

	{
		// Untyped nil
		got := InstanceOf(nil)
		if got.IsValid() {
			t.Errorf("got error value: %v, is valid %v", got, got.IsValid())
		}
	}
}

type ttTypeCastData struct {
	Name string
	Age  int
}

func (t *ttTypeCastData) Say(word string) string {
	return t.Name + " says " + word
}

func (t ttTypeCastData) Shout(word string) string {
	return t.Name + " shouts " + word + "!"
}

type ttTypeSayingInterface interface {
	Say(string) string
}

type ttTypeShoutingInterface interface {
	Shout(string) string
}

func TestMetaTypeCasting(t *testing.T) {
	data := ttTypeCastData{
		Name: "Lily",
		Age:  33,
	}

	var infSay ttTypeSayingInterface = &data
	{
		tt := reflect.ValueOf(infSay)
		if !tt.CanConvert(reflect.TypeOf(&data)) {
			t.Errorf("%s can not convert to %s", tt.Type(), reflect.TypeOf(&data))
		}
	}

	var infShout ttTypeShoutingInterface = data
	{
		tt := reflect.ValueOf(infShout)
		if !tt.CanConvert(reflect.TypeOf(data)) {
			t.Errorf("%s can not convert to %s", tt.Type(), reflect.TypeOf(&data))
		}
	}
}

func TestOrigineOf(t *testing.T) {
	data := ttTypeCastData{
		Name: "Lily",
		Age:  33,
	}

	ptr := &data

	{
		var inf ttTypeSayingInterface = ptr

		v := reflect.ValueOf(inf)
		if v.Kind() != reflect.Ptr {
			t.Errorf("%s is not pointer", v.Kind())
		}

		ins, chain := ValueInstanceChainOf(inf)
		ori := OriginOf(ins, chain)

		if *(ori.Interface().(*ttTypeCastData)) != *ptr {
			t.Errorf("origin of %v is %v, expect %v", ins, ori, ptr)
		}
	}

	{
		var inf ttTypeShoutingInterface = data

		v := reflect.ValueOf(inf)
		if v.Kind() != reflect.Struct {
			t.Errorf("%s is not Struct", v.Kind())
		}

		ins, chain := ValueInstanceChainOf(inf)
		ori := OriginOf(ins, chain)

		if ori.Interface() != data {
			t.Errorf("origin of %v is %v, expect %v", ins, ori, data)
		}
	}

	{
		var s ttTypeSayingInterface = ptr
		inf := &s

		v := reflect.ValueOf(inf)
		if v.Kind() != reflect.Ptr {
			t.Errorf("%s is not pointer", v.Kind())
		}

		ins, chain := ValueInstanceChainOf(inf)
		ori := OriginOf(ins, chain)

		p, ok := ori.Interface().(*ttTypeSayingInterface)
		if !ok {
			t.Errorf("origin of %v is %v, expect %v", ins, ori, ptr)
		}

		if *(*p).(*ttTypeCastData) != *ptr {
			t.Errorf("origin of %v is %v, expect %v", ins, ori, ptr)
		}
	}
}

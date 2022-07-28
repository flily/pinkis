package meta

import (
	"reflect"
)

func InstanceEqual(a interface{}, b interface{}) bool {
	ia := InstanceOf(a)
	ib := InstanceOf(b)

	return ValueEqual(ia, ib)
}

func iArrayEqual(a interface{}, b interface{}, compareInstance bool) (bool, error) {
	valueA := reflect.ValueOf(a)
	if valueA.Kind() != reflect.Array && valueA.Kind() != reflect.Slice {
		e := NewMetaError("a is not array or slice, but %s", valueA.Kind())
		return false, e
	}

	valueB := reflect.ValueOf(b)
	if valueB.Kind() != reflect.Array && valueB.Kind() != reflect.Slice {
		e := NewMetaError("b is not array or slice, but %s", valueA.Kind())
		return false, e
	}

	lengthA, lengthB := valueA.Len(), valueB.Len()
	if lengthA != lengthB {
		e := NewMetaError("a.Len()=%d <=> b.Len()=%d", lengthA, lengthB)
		return false, e
	}

	for i := 0; i < lengthA; i++ {
		itemA := valueA.Index(i)
		itemB := valueB.Index(i)

		if compareInstance {
			itemA = ValueInstanceOf(itemA)
			itemB = ValueInstanceOf(itemB)
		}

		if !ValueEqual(itemA, itemB) {
			e := NewMetaError("a[%d] != b[%d]: %v (%s) <=> %v (%s)",
				i, i, itemA, itemA.Type(), itemB, itemB.Type())
			return false, e
		}
	}

	return true, nil
}

func ArrayEqualInfo(a interface{}, b interface{}) (bool, error) {
	return iArrayEqual(a, b, false)
}

func ArrayEqual(a interface{}, b interface{}) bool {
	eq, _ := ArrayEqualInfo(a, b)
	return eq
}

func ArrayInstanceEqualInfo(a interface{}, b interface{}) (bool, error) {
	return iArrayEqual(a, b, true)
}

func ArrayInstanceEqual(a interface{}, b interface{}) bool {
	eq, _ := ArrayInstanceEqualInfo(a, b)
	return eq
}

func iArrayItemEqual(a interface{}, b interface{}, compareInstance bool) (bool, error) {
	valueA := reflect.ValueOf(a)
	if valueA.Kind() != reflect.Array && valueA.Kind() != reflect.Slice {
		e := NewMetaError("a is not array or slice, but %s", valueA.Kind())
		return false, e
	}

	valueB := reflect.ValueOf(b)
	if valueB.Kind() != reflect.Array && valueB.Kind() != reflect.Slice {
		e := NewMetaError("b is not array or slice, but %s", valueA.Kind())
		return false, e
	}

	lengthA, lengthB := valueA.Len(), valueB.Len()
	if lengthA != lengthB {
		e := NewMetaError("a.Len()=%d <=> b.Len()=%d", lengthA, lengthB)
		return false, e
	}

	itemMap := map[int]int{}
	for i := 0; i < lengthA; i++ {
		itemMap[i] = -1
	}

	for i := 0; i < lengthA; i++ {
		itemA := valueA.Index(i)
		if compareInstance {
			itemA = ValueInstanceOf(itemA)
		}

		matched := false
		for j := 0; j < lengthB; j++ {
			if itemMap[j] != -1 {
				// compared
				continue
			}

			itemB := valueB.Index(j)
			if compareInstance {
				itemB = ValueInstanceOf(itemB)
			}

			if ValueEqual(itemA, itemB) {
				itemMap[j] = i
				matched = true
				break
			}
		}

		if !matched {
			e := NewMetaError("a[%d] is not in b", i)
			return false, e
		}
	}

	// may be imposible to find b[j] not in a
	// for j, i := range itemMap {
	// 	if i == -1 {
	// 		e := NewMetaError("b[%d] is not in a", j)
	// 		return false, e
	// 	}
	// }

	return true, nil
}

func ArrayItemEqualInfo(a interface{}, b interface{}) (bool, error) {
	return iArrayItemEqual(a, b, false)
}

func ArrayItemEqual(a interface{}, b interface{}) bool {
	eq, _ := ArrayItemEqualInfo(a, b)
	return eq
}

func ArrayItemInstanceEqualInfo(a interface{}, b interface{}) (bool, error) {
	return iArrayItemEqual(a, b, true)
}

func ArrayItemInstanceEqual(a interface{}, b interface{}) bool {
	eq, _ := ArrayItemInstanceEqualInfo(a, b)
	return eq
}

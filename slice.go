package slicehelper

import (
	"reflect"
)

// ReplaceNilWithEmptySlice recursively scans all field of the given structure, replacing all
// nil slice with empty slice, and returns a POINTER to the given structure.
func ReplaceNilWithEmptySlice(input interface{}) interface{} {
	val := reflect.ValueOf(input)
	switch val.Kind() {
	case reflect.Ptr:
		// if, in some case, a pointer is passed in, We dereference it and do the normal stuff
		if val.IsNil() {
			if reflect.TypeOf(input).Elem().Kind() == reflect.Slice {
				res := reflect.New(reflect.TypeOf(input).Elem())
				resp := reflect.MakeSlice(reflect.TypeOf(input).Elem(), 0,0)
				res.Elem().Set(resp)
				return res.Interface()
			}
			return input
		}
		return ReplaceNilWithEmptySlice(val.Elem().Interface())
	case reflect.Slice:
		res := reflect.New(val.Type())
		resp := reflect.MakeSlice(val.Type(), 0, val.Len())
		// if this is not empty, copy it
		if !val.IsZero() {
			// iterate over each element in slice
			for i := 0; i < val.Len(); i++ {
				item := val.Index(i)
				var newItem reflect.Value
				switch item.Kind() {
				case reflect.Struct, reflect.Slice, reflect.Map:
					// recursively handle nested struct
					newItem = reflect.Indirect(reflect.ValueOf(ReplaceNilWithEmptySlice(item.Interface())))
				case reflect.Ptr:
					if item.IsNil() {
						if item.Type().Elem().Kind() == reflect.Slice {
							res := reflect.New(item.Type().Elem())
							resp := reflect.MakeSlice(item.Type().Elem(), 0, 0)
							res.Elem().Set(resp)
							newItem = res
							break
						}
						newItem = item
						break
					}
					if item.Elem().Kind() == reflect.Struct || item.Elem().Kind() == reflect.Slice {
						newItem = reflect.ValueOf(ReplaceNilWithEmptySlice(item.Elem().Interface()))
						break
					}
					fallthrough
				default:
					newItem = item
				}
				resp = reflect.Append(resp, newItem)
			}
		}
		res.Elem().Set(resp)
		return res.Interface()
	case reflect.Struct:
		resp := reflect.New(reflect.TypeOf(input))
		newVal := resp.Elem()
		// iterate over input's fields
		for i := 0; i < val.NumField(); i++ {
			newValField := newVal.Field(i)
			if !newValField.CanSet() {
				continue
			}
			valField := val.Field(i)
			updates := reflect.ValueOf(ReplaceNilWithEmptySlice(valField.Interface()))
			if updates.IsValid() {
				if valField.Kind() == reflect.Ptr {
					newValField.Set(updates)
				} else {
					newValField.Set(reflect.Indirect(updates))
				}
			}
		}
		return resp.Interface()
	case reflect.Map:
		res := reflect.New(reflect.TypeOf(input))
		resp := reflect.MakeMap(reflect.TypeOf(input))
		// iterate over every key value pair in input map
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			// recursively handle nested value
			newV := reflect.Indirect(reflect.ValueOf(ReplaceNilWithEmptySlice(v.Interface())))
			resp.SetMapIndex(k, newV)
		}
		res.Elem().Set(resp)
		return res.Interface()
	default:
		return input
	}
}

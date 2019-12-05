package main

import (
	"errors"
	"fmt"
	"reflect"
)

func i2s(in interface{}, out interface{}) error {
	inReflect := reflect.ValueOf(in)
	outReflect := reflect.ValueOf(out)
	if outReflect.Kind() != reflect.Ptr {
		return errors.New("out param isn`t a pointer")
	}
	outReflect = outReflect.Elem()
	switch inReflect.Kind() {
	case reflect.Slice:
		if outReflect.Kind() != reflect.Slice {
			return errors.New("arguments mismatch")
		}
		inSlice := reflect.ValueOf(inReflect.Interface())
		outSlice := reflect.MakeSlice(outReflect.Type(), inSlice.Len(), inSlice.Len())
		for i := 0; i < inSlice.Len(); i++ {
			err := i2s(inSlice.Index(i).Interface(), outSlice.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}
		outReflect.Set(outSlice)
		break
	case reflect.Map:
		if outReflect.Kind() != reflect.Struct {
			return errors.New(fmt.Sprintf("arguments mismatch"))
		}
		mapKeys := inReflect.MapKeys()
		if len(mapKeys) == 0 {
			errors.New("empty inMap")
		}
		for _, fieldName := range mapKeys {
			inVal := inReflect.MapIndex(fieldName).Elem()
			outVal := outReflect.FieldByName(fieldName.String())

			//fmt.Printf("\nkey: %s", fieldName)
			//fmt.Printf("\nvalue: %#v", inVal.Interface())
			//fmt.Printf("\ntype: %#v\n", inVal.Type().String())
			//fmt.Printf("out type: %#v\n", outVal.Type().String())

			switch reflect.TypeOf(inVal.Interface()).Kind() {
			case reflect.Map:
				fieldVal := reflect.New(outVal.Type()).Elem()
				err := i2s(inVal.Interface(), fieldVal.Addr().Interface())
				if err != nil {
					return err
				}
				outVal.Set(fieldVal)
				break
			case reflect.Slice:
				inSlice := reflect.ValueOf(inVal.Interface())
				outSlice := reflect.MakeSlice(outVal.Type(), inSlice.Len(), inSlice.Len())
				for i := 0; i < inSlice.Len(); i++ {
					err := i2s(inSlice.Index(i).Interface(), outSlice.Index(i).Addr().Interface())
					if err != nil {
						return err
					}
				}
				outVal.Set(outSlice)
				break
			case reflect.Float64:
				castVal, ok := inVal.Interface().(float64)
				if !ok {
					return errors.New("can`t cast to float64")
				}
				if outVal.Type().String() != "int" {
					return errors.New("argument type mismatch")
				}
				outVal.SetInt(int64(castVal))
				break
			case reflect.String:
				val, ok := inVal.Interface().(string)
				if !ok {
					return errors.New("can`t cast to string")
				}
				if outVal.Type().String() != "string" {
					return errors.New("argument type mismatch")
				}
				outVal.SetString(val)
				break
			case reflect.Bool:
				val, ok := inVal.Interface().(bool)
				if !ok {
					return errors.New("can`t cast to bool")
				}
				if outVal.Type().String() != "bool" {
					return errors.New("argument type mismatch")
				}
				outVal.SetBool(val)
				break
			}
		}
		break
	}

	return nil
}

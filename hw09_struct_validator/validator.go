package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	result := ""
	for _, err := range v {
		if _, ok := err.Err.(ValidationErrors); ok {
			result += err.Field + ":{\n" + fmt.Sprint(err.Err) + "}"
		} else {
			result = result + err.Field + ":" + fmt.Sprint(err.Err) + "!\n"
		}
	}
	return result
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	validationErrors := make(ValidationErrors, 0)
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		_, ok := rt.Field(i).Tag.Lookup("validate")
		if !ok {
			continue
		}
		lenString := reflect.StructTag(rt.Field(i).Tag.Get("validate"))
		switch rt.Field(i).Type.Kind() {
		case reflect.Struct:
			tag := reflect.StructTag(rt.Field(i).Tag.Get("validate"))
			if tag != "nested" {
				continue
			}
			validationErrors = append(validationErrors, ValidationError{
				Field: rt.Field(i).Name,
				Err:   Validate(rv.Field(i).Interface()),
			})
		case reflect.String:
			validationErrs := validateField(rv.Field(i).String(), string(lenString))
			for _, err := range validationErrs {
				validationErrors = append(validationErrors, ValidationError{
					Field: rt.Field(i).Name,
					Err:   err,
				})
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			validationErrs := validateField(rv.Field(i).Int(), string(lenString))
			for _, err := range validationErrs {
				validationErrors = append(validationErrors, ValidationError{
					Field: rt.Field(i).Name,
					Err:   err,
				})
			}
		case reflect.Slice:
			for l := 0; l < rv.Field(i).Len(); l++ {
				validationErrs := validateField(rv.Field(i).Index(l), string(lenString))
				for _, err := range validationErrs {
					validationErrors = append(validationErrors, ValidationError{
						Field: rt.Field(i).Name,
						Err:   err,
					})
				}
			}
		default:

		}

	}
	return validationErrors
}

func validateField(field interface{}, tag string) []error {
	validationTag := make([]error, 0)
	for _, prop := range strings.Split(tag, "|") {
		options := strings.Split(prop, ":")
		switch options[0] {
		case "len":
			temp, _ := strconv.Atoi(options[1])
			if temp != len(fmt.Sprint(field)) {
				validationTag = append(validationTag, fmt.Errorf("%v doesnt satisfy condition: %s", field, prop))
			}
		case "min":
			temp, _ := strconv.ParseInt(options[1], 10, 64)
			if field.(int64) < temp {
				validationTag = append(validationTag, fmt.Errorf("%v doesnt satisfy condition: %s", field, prop))
			}
		case "max":
			temp, _ := strconv.ParseInt(options[1], 10, 64)
			if field.(int64) > temp {
				validationTag = append(validationTag, fmt.Errorf("%v doesnt satisfy condition: %s", field, prop))
			}
		case "regexp":
			temp, _ := regexp.MatchString(options[1], fmt.Sprint(field))
			if !temp {
				validationTag = append(validationTag, fmt.Errorf("%v doesnt satisfy condition: %s", field, prop))
			}
		case "in":
			temp := strings.Split(options[1], ",")
			if !slices.Contains(temp, fmt.Sprint(field)) {
				validationTag = append(validationTag, fmt.Errorf("%v doesnt satisfy condition: %s", field, prop))
			}
		}
	}
	return validationTag
}

//
//func isZero(v reflect.Value) bool {
//	switch v.Kind() {
//	case reflect.String:
//		return v.String() == ""
//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//		return v.Int() == 0
//	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
//		return v.Uint() == 0
//	case reflect.Float32, reflect.Float64:
//		return v.Float() == 0
//	case reflect.Bool:
//		return !v.Bool()
//	case reflect.Interface, reflect.Ptr:
//		return v.IsNil()
//	}
//	return false
//}

package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type User struct {
	Name string `validate:"min=2,max=32"`
	Age  int    `validate:"min=0,max=130"`
}

func Validate(v interface{}) []error {
	var errors []error
	val := reflect.ValueOf(v)
	typ := val.Type()

	// Iterate over struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typ.Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}

		// Parse validation rules
		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			parts := strings.Split(rule, "=")
			if len(parts) != 2 {
				continue
			}

			name := parts[0]
			param, _ := strconv.Atoi(parts[1])

			// Apply validation rules
			switch name {
			case "min":
				switch field.Kind() {
				case reflect.String:
					if len(field.String()) < param {
						errors = append(errors, fmt.Errorf("field %s is too short", typ.Field(i).Name))
					}
				case reflect.Int:
					if field.Int() < int64(param) {
						errors = append(errors, fmt.Errorf("field %s is too small", typ.Field(i).Name))
					}
				}
			case "max":
				switch field.Kind() {
				case reflect.String:
					if len(field.String()) > param {
						errors = append(errors, fmt.Errorf("field %s is too long", typ.Field(i).Name))
					}
				case reflect.Int:
					if field.Int() > int64(param) {
						errors = append(errors, fmt.Errorf("field %s is too large", typ.Field(i).Name))
					}
				}
			}
		}
	}

	return errors
}

func main() {
	user := User{
		Name: "A",  // Too short
		Age:  150,  // Too large
	}

	if errs := Validate(user); len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	}
}

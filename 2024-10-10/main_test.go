package main

import "testing"

func TestSaySomething(t *testing.T) {
	subtests := []struct{ 
		message string
		result string
	} {
		{
			message: "Hello",
			result: "Here it is Hello",
		},
		{
			message: "Hey there",
			result: "Here it is Hey there",
		},
	}

	for _, test := range subtests {
		if result := SaySomething(test.message); test.result != result {
			t.Errorf("wanted '%s', but got '%s'", test.result, result)
		}
	}
}
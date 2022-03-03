package lgwt_reflection

import (
	"reflect"
	"testing"
)

func TestWalk(t *testing.T) {
	expected := "Chris"
	var got []string

	x := struct {
		Name string ``
	}{
		expected,
	}
	walk(x, func(input string) {
		got = append(got, input)
	})
	if len(got) != 1 {
		t.Errorf("wrong number of function calls, got %d want %d", len(got), 1)
	}
	if got[0] != expected {
		t.Errorf("got %q, want %q", got[0], expected)
	}
}

func TestWalkTable(t *testing.T) {
	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{"Struct with one string field",
			struct{ Name string }{"Chris"},
			[]string{"Chris"},
		},
		{
			"Struct with 2 string fields",
			struct {
				Name string
				City string
			}{
				"Chris",
				"London",
			},
			[]string{
				"Chris",
				"London",
			},
		},
		{
			"Struct with non string field",
			struct {
				Name string
				Age  int
			}{
				"Chris",
				33,
			},
			[]string{
				"Chris",
			},
		},
		{
			"Nested fields",
			struct {
				Name    string
				Profile struct {
					Age  int
					City string
				}
			}{
				"Chris",
				struct {
					Age  int
					City string
				}{
					33,
					"London",
				},
			},
			[]string{"Chris", "London"},
		},
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})
			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, expected %v", got, test.ExpectedCalls)
			}
		})
	}

}

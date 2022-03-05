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
        {
            "Pointers to things",
            &Person{
                "Chris",
                Profile{33, "London"},
            },
            []string{"Chris", "London"},
        },
        {
            "Slices",
            []Profile{
                {33, "London"},
                {34, "Reykjavik"},
            },
            []string{"London", "Reykjavik"},
        },
        {
            "Arrays",
            [2]Profile{
                {33, "London"},
                {34, "Reykjavik"},
            },
            []string{"London", "Reykjavik"},
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

    t.Run("with maps", func(t *testing.T) {

        aMap := map[string]string{
            "Foo": "Bar",
            "Baz": "Boz",
        }
        var got []string
        walk(aMap, func(input string) {
            got = append(got, input)
        })
        assertContains(t, got, "Bar")
        assertContains(t, got, "Boz")

    })
    t.Run("with channels", func(t *testing.T) {
        aChannel := make(chan Profile)
        go func() {
            aChannel <- Profile{33, "Berlin"}
            aChannel <- Profile{34, "Katowice"}
            close(aChannel)
        }()
        var got []string
        want := []string{"Berlin", "Katowice"}
        walk(aChannel, func(input string) {
            got = append(got, input)
        })
        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v, want %v", got, want)
        }
    })

    t.Run("with function", func(t *testing.T) {
        aFunction := func() (Profile, Profile) {
            return Profile{33, "Berlin"}, Profile{34, "Katowice"}
        }
        var got []string
        want := []string{"Berlin", "Katowice"}
        walk(aFunction, func(input string) {
            got = append(got, input)
        })
        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v, want %v", got, want)
        }
    })
}

func assertContains(t testing.TB, haystack []string, needle string) {
    t.Helper()
    for _, val := range haystack {
        if val == needle {
            return
        }
    }
    t.Errorf("expected %v to contain %q, but it didn't", haystack, needle)
}

type Person struct {
    Name    string
    Profile Profile
}

type Profile struct {
    Age  int
    City string
}

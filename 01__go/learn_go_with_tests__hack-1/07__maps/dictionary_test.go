package g3l_maps

import "testing"

func TestSearch(t *testing.T) {
    dict := Dict{"test": "this is just a test"}
    t.Run("known word", func(t *testing.T) {
        got, err := dict.Search("test")
        assertNoError(t, err)
        want := "this is just a test"
        assertStrings(t, got, want)
    })

    t.Run("unknown word", func(t *testing.T) {
        _, err := dict.Search("unknown")
        assertError(t, err, ErrNotFound)
    })
}

func TestAdd(t *testing.T) {
    t.Run("new word", func(t *testing.T) {
        dict := Dict{}
        err := dict.Add("test", "this is just a test")
        assertNoError(t, err)
        want := "this is just a test"
        got, err := dict.Search("test")
        assertNoError(t, err)
        assertStrings(t, got, want)
    })

    t.Run("existing word", func(t *testing.T) {
        key := "test"
        val := "this is just a test"
        dict := Dict{key: val}
        err := dict.Add(key, val)
        assertError(t, err, ErrWordExists)
    })
}

func TestUpdate(t *testing.T) {
    t.Run("existing word", func(t *testing.T) {
        key := "test"
        val := "this is just a test"
        dict := Dict{key: val}
        newVal := "this is a new val"
        err := dict.Update(key, newVal)
        assertNoError(t, err)
        got, err := dict.Search(key)
        assertNoError(t, err)
        assertStrings(t, got, newVal)
    })

    t.Run("new word", func(t *testing.T) {
        key := "test"
        val := "this is just a test"
        dict := Dict{}
        err := dict.Update(key, val)
        assertError(t, err, ErrWordDoesNotExist)
    })
}

func TestDelete(t *testing.T) {
    t.Run("existing key", func(t *testing.T) {
        key := "test"
        val := "this is just a test"
        dict := Dict{key: val}
        err := dict.Delete(key)
        assertNoError(t, err)
        _, err = dict.Search(key)
        assertError(t, err, ErrNotFound)
    })

    t.Run("unknow key", func(t *testing.T) {
        key := "test"
        dict := Dict{}
        err := dict.Delete(key)
        assertError(t, err, ErrWordDoesNotExist)
    })
}

func assertStrings(t *testing.T, got string, want string) {
    t.Helper()
    if got != want {
        t.Errorf("got %q want %q given %q", got, want, "test")
    }
}

func assertError(t testing.TB, got, want error) {
    t.Helper()
    if got != want {
        t.Errorf("got %v want %v", got, want)
    }
}

func assertNoError(t testing.TB, got error) {
    t.Helper()
    if got != nil {
        t.Errorf("expected no error but got %v", got)
    }
}

package _5__structs_methods_and_interfaces

import "testing"

//func TestPerimeter(t *testing.T) {
//    got := Perimeter(10.0, 10.0)
//    want := 40.0
//    if got != want {
//        t.Errorf("got %.2f want %.2f", got, want)
//    }
//}
//
//func TestArea(t *testing.T) {
//    got := Area(12.0, 16.0)
//    want := 72.0
//    if got != want {
//        t.Errorf("got %.2f want %.2f", got, want)
//    }
//}

func TestPerimeter(t *testing.T) {
    rectangle := Rectangle{
        Width:  10.0,
        Height: 10.0,
    }
    got := Perimeter(rectangle)
    want := 40.0
    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}

func TestArea(t *testing.T) {
    t.Run("rectangles", func(t *testing.T) {

        rectangle := Rectangle{
            Width:  10.0,
            Height: 10.0,
        }
        got := rectangle.Area()
        want := 100.0
        if got != want {
            t.Errorf("got %.2f want %.2f", got, want)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{
            radius: 10.0,
        }
        got := circle.Area()
        want := 314.1592653589793
        if got != want {
            t.Errorf("got %g want %g", got, want)
        }
    })
}

func TestAreaWithInterface(t *testing.T) {
    checkArea := func(t testing.TB, shape Shape, want float64) {
        t.Helper()
        got := shape.Area()
        if got != want {
            t.Errorf("got %g want %g", got, want)
        }
    }

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{Width: 12, Height: 6}
        checkArea(t, rectangle, 72.0)
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        checkArea(t, circle, 314.1592653589793)
    })
}

func TestAreaTable(t *testing.T) {
    areaTests := []struct {
        shape Shape
        want  float64
    }{
        {Rectangle{6, 12}, 72},
        {Circle{10}, 314.1592653589793},
        {Triangle{12, 6}, 36.0},
    }
    for _, test := range areaTests {
        got := test.shape.Area()
        if got != test.want {
            t.Errorf("%#v, got %g, wang %g", test.shape, got, test.want)
        }
    }
}

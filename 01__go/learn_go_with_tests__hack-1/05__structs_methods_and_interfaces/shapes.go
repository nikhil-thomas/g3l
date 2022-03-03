package _5__structs_methods_and_interfaces

import (
    "math"
)

type Shape interface {
    Area() float64
}

// Rectangle represents a rectangle
type Rectangle struct {
    Width  float64
    Height float64
}

// Circle represents a circle
type Circle struct {
    radius float64
}

type Triangle struct {
    Base   float64
    Height float64
}

// Perimeter returns perimeter of a rectangle
func Perimeter(r Rectangle) float64 {
    return 2 * (r.Width + r.Height)
}

//// Area returns area of a rectangle
//func Area(r Rectangle) float64 {
//    return r.Width * r.Height
//}

func (r Rectangle) Area() float64 {
    return r.Height * r.Width
}

func (c Circle) Area() float64 {
    return math.Pi * c.radius * c.radius
}

func (t Triangle) Area() float64 {
    return (t.Height * t.Base) / 2
}

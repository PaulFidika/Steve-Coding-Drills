package main

import (
	"fmt"
	"math"
)

type Point struct {
	X float64
	Y float64
}

func (p Point) Offset(x float64, y float64) Point {
	return Point{p.X + x, p.Y + y}
}

func (p *Point) Move(x float64, y float64) {
	p.X += x
	p.Y += y
}

type Line struct {
	start Point
	end Point
}

func (line Line) Length() float64 {
	return math.Hypot(line.start.X - line.end.X, line.start.Y - line.end.Y)
}

type Polygon []Line

func runCoordinates() Polygon {
	var line Line = Line{Point{0, 0}, Point{1, 1}}
	fmt.Println("Line: ", line, "Length: ", line.Length())

	// Unit box
	var unitBox = Polygon{
		Line{Point{0, 0}, Point{1, 0}},
		Line{Point{1, 0}, Point{1, 1}},
		Line{Point{1, 1}, Point{0, 1}},
		Line{Point{0, 1}, Point{0, 0}},
	}

	return unitBox
}
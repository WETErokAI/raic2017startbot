package codewars

import (
	"math"
)

type Unit struct {
	Id   int64
	X, Y float64
}

func (u *Unit) GetId() int64  { return u.Id }
func (u *Unit) GetX() float64 { return u.X }
func (u *Unit) GetY() float64 { return u.Y }

func (u *Unit) GetDistanceTo(x, y float64) float64 {
	return math.Hypot(x-u.X, y-u.Y)
}

func (u *Unit) GetDistanceToUnit(unit *Unit) float64 {
	return u.GetDistanceTo(unit.X, unit.Y)
}

type Point interface {
	GetX() float64
	GetY() float64
}

func (u *Unit) GetDistanceToPoint(p Point) float64 {
	return u.GetDistanceTo(p.GetX(), p.GetY())
}

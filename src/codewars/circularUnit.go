package codewars

type CircularUnit struct {
	Unit
	Radius float64
}

func (u *CircularUnit) GetRadius() float64            { return u.Radius }
func (u *CircularUnit) AsCircularUnit() *CircularUnit { return u }

package codewars

type Move struct {
	Action            ActionType
	Group             int
	Left              float64
	Top               float64
	Right             float64
	Bottom            float64
	X                 float64
	Y                 float64
	Angle             float64
	Max_speed         float64
	Max_angular_speed float64
	Vehicle_type      VehicleType
	Facility_id       int64
}

func NewMove() *Move {
	return &Move{
		Action:       Action_None,
		Vehicle_type: Vehicle_Unknown,
		Facility_id:  0,
	}
}

package codewars

type FacilityType int

const (
	Facility_Control_Center FacilityType = iota
	Facility_Vehicle_Factory
)

type Facility struct {
	Id                 int64
	Type               FacilityType
	OwnerPlayerId      int64
	Left               float64
	Top                float64
	CapturePoints      float64
	VehicleType        VehicleType
	ProductionProgress int
}

/*
func NewFacility() *Facility {
	return &Facility{
		Facility_id: -1,
	}
}*/

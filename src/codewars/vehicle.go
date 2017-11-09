package codewars

import (
	"errors"
)

type VehicleType int

const (
	Vehicle_Unknown VehicleType = -1
	Vehicle_Arrv    VehicleType = iota
	Vehicle_Fighter
	Vehicle_Helicopter
	Vehicle_Ifv
	Vehicle_Tank
)

type Vehicle struct {
	CircularUnit
	PlayerId                     int64
	Durability                   int
	MaxDurability                int
	MaxSpeed                     float64
	VisionRange                  float64
	SquaredVisionRange           float64
	GroundAttackRange            float64
	SquaredGroundAttackRange     float64
	AerialAttackRange            float64
	SquaredAerialAttackRange     float64
	GroundDamage                 int
	AerialDamage                 int
	GroundDefence                int
	AerialDefence                int
	AttackCooldownTicks          int
	RemainingAttackCooldownTicks int
	VehicleType                  VehicleType
	Aerial                       bool
	Selected                     bool
	Groups                       []int
}

func (v *Vehicle) update(vehicle_update *VehicleUpdate) {
	if v.Id != vehicle_update.Id {
		panic(errors.New("Vehicle ID mismatch"))
		//TODO
		//raise ValueError("Vehicle ID mismatch [actual=%s, expected=%s]." % (vehicle_update.id, f.id))
	}

	v.X = vehicle_update.X
	v.Y = vehicle_update.Y
	v.Durability = vehicle_update.Durability
	v.RemainingAttackCooldownTicks = vehicle_update.RemainingAttackCooldownTicks
	v.Selected = vehicle_update.Selected
	v.Groups = vehicle_update.Groups
}

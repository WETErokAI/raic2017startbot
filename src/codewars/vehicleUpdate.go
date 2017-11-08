package codewars

type VehicleUpdate struct {
	CircularUnit
	Durability                   int
	RemainingAttackCooldownTicks int
	Selected                     bool
	Groups                       []int
}

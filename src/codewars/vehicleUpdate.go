package codewars

type VehicleUpdate struct {
	Unit
	Durability                   int
	RemainingAttackCooldownTicks int
	Selected                     bool
	Groups                       []int
}

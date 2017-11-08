package codewars

type Player struct {
	Id                           int64
	Me                           bool
	Name                         string
	StrategyCrashed              bool
	Score                        int
	RemainingActionCooldownTicks int
}

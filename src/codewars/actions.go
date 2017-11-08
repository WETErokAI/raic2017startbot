package codewars

type ActionType int

const (
	Action_None ActionType = iota
	Action_Clear_And_Select
	Action_Add_To_Selection
	Action_Deselect
	Action_Assign
	Action_Dismiss
	Action_Disband
	Action_Move
	Action_Rotate
	Action_Setup_Vehicle_Production
)

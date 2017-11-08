package main

import (
	. "codewars"
)

type MyStrategy struct{}

func New() Strategy {
	return &MyStrategy{}
}

func (s *MyStrategy) Move(me *Player, world *World, game *Game, move *Move) {
	// put your code here
	if world.TickIndex == 0 {
		move.Action = Action_Clear_And_Select
		move.Right = world.Width
		move.Bottom = world.Height
	}

	if world.TickIndex == 1 {
		move.Action = Action_Move
		move.X = world.Width / 2.0
		move.Y = world.Height / 2.0
	}

}

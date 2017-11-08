package codewars

type Strategy interface {
	Move(me *Player, world *World, game *Game, move *Move)
}

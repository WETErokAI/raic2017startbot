package codewars

type World struct {
	TickIndex       int
	TickCount       int
	Width, Height   float64
	Players         []*Player
	NewVehicles     []*Vehicle
	VehicleUpdate   []*VehicleUpdate
	TerrainByCellXY [][]TerrainType
	WeatherByCellXY [][]WeatherType
	Facilities      []*Facility
}

func (w *World) GetMyPlayer() *Player {
	for _, player := range w.Players {
		if player.Me {
			return player
		}
	}
	return nil
}

func (w *World) GetOpponentPlayer() *Player {
	for _, player := range w.Players {
		if !player.Me {
			return player
		}
	}
	return nil
}

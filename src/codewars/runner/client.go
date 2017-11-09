package runner

import (
	"bufio"
	. "codewars"
	"encoding/binary"
	"errors"
	//"fmt"
	"net"
)

var (
	Order = binary.LittleEndian
)

type PlayerContext struct {
	Player *Player
	World  *World
}

type Client struct {
	conn net.Conn
	w    *bufio.Writer
	r    *bufio.Reader

	previousPlayers    []*Player
	previousFacilities []*Facility

	TerrainByCellXY [][]TerrainType
	WeatherByCellXY [][]WeatherType

	previousPlayerById map[int64]*Player
	prevoiusUnitById   map[int64]interface{}
}

type MessageType int

const (
	Message_Unknown MessageType = iota
	Message_GameOver
	Message_AuthToken
	Message_TeamSize
	Message_ProtoVersion
	Message_GameContext
	Message_PlayerContext
	Message_Moves
)

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn, bufio.NewWriter(conn), bufio.NewReader(conn),
		nil, nil, nil, nil, // previous
		make(map[int64]*Player), make(map[int64]interface{}), // previousById
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) WriteToken(token string) {
	c.writeOpcode(Message_AuthToken)
	c.writeString(token)
	c.flush()
}

func (c *Client) WriteProtocolVersion(ver int) {
	c.writeOpcode(Message_ProtoVersion)
	c.writeInt(ver)
	c.flush()
}

func (c *Client) WriteMovesMessage(move *Move) {
	c.writeOpcode(Message_Moves)
	c.writeMove(move)
	c.flush()
}

func (c *Client) writeMove(move *Move) {
	c.writeByte(1)

	c.writeByte(byte(move.Action))
	c.writeInt(move.Group)
	c.writeFloat64(move.Left)
	c.writeFloat64(move.Top)
	c.writeFloat64(move.Right)
	c.writeFloat64(move.Bottom)
	c.writeFloat64(move.X)
	c.writeFloat64(move.Y)
	c.writeFloat64(move.Angle)
	c.writeFloat64(move.Max_speed)
	c.writeFloat64(move.Max_angular_speed)
	c.writeByte(byte(move.Vehicle_type))
	c.writeInt64(move.Facility_id)

}

func (c *Client) ReadTeamSize() int {
	c.ensureMessageType(c.readByte(), Message_TeamSize)
	return c.readInt()
}

func (c *Client) ReadPlayerContext() *PlayerContext {
	opcode := c.readByte()
	if opcode == byte(Message_GameOver) {
		return nil
	}
	c.ensureMessageType(opcode, Message_PlayerContext)
	if !c.readBool() {
		return nil
	}
	return &PlayerContext{
		Player: c.ReadPlayer(),
		World:  c.ReadWorld(),
	}
}

func (c *Client) ReadWorld() *World {
	if !c.readBool() {
		return nil
	}
	w := World{
		TickIndex:     c.readInt(),
		TickCount:     c.readInt(),
		Width:         c.readFloat64(),
		Height:        c.readFloat64(),
		Players:       c.ReadPlayers(),
		NewVehicles:   c.ReadVehicles(),
		VehicleUpdate: c.ReadVehicleUpdates(),
	}
	if c.TerrainByCellXY == nil {
		w.TerrainByCellXY = c.ReadTerrainByCellXY()
	}
	if c.WeatherByCellXY == nil {
		w.WeatherByCellXY = c.ReadWeatherByCellXY()
	}
	w.Facilities = c.ReadFacilities()

	return &w
}

func (c *Client) ReadPlayers() []*Player {
	l := c.readInt()
	if l < 0 {
		return c.previousPlayers
	}
	r := make([]*Player, l)
	for i := range r {
		r[i] = c.ReadPlayer()
	}
	c.previousPlayers = r
	return r
}

func (c *Client) ReadPlayer() *Player {
	switch c.readByte() {
	case 0:
		return nil
	case 127:
		return c.previousPlayerById[c.readInt64()]
	default:
		p := &Player{
			Id: c.readInt64(),
			Me: c.readBool(),
			//Name:            c.readString(),
			StrategyCrashed: c.readBool(),
			Score:           c.readInt(),
			RemainingActionCooldownTicks: c.readInt(),
		}
		c.previousPlayerById[p.Id] = p
		return p
	}
}

func (c *Client) ReadVehicles() []*Vehicle {
	l := c.readInt()
	r := make([]*Vehicle, l)
	for i := range r {
		r[i] = c.ReadVehicle()
	}
	return r
}

func (c *Client) ReadVehicle() *Vehicle {
	if !c.readBool() {
		return nil
	}
	return &Vehicle{
		CircularUnit:                 c.readCircularUnit(),
		PlayerId:                     c.readInt64(),
		Durability:                   c.readInt(),
		MaxDurability:                c.readInt(),
		MaxSpeed:                     c.readFloat64(),
		VisionRange:                  c.readFloat64(),
		SquaredVisionRange:           c.readFloat64(),
		GroundAttackRange:            c.readFloat64(),
		SquaredGroundAttackRange:     c.readFloat64(),
		AerialAttackRange:            c.readFloat64(),
		SquaredAerialAttackRange:     c.readFloat64(),
		GroundDamage:                 c.readInt(),
		AerialDamage:                 c.readInt(),
		GroundDefence:                c.readInt(),
		AerialDefence:                c.readInt(),
		AttackCooldownTicks:          c.readInt(),
		RemainingAttackCooldownTicks: c.readInt(),
		VehicleType:                  VehicleType(c.readByte()),
		Aerial:                       c.readBool(),
		Selected:                     c.readBool(),
		Groups:                       c.readIntArray(),
	}
}

func (c *Client) ReadVehicleUpdates() []*VehicleUpdate {
	l := c.readInt()
	r := make([]*VehicleUpdate, l)
	for i := range r {
		r[i] = c.ReadVehicleUpdate()
	}
	return r
}

func (c *Client) ReadVehicleUpdate() *VehicleUpdate {
	if !c.readBool() {
		return nil
	}
	return &VehicleUpdate{
		Unit: Unit{
			Id: c.readInt64(),
			X:  c.readFloat64(),
			Y:  c.readFloat64(),
		},
		Durability:                   c.readInt(),
		RemainingAttackCooldownTicks: c.readInt(),
		Selected:                     c.readBool(),
		Groups:                       c.readIntArray(),
	}
}

func (c *Client) ReadTerrainByCellXY() [][]TerrainType {
	countX := c.readInt()
	rX := make([][]TerrainType, countX)
	for i := range rX {

		countY := c.readInt()

		rY := make([]TerrainType, countY)
		for i := range rY {
			rY[i] = TerrainType(c.readByte())
		}
		rX[i] = rY
	}

	c.TerrainByCellXY = rX
	return c.TerrainByCellXY
}

func (c *Client) ReadWeatherByCellXY() [][]WeatherType {
	countX := c.readInt()
	rX := make([][]WeatherType, countX)
	for i := range rX {

		countY := c.readInt()

		rY := make([]WeatherType, countY)
		for i := range rY {
			rY[i] = WeatherType(c.readByte())
		}
		rX[i] = rY
	}

	c.WeatherByCellXY = rX

	return c.WeatherByCellXY
}

func (c *Client) ReadFacilities() []*Facility {
	l := c.readInt()
	if l < 0 {
		return c.previousFacilities
	}
	f := make([]*Facility, l)
	for i := range f {
		f[i] = c.ReadFacility()
	}
	c.previousFacilities = f
	return f
}

func (c *Client) ReadFacility() *Facility {
	if !c.readBool() {
		return nil
	}
	return &Facility{
		Id:                 c.readInt64(),
		Type:               FacilityType(c.readByte()),
		OwnerPlayerId:      c.readInt64(),
		Left:               c.readFloat64(),
		Top:                c.readFloat64(),
		CapturePoints:      c.readFloat64(),
		VehicleType:        VehicleType(c.readByte()),
		ProductionProgress: c.readInt(),
	}
}

func (c *Client) readCircularUnit() CircularUnit {
	return CircularUnit{
		Unit: Unit{
			Id: c.readInt64(),
			X:  c.readFloat64(),
			Y:  c.readFloat64(),
		},
		Radius: c.readFloat64(),
	}
}

func (c *Client) ReadGameContext() *Game {
	c.ensureMessageType(c.readByte(), Message_GameContext)
	if !c.readBool() {
		return nil
	}
	return &Game{
		RandomSeed:                             c.readInt64(),
		TickCount:                              c.readInt(),
		WorldWidth:                             c.readFloat64(),
		WorldHeight:                            c.readFloat64(),
		FogOfWarEnabled:                        c.readBool(),
		VictoryScore:                           c.readInt(),
		FacilityCaptureScore:                   c.readInt(),
		VehicleEliminationScore:                c.readInt(),
		ActionDetectionInterval:                c.readInt(),
		BaseActionCount:                        c.readInt(),
		AdditionalActionCountPerControlCenter:  c.readInt(),
		MaxUnitGroup:                           c.readInt(),
		TerrainWeatherMapColumnCount:           c.readInt(),
		TerrainWeatherMapRowCount:              c.readInt(),
		PlainTerrainVisionFactor:               c.readFloat64(),
		PlainTerrainStealthFactor:              c.readFloat64(),
		PlainTerrainSpeedFactor:                c.readFloat64(),
		SwampTerrainVisionFactor:               c.readFloat64(),
		SwampTerrainStealthFactor:              c.readFloat64(),
		SwampTerrainSpeedFactor:                c.readFloat64(),
		ForestTerrainVisionFactor:              c.readFloat64(),
		ForestTerrainStealthFactor:             c.readFloat64(),
		ForestTerrainSpeedFactor:               c.readFloat64(),
		ClearWeatherVisionFactor:               c.readFloat64(),
		ClearWeatherStealthFactor:              c.readFloat64(),
		ClearWeatherSpeedFactor:                c.readFloat64(),
		CloudWeatherVisionFactor:               c.readFloat64(),
		CloudWeatherStealthFactor:              c.readFloat64(),
		CloudWeatherSpeedFactor:                c.readFloat64(),
		RainWeatherVisionFactor:                c.readFloat64(),
		RainWeatherStealthFactor:               c.readFloat64(),
		RainWeatherSpeedFactor:                 c.readFloat64(),
		VehicleRadius:                          c.readFloat64(),
		TankDurability:                         c.readInt(),
		TankSpeed:                              c.readFloat64(),
		TankVisionRange:                        c.readFloat64(),
		TankGroundAttackRange:                  c.readFloat64(),
		TankAerialAttackRange:                  c.readFloat64(),
		TankGroundDamage:                       c.readInt(),
		TankAerialDamage:                       c.readInt(),
		TankGroundDefence:                      c.readInt(),
		TankAerialDefence:                      c.readInt(),
		TankAttackCooldownTicks:                c.readInt(),
		TankProductionCost:                     c.readInt(),
		IfvDurability:                          c.readInt(),
		IfvSpeed:                               c.readFloat64(),
		IfvVisionRange:                         c.readFloat64(),
		IfvGroundAttackRange:                   c.readFloat64(),
		IfvAerialAttackRange:                   c.readFloat64(),
		IfvGroundDamage:                        c.readInt(),
		IfvAerialDamage:                        c.readInt(),
		IfvGroundDefence:                       c.readInt(),
		IfvAerialDefence:                       c.readInt(),
		IfvAttackCooldownTicks:                 c.readInt(),
		IfvProductionCost:                      c.readInt(),
		ArrvDurability:                         c.readInt(),
		ArrvSpeed:                              c.readFloat64(),
		ArrvVisionRange:                        c.readFloat64(),
		ArrvGroundDefence:                      c.readInt(),
		ArrvAerialDefence:                      c.readInt(),
		ArrvProductionCost:                     c.readInt(),
		ArrvRepairRange:                        c.readFloat64(),
		ArrvRepairSpeed:                        c.readFloat64(),
		HelicopterDurability:                   c.readInt(),
		HelicopterSpeed:                        c.readFloat64(),
		HelicopterVisionRange:                  c.readFloat64(),
		HelicopterGroundAttackRange:            c.readFloat64(),
		HelicopterAerialAttackRange:            c.readFloat64(),
		HelicopterGroundDamage:                 c.readInt(),
		HelicopterAerialDamage:                 c.readInt(),
		HelicopterGroundDefence:                c.readInt(),
		HelicopterAerialDefence:                c.readInt(),
		HelicopterAttackCooldownTicks:          c.readInt(),
		HelicopterProductionCost:               c.readInt(),
		FighterDurability:                      c.readInt(),
		FighterSpeed:                           c.readFloat64(),
		FighterVisionRange:                     c.readFloat64(),
		FighterGroundAttackRange:               c.readFloat64(),
		FighterAerialAttackRange:               c.readFloat64(),
		FighterGroundDamage:                    c.readInt(),
		FighterAerialDamage:                    c.readInt(),
		FighterGroundDefence:                   c.readInt(),
		FighterAerialDefence:                   c.readInt(),
		FighterAttackCooldownTicks:             c.readInt(),
		FighterProductionCost:                  c.readInt(),
		MaxFacilityCapturePoints:               c.readFloat64(),
		FacilityCapturePointsPerVehiclePerTick: c.readFloat64(),
		FacilityWidth:                          c.readFloat64(),
		FacilityHeight:                         c.readFloat64(),
	}
}

func (c *Client) readIntArray() []int {
	count := c.readInt()
	r := make([]int, count)
	for i := range r {
		r[i] = c.readInt()
	}
	return r
}

func (c *Client) readIntArray2D() [][]int {
	count := c.readInt()
	r := make([][]int, count)
	for i := range r {
		r[i] = c.readIntArray()
	}
	return r
}

func (c *Client) readInt() int {
	var v int32
	if err := binary.Read(c.r, Order, &v); err != nil {
		panic(err)
	}
	return int(v)
}

func (c *Client) readInt64() int64 {
	var v int64
	if err := binary.Read(c.r, Order, &v); err != nil {
		panic(err)
	}
	return v
}

func (c *Client) readFloat64() float64 {
	var v float64
	if err := binary.Read(c.r, Order, &v); err != nil {
		panic(err)
	}
	return v
}

func (c *Client) readBool() bool {
	return c.readByte() != 0
}

func (c *Client) readByte() byte {
	b, err := c.r.ReadByte()
	if err != nil {
		panic(err)
	}
	return b
}

func (c *Client) readBytes() []byte {
	l := c.readInt()
	r := make([]byte, l)
	for i := range r {
		r[i] = c.readByte()
	}
	return r
}

func (c *Client) readString() string {
	return string(c.readBytes())
}

func (c *Client) ensureMessageType(v byte, m MessageType) {
	if v != byte(m) {
		panic(errors.New("unexpected message"))
	}
}

func (c *Client) writeOpcode(m MessageType) {
	c.writeByte(byte(m))
}

func (c *Client) writeInt(v int) {
	if err := binary.Write(c.w, Order, int32(v)); err != nil {
		panic(err)
	}
}

func (c *Client) writeFloat64(v float64) {
	if err := binary.Write(c.w, Order, v); err != nil {
		panic(err)
	}
}

func (c *Client) writeInt64(v int64) {
	if err := binary.Write(c.w, Order, v); err != nil {
		panic(err)
	}
}

func (c *Client) writeByte(v byte) {
	if err := c.w.WriteByte(v); err != nil {
		panic(err)
	}
}

func (c *Client) writeBytes(v []byte) {
	c.writeInt(len(v))
	if _, err := c.w.Write(v); err != nil {
		panic(err)
	}
}

func (c *Client) writeString(v string) {
	c.writeInt(len(v))
	if _, err := c.w.WriteString(v); err != nil {
		panic(err)
	}
}

func (c *Client) flush() {
	if err := c.w.Flush(); err != nil {
		panic(err)
	}
}

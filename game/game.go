// Copyright (c) 2012 by Lecture Hall Games Authors.
// All source files are distributed under the Simplified BSD License.
package game

import (
	"errors"
	"github.com/banthar/Go-SDL/mixer"
	"github.com/banthar/Go-SDL/sdl"
	"log"
	"time"
)

const Ghost Player = Player(1)
const Human Player = Player(0)
const GhostOrHuman Player = Player(-1)

func FillRect(x0, y0, x1, y1 int, pos []MapPosition) []MapPosition {
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			pos = append(pos, MapPosition{x, y})
		}
	}

	return pos
}

func NewMapConfig() *MapConfig {
	walls := []MapPosition{
		MapPosition{0, 0},
		MapPosition{0, 1},
		MapPosition{0, 2},
		MapPosition{0, 3},
		MapPosition{0, 4},
		MapPosition{0, 5},
		MapPosition{0, 6},
		MapPosition{0, 7},
		MapPosition{0, 8},
		MapPosition{0, 9},
		MapPosition{0, 10},
		MapPosition{0, 11},
		MapPosition{0, 12},
		MapPosition{0, 13},
		MapPosition{0, 14},
		MapPosition{0, 15},
		MapPosition{0, 16},

		MapPosition{1, 0},
		MapPosition{2, 0},
		MapPosition{3, 0},
		MapPosition{4, 0},
		MapPosition{5, 0},
		MapPosition{6, 0},
		MapPosition{7, 0},
		MapPosition{8, 0},
		MapPosition{9, 0},
		MapPosition{10, 0},
		MapPosition{11, 0},
		MapPosition{12, 0},
		MapPosition{13, 0},
		MapPosition{14, 0},

		MapPosition{14, 0},
		MapPosition{14, 1},
		MapPosition{14, 2},
		MapPosition{14, 3},
		MapPosition{14, 4},
		MapPosition{14, 5},
		MapPosition{14, 6},
		MapPosition{14, 7},
		MapPosition{14, 8},
		MapPosition{14, 9},
		MapPosition{14, 10},
		MapPosition{14, 11},
		MapPosition{14, 12},
		MapPosition{14, 13},
		MapPosition{14, 14},
		MapPosition{14, 15},
		MapPosition{14, 16},

		MapPosition{1, 16},
		MapPosition{2, 16},
		MapPosition{3, 16},
		MapPosition{4, 16},
		MapPosition{5, 16},
		MapPosition{6, 16},
		MapPosition{7, 16},
		MapPosition{8, 16},
		MapPosition{9, 16},
		MapPosition{10, 16},
		MapPosition{11, 16},
		MapPosition{12, 16},
		MapPosition{13, 16},

		MapPosition{1, 1},
		MapPosition{1, 2},
		MapPosition{1, 4},
		MapPosition{1, 5},
		MapPosition{1, 6},
		MapPosition{2, 6},
		MapPosition{3, 6},
		MapPosition{4, 6},
		MapPosition{5, 6},
		MapPosition{6, 6},
		MapPosition{7, 6},
		MapPosition{8, 6},
		MapPosition{8, 4},
		MapPosition{8, 3},
		MapPosition{8, 2},
		MapPosition{8, 1},

		MapPosition{8, 7},
		MapPosition{8, 8},
		MapPosition{8, 9},
		MapPosition{8, 10},
		MapPosition{8, 12},
		MapPosition{8, 13},

		MapPosition{8, 13},
		MapPosition{7, 13},
		MapPosition{6, 13},
		MapPosition{5, 13},
		MapPosition{3, 13},
		MapPosition{1, 13},

		MapPosition{1, 12},
		MapPosition{1, 11},

		MapPosition{3, 12},
		MapPosition{3, 11},

		MapPosition{5, 12},
		MapPosition{5, 11},

		MapPosition{6, 14},
		MapPosition{6, 15},

		MapPosition{9, 10},
		MapPosition{10, 10},

		MapPosition{10, 11},
		MapPosition{10, 12},
		MapPosition{10, 13},
		MapPosition{10, 15},

		MapPosition{11, 13},
		MapPosition{13, 13},
		MapPosition{12, 15},

		MapPosition{9, 8},
		MapPosition{11, 8},
		MapPosition{12, 8},
		MapPosition{13, 8},

		MapPosition{5, 5},
		MapPosition{5, 4},
		MapPosition{5, 3},

		MapPosition{7, 3},

		MapPosition{3, 14},
		MapPosition{3, 15},
	}

	triggers := []CfgTriggerData{
		// triggerID doorID bannID
		NewCfgTriggerData(1, 2, -1, NewMapPosition(4, 15), DirWest, Human, GhostOrHuman, -1),
		NewCfgTriggerData(2, 1, -1, NewMapPosition(2, 15), DirWest, Ghost, GhostOrHuman, -1),
		// boulder drop
		NewCfgTriggerData(3, -1, -1, NewMapPosition(7, 12), DirWest, Human, GhostOrHuman, 2), // x1
		// NewCfgTriggerData(3, 3, -1, NewMapPosition(7, 12), DirEast, Human, GhostOrHuman, 2), // x1

		NewCfgTriggerData(5, 6, -1, NewMapPosition(11, 12), DirWest, Human, Ghost, -1),        // g4
		NewCfgTriggerData(6, 5, -1, NewMapPosition(13, 11), DirWest, Ghost, GhostOrHuman, -1), // g5
		NewCfgTriggerData(7, 7, -1, NewMapPosition(13, 5), DirWest, Human, GhostOrHuman, -1),

		// inactive
		NewCfgTriggerData(8, -1, -1, NewMapPosition(13, 1), DirWest, Ghost, Human, -1),

		NewCfgTriggerData(9, -1, -1, NewMapPosition(2, 1), DirWest, Human, Ghost, -1),
		NewCfgTriggerData(10, -1, -1, NewMapPosition(2, 5), DirWest, GhostOrHuman, Human, -1),
		NewCfgTriggerData(11, -1, -1, NewMapPosition(7, 1), DirWest, Human, GhostOrHuman, -1),

		// mid down
		// NewCfgTriggerData(12, -1, -1, NewMapPosition(8, 15), DirSouth, Ghost, Human),
		NewCfgTriggerData(12, 4, -1, NewMapPosition(8, 15), DirWest, Ghost, Human, -1),
	}

	plates := []CfgPlateData{
		NewCfgPlateData(1, -1, -1, NewMapPosition(1, 10)),
		NewCfgPlateData(1, 3, -1, NewMapPosition(7, 7)),
	}

	doors := []CfgDoorData{
		NewCfgDoorData(1, 3, NewMapPosition(2, 13)),
		NewCfgDoorData(2, 3, NewMapPosition(4, 13)),
		NewCfgDoorData(3, 4, NewMapPosition(8, 11)),
		NewCfgDoorData(4, 7, NewMapPosition(10, 14)), /// t4
		NewCfgDoorData(5, 5, NewMapPosition(12, 13)), // t5
		NewCfgDoorData(6, 6, NewMapPosition(10, 8)),  // t6
		NewCfgDoorData(7, 8, NewMapPosition(8, 5)),
		NewCfgDoorData(8, 9, NewMapPosition(6, 3)),
		NewCfgDoorData(9, -1, NewMapPosition(1, 3)),
	}

	boulders := []CfgBoulderData{
		NewCfgBoulderData(1, true, NewMapPosition(4, 12)),
		NewCfgBoulderData(2, false, NewMapPosition(3, 8)),
	}

	bannWalls := []CfgBannWallData{
		NewCfgBannWallData(1, NewMapPosition(3, 1), 0),
		NewCfgBannWallData(2, NewMapPosition(4, 1), 1),
		NewCfgBannWallData(3, NewMapPosition(3, 2), 2),
		NewCfgBannWallData(4, NewMapPosition(4, 2), 3),

		NewCfgBannWallData(5, NewMapPosition(5, 1), 0),
		NewCfgBannWallData(6, NewMapPosition(6, 1), 1),
		NewCfgBannWallData(7, NewMapPosition(5, 2), 2),
		NewCfgBannWallData(8, NewMapPosition(6, 2), 3),

		NewCfgBannWallData(9, NewMapPosition(3, 3), 0),
		NewCfgBannWallData(10, NewMapPosition(4, 3), 1),
		NewCfgBannWallData(11, NewMapPosition(3, 4), 2),
		NewCfgBannWallData(13, NewMapPosition(4, 4), 3),
	}

	room3 := make([]MapPosition, 0)
	room3 = FillRect(0, 6, 8, 13, room3)

	room4 := make([]MapPosition, 0)
	room4 = FillRect(6, 13, 10, 16, room4)
	room4 = FillRect(8, 10, 10, 13, room4)

	room5 := make([]MapPosition, 0)
	room5 = FillRect(8, 8, 14, 10, room5)
	room5 = FillRect(10, 10, 14, 13, room5)

	room7 := make([]MapPosition, 0)
	room7 = FillRect(10, 13, 14, 16, room7)

	room6 := make([]MapPosition, 0)
	room6 = FillRect(8, 0, 14, 8, room6)

	room8 := make([]MapPosition, 0)
	room8 = FillRect(5, 3, 8, 6, room8)

	room9 := make([]MapPosition, 0)
	room9 = FillRect(0, 0, 5, 6, room9)
	room9 = FillRect(5, 0, 8, 3, room9)

	rooms := []CfgRoomData{
		NewCfgRoomData(1, []MapPosition{
			MapPosition{0, 16},
			MapPosition{1, 16},
			MapPosition{2, 16},
			MapPosition{3, 16},

			MapPosition{0, 15},
			MapPosition{1, 15},
			MapPosition{2, 15},
			MapPosition{3, 15},

			MapPosition{0, 14},
			MapPosition{1, 14},
			MapPosition{2, 14},
			MapPosition{3, 14},

			MapPosition{0, 13},
			MapPosition{1, 13},
			MapPosition{2, 13},
			MapPosition{3, 13},
		}),

		NewCfgRoomData(2, []MapPosition{
			MapPosition{0 + 3, 16},
			MapPosition{1 + 3, 16},
			MapPosition{2 + 3, 16},
			MapPosition{3 + 3, 16},

			MapPosition{0 + 3, 15},
			MapPosition{1 + 3, 15},
			MapPosition{2 + 3, 15},
			MapPosition{3 + 3, 15},

			MapPosition{0 + 3, 14},
			MapPosition{1 + 3, 14},
			MapPosition{2 + 3, 14},
			MapPosition{3 + 3, 14},

			MapPosition{0 + 3, 13},
			MapPosition{1 + 3, 13},
			MapPosition{2 + 3, 13},
			MapPosition{3 + 3, 13},
		}),

		NewCfgRoomData(3, room3),

		NewCfgRoomData(4, room4),
		NewCfgRoomData(5, room5),
		NewCfgRoomData(6, room6),
		NewCfgRoomData(7, room7),
		NewCfgRoomData(8, room8),
		NewCfgRoomData(9, room9),
	}

	doorMu := mixer.LoadMUS("data/door.ogg")
	if doorMu == nil {
		log.Println(sdl.GetError())
	}

	triggerMu1 := mixer.LoadMUS("data/trigger1.wav")
	if triggerMu1 == nil {
		log.Println(sdl.GetError())
	}

	triggerMu2 := mixer.LoadMUS("data/trigger2.wav")
	if triggerMu2 == nil {
		log.Println(sdl.GetError())
	}

	return &MapConfig{
		playerStartPos: []MapPosition{
			MapPosition{5, 15},
			MapPosition{1, 15},
		},
		playerStartLook: []Direction{
			DirEast,
			DirWest,
		},

		walkTime:   200 * time.Millisecond,
		rollTime:   200 * time.Millisecond,
		actionTime: 200 * time.Millisecond,

		mapWidth:  15,
		mapHeight: 17,

		walls: walls,

		triggerData:  triggers,
		plateData:    plates,
		doorData:     doors,
		boulderData:  boulders,
		bannWallData: bannWalls,
		roomData:     rooms,

		rooms:     make(map[RoomID]*Room),
		plates:    make(map[PlateID]*Plate),
		triggers:  make(map[TriggerID]*Trigger),
		doors:     make(map[DoorID]*Door),
		boulders:  make(map[BoulderID]*Boulder),
		bannWalls: make(map[BannWallID]*BannWall),

		doorMusic:     doorMu,
		trigger1Music: triggerMu1,
		trigger2Music: triggerMu2,
	}
}

type CellType int

const (
	CellTypeEmpty CellType = iota
	CellTypeWall
	CellTypeDoor
)

var GlobalConfig *MapConfig = NewMapConfig()

type RoomID int
type TriggerID int
type DoorID int
type BoulderID int
type BannWallID int
type PlateID int

type CfgTriggerData struct {
	id             TriggerID
	targetDoor     DoorID
	targetBannWall BannWallID
	pos            MapPosition
	dir            Direction
	canTrigger     Player
	canVis         Player
	targetBoulder  BoulderID
}

type CfgPlateData struct {
	id             PlateID
	targetDoor     DoorID
	targetBannWall BannWallID
	pos            MapPosition
}

type CfgDoorData struct {
	id         DoorID
	pos        MapPosition
	targetRoom RoomID
}

type CfgRoomData struct {
	id    RoomID
	cells []MapPosition
}

type CfgBoulderData struct {
	id      BoulderID
	pos     MapPosition
	spawned bool
}

type CfgBannWallData struct {
	id           BannWallID
	pos          MapPosition
	bannWallType int
}

func NewCfgBannWallData(id BannWallID, pos MapPosition, bannWallType int) CfgBannWallData {
	return CfgBannWallData{
		id:           id,
		pos:          pos,
		bannWallType: bannWallType,
	}
}

func NewCfgBoulderData(id BoulderID, spawned bool, pos MapPosition) CfgBoulderData {
	return CfgBoulderData{
		id:      id,
		pos:     pos,
		spawned: spawned,
	}
}

func NewCfgTriggerData(id TriggerID, door DoorID, bannWall BannWallID, pos MapPosition, dir Direction, canTrigger Player, canVis Player, boulder BoulderID) CfgTriggerData {
	return CfgTriggerData{
		id:             id,
		targetDoor:     door,
		pos:            pos,
		dir:            dir,
		targetBannWall: bannWall,
		canTrigger:     canTrigger,
		canVis:         canVis,
		targetBoulder:  boulder,
	}
}

func NewCfgPlateData(id PlateID, door DoorID, bannWall BannWallID, pos MapPosition) CfgPlateData {
	return CfgPlateData{
		id:             id,
		targetDoor:     door,
		targetBannWall: bannWall,
		pos:            pos,
	}
}

func NewCfgDoorData(id DoorID, room RoomID, pos MapPosition) CfgDoorData {
	return CfgDoorData{
		id:         id,
		pos:        pos,
		targetRoom: room,
	}
}

func NewCfgRoomData(id RoomID, cells []MapPosition) CfgRoomData {
	return CfgRoomData{
		id:    id,
		cells: cells,
	}
}

type MapConfig struct {
	playerStartPos  []MapPosition
	playerStartLook []Direction

	walkTime   time.Duration
	rollTime   time.Duration
	actionTime time.Duration

	walls     []MapPosition
	mapWidth  int
	mapHeight int

	triggerData  []CfgTriggerData
	plateData    []CfgPlateData
	doorData     []CfgDoorData
	roomData     []CfgRoomData
	boulderData  []CfgBoulderData
	bannWallData []CfgBannWallData

	rooms     map[RoomID]*Room
	plates    map[PlateID]*Plate
	triggers  map[TriggerID]*Trigger
	doors     map[DoorID]*Door
	boulders  map[BoulderID]*Boulder
	bannWalls map[BannWallID]*BannWall

	doorMusic     *mixer.Music
	trigger1Music *mixer.Music
	trigger2Music *mixer.Music
}

func BuildGame(g *Game) {
	// g, _ := NewGame()

	for _, pos := range GlobalConfig.walls {
		g.SetWall(pos)
	}

	for _, doorData := range GlobalConfig.doorData {
		door := g.SetDoor(doorData.pos)
		GlobalConfig.doors[doorData.id] = door
	}

	for _, boulderData := range GlobalConfig.boulderData {
		boulder := g.SetBoulder(boulderData.pos, boulderData.spawned)
		GlobalConfig.boulders[boulderData.id] = boulder
	}

	for _, plateData := range GlobalConfig.plateData {
		plate := g.SetPlate(plateData.pos)
		GlobalConfig.plates[plateData.id] = plate
	}

	for _, triggerData := range GlobalConfig.triggerData {
		var boulder *Boulder = nil
		if triggerData.targetBoulder >= 0 {
			boulder = GlobalConfig.boulders[triggerData.targetBoulder]
		}

		log.Println("Init boulder", boulder)

		trigger := g.SetTrigger(triggerData.pos, triggerData.dir, triggerData.canTrigger, triggerData.canVis, boulder)
		GlobalConfig.triggers[triggerData.id] = trigger
	}

	for _, roomData := range GlobalConfig.roomData {
		room := g.NewRoom(roomData.cells)
		GlobalConfig.rooms[roomData.id] = room
	}

	for _, bannWallData := range GlobalConfig.bannWallData {
		bannWall := g.SetBannWall(bannWallData.pos, bannWallData.bannWallType)
		GlobalConfig.bannWalls[bannWallData.id] = bannWall
	}

	GlobalConfig.rooms[1].isVisible = true
	GlobalConfig.rooms[2].isVisible = true

	ConnectEverything(g)
}

func ConnectEverything(g *Game) {
	for _, triggerData := range GlobalConfig.triggerData {
		trigger := GlobalConfig.triggers[triggerData.id]
		if triggerData.targetDoor > 0 {
			targetDoor := GlobalConfig.doors[triggerData.targetDoor]
			trigger.linkedDoor = targetDoor
		}

		if triggerData.targetBannWall > 0 {
			targetBannWall := GlobalConfig.bannWalls[triggerData.targetBannWall]
			trigger.linkedBannWall = targetBannWall
		}
	}

	for _, plateData := range GlobalConfig.plateData {
		plate := GlobalConfig.plates[plateData.id]
		if plateData.targetDoor > 0 {
			targetDoor := GlobalConfig.doors[plateData.targetDoor]
			plate.linkedDoor = targetDoor
		}

		if plateData.targetBannWall > 0 {
			targetBannWall := GlobalConfig.bannWalls[plateData.targetBannWall]
			plate.linkedBannWall = targetBannWall
		}
	}

	for _, doorData := range GlobalConfig.doorData {
		targetRoom := GlobalConfig.rooms[doorData.targetRoom]
		door := GlobalConfig.doors[doorData.id]
		door.linkedRoom = targetRoom
	}
}

// func (g *Game) CellIsVisible(pos MapPosition) bool {
// return true
// }

type Player int

type Room struct {
	isVisible bool
	cells     []MapPosition
}

func (g *Game) NewRoom(cells []MapPosition) *Room {
	r := &Room{
		isVisible: false,
		cells:     cells,
	}

	g.rooms = append(g.rooms, r)

	return r
}

type Door struct {
	isOpen     bool
	linkedRoom *Room
}

func (g *Game) MakeRoomVisible(room *Room) {
	room.isVisible = true
	for _, player := range g.players {
		for _, cellPos := range room.cells {
			g.playerVis[player].visCell[cellPos] = true
		}
	}
}

func (g *Game) ToggleDoor(d *Door) {
	if d.isOpen {
		// d.linkedRoom.isVisible = false
		d.isOpen = false
	} else {
		if GlobalConfig.doorMusic != nil {
			GlobalConfig.doorMusic.PlayMusic(-1)
		}
		g.MakeRoomVisible(d.linkedRoom)
		d.isOpen = true
	}
}

func (d *Door) IsOpen() bool {
	return d.isOpen
}
func (d *Door) IsClosed() bool {
	return !d.IsOpen()
}

func NewDoor() *Door {
	return &Door{
		isOpen:     false,
		linkedRoom: nil,
	}
}

type Trigger struct {
	isActive       bool
	staysActive    time.Duration
	linkedDoor     *Door
	linkedBannWall *BannWall
	needsTrigger   *Trigger

	canTrigger   Player
	canVis       Player
	spawnBoulder *Boulder
}

func (t *Trigger) IsActive() bool {
	return t.isActive
}

func NewTrigger(canTrigger, canVis Player, boulder *Boulder) *Trigger {
	return &Trigger{
		isActive:       false,
		staysActive:    0,
		linkedDoor:     nil,
		linkedBannWall: nil,
		needsTrigger:   nil,
		canTrigger:     canTrigger,
		canVis:         canVis,
		spawnBoulder:   boulder,
	}
}

type Boulder struct {
	active bool
}

func NewBoulder(active bool) *Boulder {
	return &Boulder{
		active: active,
	}
}

func (b *Boulder) IsActive() bool {
	return b.active
}

type Plate struct {
	isActive       bool
	linkedBannWall *BannWall
	linkedDoor     *Door
}

func NewPlate() *Plate {
	return &Plate{
		isActive:       false,
		linkedBannWall: nil,
		linkedDoor:     nil,
	}
}

type BannWall struct {
	isActive     bool
	bannWallType int
}

func (bw *BannWall) Type() int {
	return bw.bannWallType
}

func NewBannWall(bannWallType int) *BannWall {
	return &BannWall{
		isActive:     true,
		bannWallType: bannWallType,
	}
}

type MapPosition struct {
	x int
	y int
}

func NewMapPosition(x int, y int) MapPosition {
	return MapPosition{x: x, y: y}
}

func (pos MapPosition) X() int {
	return pos.x
}

func (pos MapPosition) Y() int {
	return pos.y
}

func (mp MapPosition) Neighbor(direction Direction) MapPosition {
	switch direction {
	case DirNorth:
		return MapPosition{mp.x, mp.y - 1}
	case DirWest:
		return MapPosition{mp.x - 1, mp.y}
	case DirSouth:
		return MapPosition{mp.x, mp.y + 1}
	case DirEast:
		return MapPosition{mp.x + 1, mp.y}
	}
	log.Fatal("Not reached")
	return MapPosition{-1, -1}
}

type Position struct {
	x float32
	y float32
}

func (p Position) X() float32 {
	return p.x
}

func (p Position) Y() float32 {
	return p.y
}

type Cell struct {
	accessibleTriggers map[Direction]*Trigger
	isWall             bool
}

func (c *Cell) DirOfTrigger(t *Trigger) (Direction, error) {
	for dir, trigger := range c.accessibleTriggers {
		if trigger == t {
			return dir, nil
		}
	}

	return DirNorth, errors.New("No trigger")
}

func (c *Cell) TriggerByDir(dir Direction) *Trigger {
	if trigger, ok := c.accessibleTriggers[dir]; ok {
		return trigger
	}

	return nil
}

func NewCell() *Cell {
	return &Cell{
		accessibleTriggers: make(map[Direction]*Trigger),
		isWall:             false,
	}
}

func (c *Cell) IsWall() bool {
	return c.isWall
}

func (c *Cell) SetWall() {
	c.isWall = true
}

// Index of type Position
type Map struct {
	cells [][]*Cell
}

func (m *Map) Cell(pos MapPosition) *Cell {
	return m.cells[pos.y][pos.x]
}

func (m *Map) NeighborOfCell(pos MapPosition, direction Direction) *Cell {
	return m.Cell(pos.Neighbor(direction))
}

func (g *Game) SetTrigger(pos MapPosition, direction Direction, canTrigger, canVis Player, boulder *Boulder) *Trigger {
	t := NewTrigger(canTrigger, canVis, boulder)
	g.triggers[pos] = t
	cell := g.gameMap.Cell(pos)
	cell.accessibleTriggers[direction] = t

	return t
}

type Direction int

const (
	DirNorth Direction = iota
	DirWest
	DirSouth
	DirEast
)

func Dirs() []Direction {
	dirs := make([]Direction, 4)
	dirs[0] = DirNorth
	dirs[1] = DirEast
	dirs[2] = DirSouth
	dirs[3] = DirWest
	return dirs
}

type PlayerState struct {
	mapPos  MapPosition
	looksIn Direction
}

func NewPlayerState(pos MapPosition, looksIn Direction) *PlayerState {
	return &PlayerState{
		mapPos:  pos,
		looksIn: looksIn,
	}
}

type DoorTransition struct {
	door    *Door
	dtime   time.Duration
	toState bool
}

type Transition interface {
	Update(time.Duration)
	IsFinished() bool
	UpdateGameState(*Game)
}

type AnimTransition interface {
	Transition
	Frame() int
}

type MoveableTransition interface {
	AnimTransition
	OriginPos() MapPosition
	TargetPos() MapPosition
	InterpPos() Position
}

func (g *Game) PlayerIsWalking(player Player) bool {
	_, ok := g.playerMoveTransition[player]
	return ok
}

// func (g *Game) BoulderIsMoving(boulder *Boulder) bool {
// 	_, ok := g.boulderTransition[boulder]
// 	return ok
// }

// func (g *Game) BoulderMoveFrame(boulder *Boulder) int {
// 	return 0
// }

func (g *Game) PlayerActionFrame(player Player) int {
	if g.PlayerDoesAction(player) {
		return g.playerActionTransition[player].Frame()
	}

	return 0
}

func (g *Game) PlayerWalkFrame(player Player) int {
	if g.PlayerIsWalking(player) {
		return g.playerMoveTransition[player].Frame()
	}
	return 0
}

func (g *Game) PerformPlayerAction(player Player, action ActionType) error {
	switch action {
	case ActionMoveNorth:
		log.Println("Move player")
		return g.PlayerMove(player, DirNorth)
	case ActionMoveEast:
		log.Println("Move player")
		return g.PlayerMove(player, DirEast)
	case ActionMoveSouth:
		log.Println("Move player")
		return g.PlayerMove(player, DirSouth)
	case ActionMoveWest:
		log.Println("Move player")
		return g.PlayerMove(player, DirWest)

	case ActionLookNorth:
		log.Println("Player look")
		return g.PlayerLookIn(player, DirNorth)
	case ActionLookEast:
		log.Println("Player look")
		return g.PlayerLookIn(player, DirEast)
	case ActionLookSouth:
		log.Println("Player look")
		return g.PlayerLookIn(player, DirSouth)
	case ActionLookWest:
		log.Println("Player look")
		return g.PlayerLookIn(player, DirWest)

	case ActionAction:
		return g.PlayerAction(player, g.playerState[player].looksIn) // action in look direction
	case ActionToggleVisibility:
		return g.PlayerSetOtherVisible(player)
	case ActionNoAction:
		return nil
	}
	log.Fatal("Not reached")
	return nil
}

type VisStateTransition struct {
	dtime  time.Duration
	player Player
}

type VisDelay struct {
	dtime  time.Duration
	player Player
}

func NewVisDelay(player Player) *VisDelay {
	return &VisDelay{
		dtime:  0,
		player: player,
	}
}

func (vt *VisDelay) Update(dt time.Duration) {
	vt.dtime += dt
}

func (vt *VisDelay) IsFinished() bool {
	return vt.dtime >= 5000*time.Millisecond
}

func (vt *VisDelay) UpdateGameState(g *Game) {
	g.playerVis[vt.player].visPlayer = false
}

func NewVisStateTransition(player Player) *VisStateTransition {
	return &VisStateTransition{
		dtime:  0,
		player: player,
	}
}

func (vt *VisStateTransition) Update(dt time.Duration) {
	vt.dtime += dt
}

func (vt *VisStateTransition) IsFinished() bool {
	return vt.dtime >= 500*time.Millisecond
}

func (vt *VisStateTransition) UpdateGameState(g *Game) {
	g.playerVis[vt.player].visPlayer = false
}

func (g *Game) PlayerSetOtherVisible(player Player) error {
	// var otherPlayer Player
	// if player == Human {
	// 	otherPlayer = Ghost
	// } else {
	// 	otherPlayer = Human
	// }

	if _, ok := g.playerVisTransition[player]; ok {
		return errors.New("Already visible")
	}

	if _, ok := g.visDelay[player]; ok {
		return errors.New("Delay")
	}

	g.playerVisTransition[player] = NewVisStateTransition(player)

	g.playerVis[player].visPlayer = true

	return nil
}

func (g *Game) PlayerLookIn(player Player, direction Direction) error {
	g.playerState[player].looksIn = direction
	return nil
}

func (g *Game) PlayerDirection(player Player) Direction {
	return g.playerState[player].looksIn
}

func (g *Game) PlayerMove(player Player, direction Direction) error {
	if _, ok := g.playerMoveTransition[player]; ok {
		return errors.New("Player in action")
	}

	playerPos := g.playerState[player].mapPos
	targetPos := playerPos.Neighbor(direction)

	targetCell := g.gameMap.Cell(targetPos)
	if targetCell.IsWall() {
		return errors.New("Is Wall")
	}

	if !g.IsEmpty(targetPos) {
		if g.IsDoor(targetPos) && g.Door(targetPos).IsOpen() {
			// door is open
		} else if g.IsDoor(targetPos) && g.PlayerCanPassDoor(player, g.doors[targetPos]) {
			// block door close
		} else if g.IsBannWall(targetPos) && g.PlayerCanPassBannWall(player, g.bannWalls[targetPos]) {
			log.Println("Bann wall passable")
		} else if g.IsBoulder(targetPos) && g.PlayerCanPassBoulder(player, g.boulders[targetPos]) {
			log.Println("Boulder passable")
		} else if g.PosEmptyInFuture(targetPos) {
			log.Println("Will be empty - something moving away nothing in")
		} else {
			return errors.New("Is not empty and will not be empty")
		}

		if g.IsPlayer(targetPos) {
			return errors.New("Player on field")
		} else {

		}

	}

	g.playerMoveTransition[player] =
		NewPlayerMoveTransition(player, playerPos, targetPos)

	return nil
}

func (g *Game) ActivateTrigger(player Player, trigger *Trigger) {
	log.Println("Trigger activated", trigger)
	trigger.isActive = !trigger.isActive

	if trigger.linkedDoor != nil {
		log.Println("Door toogled")
		g.ToggleDoor(trigger.linkedDoor)
	}

	if trigger.linkedBannWall != nil {
		trigger.linkedBannWall.isActive = !trigger.linkedBannWall.isActive
	}

	if trigger.spawnBoulder != nil {
		trigger.spawnBoulder.active = true
		for _, player := range g.players {
			if Human == player {
				g.playerCans[player].canPassBoulder[trigger.spawnBoulder] = false
				g.playerCans[player].canPush[trigger.spawnBoulder] = true
				g.playerVis[player].visBoulder[trigger.spawnBoulder] = true
			} else {
				g.playerCans[player].canPassBoulder[trigger.spawnBoulder] = true
				g.playerCans[player].canPush[trigger.spawnBoulder] = false
				g.playerVis[player].visBoulder[trigger.spawnBoulder] = true
			}
		}
	}
}

func (g *Game) PlayerDoesAction(player Player) bool {
	_, ok := g.playerActionTransition[player]
	return ok
}

func (g *Game) PlayerAction(player Player, direction Direction) error {
	if g.PlayerIsWalking(player) {
		return errors.New("Player in action")
	}

	if g.PlayerDoesAction(player) {
		return errors.New("Player in action")
	}

	playerPos := g.playerState[player].mapPos
	playerCell := g.gameMap.Cell(playerPos)

	// fast fix
	if trigger, ok := playerCell.accessibleTriggers[DirWest]; ok {
		if _, ok := g.triggerTransition[trigger]; ok {
			return errors.New("Not possible")
			// trigger in transition
		}

		if _, ok := g.doorTransition[trigger.linkedDoor]; ok {
			return errors.New("Not possible")
			// door in transition
		}

		if g.PlayerCanTrigger(player, trigger) {
			g.ActivateTrigger(player, trigger)
			g.playerActionTransition[player] = NewPlayerActionTransition(player)
			return nil
		} else {
			return errors.New("Not authorized")
		}

		// feedback - cannot do
	}

	// actionOnCell := g.gameMap.NeighborOfCell(playerCell, direction)
	boulderPos := playerPos.Neighbor(direction)
	if g.IsBoulder(boulderPos) {
		boulder := g.boulders[boulderPos]
		if _, ok := g.boulderTransition[boulder]; ok {
			return errors.New("Boulder already moving")
		}
		targetPos := playerPos.Neighbor(direction).Neighbor(direction)

		if _, ok := g.boulderTransition[boulder]; ok {
			// boulder in transition
		}

		// is empty
		if !g.IsEmpty(targetPos) {
			if g.IsDoor(targetPos) && g.Door(targetPos).IsOpen() {
				// door is open
			} else if g.PosEmptyInFuture(targetPos) {
				// but something moves away from it
			} else {
				return errors.New("Is not empty and will not be empty")
			}
		}

		if g.PlayerCanPush(player, boulder) {
		} else {
			return errors.New("Not authorized")
		}

		g.boulderTransition[boulder] =
			NewBoulderTransition(boulder, boulderPos, targetPos)
		g.playerActionTransition[player] = NewPlayerActionTransition(player)
		return nil

		// feedback - cannot do
	}
	return errors.New("No action")
}

type PlayerCans struct {
	canPassDoor     map[*Door]bool
	canPassBoulder  map[*Boulder]bool
	canPassBannWall map[*BannWall]bool
	canPush         map[*Boulder]bool
	canTrigger      map[*Trigger]bool
	canPressure     map[*Plate]bool
}

type PlayerVis struct {
	visDoor     map[*Door]bool
	visTrigger  map[*Trigger]bool
	visCell     map[MapPosition]bool
	visBannWall map[*BannWall]bool
	visBoulder  map[*Boulder]bool
	visPlate    map[*Plate]bool
	visPlayer   bool
}

func (g *Game) PlayerCanSeeOtherPlayer(player Player) bool {
	return g.playerVis[player].visPlayer
}

func (g *Game) PlayerCanSeeCell(player Player, pos MapPosition) bool {
	return g.playerVis[player].visCell[pos]
}

func (g *Game) PlayerCanSeeTrigger(player Player, t *Trigger) bool {
	trigPos := NewMapPosition(-1, -1)
	for pos, trig := range g.triggers {
		if trig == t {
			trigPos = pos
		}
	}
	return g.playerVis[player].visTrigger[t] && g.PlayerCanSeeCell(player, trigPos)
}
func (g *Game) PlayerCanSeeDoor(player Player, d *Door) bool {
	trigPos := NewMapPosition(-1, -1)
	for pos, trig := range g.doors {
		if trig == d {
			trigPos = pos
		}
	}
	return g.playerVis[player].visDoor[d] && g.PlayerCanSeeCell(player, trigPos)
}

func (g *Game) PlayerCanSeeBoulder(player Player, b *Boulder) bool {
	trigPos := NewMapPosition(-1, -1)
	for pos, bw2 := range g.boulders {
		if bw2 == b {
			trigPos = pos
		}
	}
	return g.playerVis[player].visBoulder[b] && g.PlayerCanSeeCell(player, trigPos)
}

func (g *Game) PlayerCanSeePlate(player Player, p *Plate) bool {
	trigPos := NewMapPosition(-1, -1)
	for pos, bw2 := range g.plates {
		if bw2 == p {
			trigPos = pos
		}
	}
	return g.playerVis[player].visPlate[p] && g.PlayerCanSeeCell(player, trigPos)
}

func (g *Game) PlayerCanSeeBannWall(player Player, bw *BannWall) bool {
	trigPos := NewMapPosition(-1, -1)
	for pos, bw2 := range g.bannWalls {
		if bw2 == bw {
			trigPos = pos
		}
	}
	return g.playerVis[player].visBannWall[bw] && g.PlayerCanSeeCell(player, trigPos)
}

func NewPlayerVis() *PlayerVis {
	return &PlayerVis{
		visDoor:     make(map[*Door]bool),
		visTrigger:  make(map[*Trigger]bool),
		visCell:     make(map[MapPosition]bool),
		visBannWall: make(map[*BannWall]bool),
		visBoulder:  make(map[*Boulder]bool),
		visPlate:    make(map[*Plate]bool),
		visPlayer:   false,
	}
}

func NewPlayerCans() *PlayerCans {
	return &PlayerCans{
		canPassDoor:     make(map[*Door]bool),
		canPassBoulder:  make(map[*Boulder]bool),
		canPassBannWall: make(map[*BannWall]bool),
		canPush:         make(map[*Boulder]bool),
		canTrigger:      make(map[*Trigger]bool),
		canPressure:     make(map[*Plate]bool),
	}
}

type TriggerTransition struct {
	dtime   time.Duration
	toState bool
	// afterActive
}

type BoulderTransition struct {
	dtime   time.Duration
	boulder *Boulder
	fromPos MapPosition
	toPos   MapPosition
}

type BannWallTransition struct {
	dtime   time.Duration
	toState bool
}

type PlayerMoveTransition struct {
	dtime   time.Duration
	player  Player
	fromPos MapPosition
	toPos   MapPosition
}

type PlayerActionTransition struct {
	dtime  time.Duration
	player Player
}

func (pat *PlayerActionTransition) IsFinished() bool {
	return pat.dtime > GlobalConfig.actionTime
}

func (pmt *PlayerMoveTransition) IsFinished() bool {
	return pmt.dtime > GlobalConfig.walkTime
}

func (dt *DoorTransition) IsFinished() bool {
	return dt.dtime > 500*time.Millisecond
}

func (bt *BoulderTransition) IsFinished() bool {
	return bt.dtime > GlobalConfig.rollTime
}

func (bwt *BannWallTransition) IsFinished() bool {
	return bwt.dtime > 500*time.Millisecond
}

func (pat *PlayerActionTransition) Update(dt time.Duration) {
	pat.dtime += dt
}

func (doort *DoorTransition) Update(dt time.Duration) {
	doort.dtime += dt
}

func (pmt *BoulderTransition) Update(dt time.Duration) {
	pmt.dtime += dt
}

func (pmt *PlayerMoveTransition) Update(dt time.Duration) {
	pmt.dtime += dt
}

func (pmt *PlayerMoveTransition) UpdateGameState(g *Game) {
	log.Println("Update game state", pmt.player)
	log.Println("Old pos was", g.playerState[pmt.player].mapPos)
	log.Println("transition origin", pmt.OriginPos())
	log.Println("New pos is", pmt.TargetPos())

	g.playerState[pmt.player].mapPos = pmt.TargetPos()
}

func (pmt *PlayerActionTransition) UpdateGameState(g *Game) {
	// g.playerState[pmt.player].
}

func (dt *DoorTransition) UpdateGameState(g *Game) {
	dt.door.isOpen = dt.toState
}

func (bt *BoulderTransition) UpdateGameState(g *Game) {
	currentPos, _ := g.BoulderPos(bt.boulder)
	delete(g.boulders, currentPos)
	g.boulders[bt.TargetPos()] = bt.boulder
}

func NewBoulderTransition(b *Boulder, from, to MapPosition) *BoulderTransition {
	return &BoulderTransition{
		dtime:   0,
		boulder: b,
		fromPos: from,
		toPos:   to,
	}
}

func NewPlayerMoveTransition(player Player, from, to MapPosition) *PlayerMoveTransition {
	return &PlayerMoveTransition{
		player:  player,
		dtime:   0,
		fromPos: from,
		toPos:   to,
	}
}

func NewDoorTransition(door *Door) *DoorTransition {
	return &DoorTransition{
		dtime:   0,
		toState: !door.isOpen,
	}
}

func NewBannWallTransition(bw *BannWall) *BannWallTransition {
	return &BannWallTransition{
		dtime:   0,
		toState: !bw.isActive,
	}
}

func NewPlayerActionTransition(player Player) *PlayerActionTransition {
	return &PlayerActionTransition{
		dtime:  0,
		player: player,
	}
}

func (pmt *BoulderTransition) OriginPos() MapPosition {
	return pmt.fromPos
}

func (pmt *BoulderTransition) TargetPos() MapPosition {
	return pmt.toPos
}

func (pmt *BoulderTransition) InterpPos() Position {
	maxTime := GlobalConfig.walkTime
	x, y := pmt.TargetPos().X(), pmt.TargetPos().Y()
	ox, oy := pmt.OriginPos().X(), pmt.OriginPos().Y()
	return Position{
		x: float32(ox) + float32(x-ox)*float32(pmt.dtime)/float32(maxTime),
		y: float32(oy) + float32(y-oy)*float32(pmt.dtime)/float32(maxTime),
	}
}

func (pmt *BoulderTransition) Frame() int {
	return 0
}

func (pat *PlayerActionTransition) Frame() int {
	return SplitTimeEven(4, GlobalConfig.actionTime, pat.dtime)
}

func (pmt *PlayerMoveTransition) OriginPos() MapPosition {
	return pmt.fromPos
}

func (pmt *PlayerMoveTransition) TargetPos() MapPosition {
	return pmt.toPos
}

func (pmt *PlayerMoveTransition) InterpPos() Position {
	maxTime := GlobalConfig.walkTime
	x, y := pmt.TargetPos().X(), pmt.TargetPos().Y()
	ox, oy := pmt.OriginPos().X(), pmt.OriginPos().Y()
	return Position{
		x: float32(ox) + float32(x-ox)*float32(pmt.dtime)/float32(maxTime),
		y: float32(oy) + float32(y-oy)*float32(pmt.dtime)/float32(maxTime),
	}
}

func SplitTimeEven(num int, total, local time.Duration) int {
	i := 0
	for dt := float32(0); dt < float32(total); dt += float32(total) / float32(num) {
		if float32(local) < dt {
			return i
		}
		i += 1
	}
	return 0
}

func (pmt *PlayerMoveTransition) Frame() int {
	return SplitTimeEven(4, GlobalConfig.walkTime, pmt.dtime)
}

func (g *Game) SetRoomVisible(room *Room) {
	room.isVisible = true
	for _, pos := range room.cells {
		for _, player := range g.players {
			g.playerVis[player].visCell[pos] = true
		}
	}
}

func (g *Game) PlateAtPos(pos MapPosition) bool {
	_, ok := g.plates[pos]
	return ok
}

func (g *Game) PlateCanBeActivated(pos MapPosition) bool {
	if !g.PlateAtPos(pos) {
		return false
	}

	plate := g.plates[pos]

	for _, player := range g.players {
		// player on field and player can pressure
		if g.playerState[player].mapPos == pos &&
			g.playerCans[player].canPressure[plate] {
			return true
		}
	}

	// boulder on field
	if _, ok := g.boulders[pos]; ok {
		return true
	}

	return false
}

func (g *Game) ActivatePlate(pos MapPosition) {
	plate := g.plates[pos]
	if plate.linkedBannWall != nil {
		plate.linkedBannWall.isActive = false
	}

	if plate.linkedDoor != nil {
		g.ToggleDoor(plate.linkedDoor)
	}
}

func (g *Game) Update(t time.Duration) {
	for player, moveTransition := range g.playerMoveTransition {
		moveTransition.Update(t)

		if moveTransition.IsFinished() {
			moveTransition.UpdateGameState(g)
			// check if plate underneath new position
			delete(g.playerMoveTransition, player)
			door, ok := g.doors[moveTransition.TargetPos()]
			if ok && player == Ghost {
				room := door.linkedRoom
				g.MakeRoomVisible(room)
			}
		}
	}

	for door, doorTransition := range g.doorTransition {
		doorTransition.Update(t)

		if doorTransition.IsFinished() {
			doorTransition.UpdateGameState(g)
			g.SetRoomVisible(door.linkedRoom)
			delete(g.doorTransition, door)
		}

	}

	for boulder, boulderTransition := range g.boulderTransition {
		boulderTransition.Update(t)

		if boulderTransition.IsFinished() {
			boulderTransition.UpdateGameState(g)
			// check if plate underneath new position
			boulderPos, _ := g.BoulderPos(boulder)
			if g.PlateCanBeActivated(boulderPos) {
				g.ActivatePlate(boulderPos)
			}
			delete(g.boulderTransition, boulder)
		}
	}

	for player, actionTransition := range g.playerActionTransition {
		actionTransition.Update(t)

		if actionTransition.IsFinished() {
			actionTransition.UpdateGameState(g)
			delete(g.playerActionTransition, player)
		}
	}

	for _, visDelay := range g.visDelay {
		visDelay.Update(t)
		if visDelay.IsFinished() {
			visDelay.UpdateGameState(g)
			delete(g.visDelay, visDelay.player)
		}
	}

	for _, visTrans := range g.playerVisTransition {
		visTrans.Update(t)
		if visTrans.IsFinished() {
			visTrans.UpdateGameState(g)
			delete(g.playerVisTransition, visTrans.player)
			g.visDelay[visTrans.player] = NewVisDelay(visTrans.player)
		}
	}
}

func (g *Game) PosEmptyInFuture(pos MapPosition) bool {
	canBeEmpty := false
	for _, playerMoveTrans := range g.playerMoveTransition {
		if playerMoveTrans.OriginPos() == pos {
			canBeEmpty = true
		}
	}

	for _, boulderTransition := range g.boulderTransition {
		if boulderTransition.OriginPos() == pos {
			canBeEmpty = true
		}
	}

	if canBeEmpty && !g.PosOccupiedInFuture(pos) {
		return true
	} else {
		return false
	}
}

func (g *Game) PosOccupiedInFuture(pos MapPosition) bool {
	for _, playerMoveTrans := range g.playerMoveTransition {
		if playerMoveTrans.TargetPos() == pos {
			return true
		}
	}

	for _, boulderTransition := range g.boulderTransition {
		if boulderTransition.TargetPos() == pos {
			return true
		}
	}

	return false
}

func (g *Game) Width() int {
	return len(g.gameMap.cells[0])
}

func (g *Game) Height() int {
	return len(g.gameMap.cells)
}

func (g *Game) Cell(pos MapPosition) *Cell {
	return g.gameMap.Cell(pos)
}

func (g *Game) Players() []Player {
	return g.players
}

func (g *Game) PlayerRenderPos(player Player) Position {
	if transition, ok := g.playerMoveTransition[player]; ok {
		return transition.InterpPos()
	}

	return Position{
		x: float32(g.playerState[player].mapPos.X()),
		y: float32(g.playerState[player].mapPos.Y()),
	}
}

func (g *Game) BoulderRenderPos(boulder *Boulder) Position {
	if transition, ok := g.boulderTransition[boulder]; ok {
		return transition.InterpPos()
	}
	pos, _ := g.BoulderPos(boulder)
	return Position{
		x: float32(pos.X()),
		y: float32(pos.Y()),
	}
}

func (g *Game) BannWalls() map[MapPosition]*BannWall {
	return g.bannWalls
}

func (g *Game) Boulders() map[MapPosition]*Boulder {
	return g.boulders
}

func (g *Game) Doors() map[MapPosition]*Door {
	return g.doors
}

func (g *Game) Triggers() map[MapPosition]*Trigger {
	return g.triggers
}

func (g *Game) Plates() map[MapPosition]*Plate {
	return g.plates
}

type Game struct {
	players []Player
	gameMap *Map

	playerState map[Player]*PlayerState

	doors     map[MapPosition]*Door
	boulders  map[MapPosition]*Boulder
	bannWalls map[MapPosition]*BannWall

	rooms    []*Room
	triggers map[MapPosition]*Trigger
	plates   map[MapPosition]*Plate

	playerCans map[Player]*PlayerCans
	playerVis  map[Player]*PlayerVis

	playerMoveTransition   map[Player]MoveableTransition
	visDelay               map[Player]*VisDelay
	playerVisTransition    map[Player]*VisStateTransition
	playerActionTransition map[Player]AnimTransition
	boulderTransition      map[*Boulder]MoveableTransition

	triggerTransition map[*Trigger]*TriggerTransition
	doorTransition    map[*Door]*DoorTransition

	// spriteCarBG   *Sprite
	// spriteWaiting *Sprite

	running bool

	music *mixer.Music
	// font  *ttf.Font
}

func (g *Game) IsEmpty(pos MapPosition) bool {
	if !(g.IsDoor(pos) || g.IsBoulder(pos) || g.IsWall(pos) || g.IsPlayer(pos)) {
		log.Println("Is states", g.IsDoor(pos), g.IsBoulder(pos), g.IsWall(pos), g.IsPlayer(pos))
		return true
	}
	log.Println("Is states", g.IsDoor(pos), g.IsBoulder(pos), g.IsWall(pos), g.IsPlayer(pos))
	return false
}

func (g *Game) IsPlayer(pos MapPosition) bool {
	for _, state := range g.playerState {
		if state.mapPos == pos {
			return true
		}
	}

	return false
}

func (g *Game) IsWall(pos MapPosition) bool {
	return g.gameMap.Cell(pos).IsWall()
}

func (g *Game) IsDoor(pos MapPosition) bool {
	if _, ok := g.doors[pos]; ok {
		return true
	}
	return false
}

func (g *Game) Door(pos MapPosition) *Door {
	if door, ok := g.doors[pos]; ok {
		return door
	}

	return nil
}

// func (g *Game) IsPassableDoor(pos MapPosition) bool {
// 	if door, ok := g.doors[pos]; ok {
// 		if door.isOpen {
// 			return true
// 		}

// 		return false
// 	}
// }

func (g *Game) IsBoulder(pos MapPosition) bool {
	if _, ok := g.boulders[pos]; ok {
		return true
	}
	return false
}

func (g *Game) IsBannWall(pos MapPosition) bool {
	if _, ok := g.bannWalls[pos]; ok {
		return true
	}
	return false
}

// func (g *Game) Trigger(*Trigger)

func (g *Game) Render(screen *sdl.Surface) {
	// render cells
	// render walls
	// render boulders, doors and triggers
	// render player
}

func NewMap(width, height int) *Map {
	cells := make([][]*Cell, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]*Cell, width)
		for x := 0; x < width; x++ {
			cells[y][x] = NewCell()
		}
	}

	return &Map{
		cells: cells,
	}
}

func (g *Game) SetWall(pos MapPosition) {
	g.gameMap.Cell(pos).SetWall()
}

func (g *Game) SetBoulder(pos MapPosition, active bool) *Boulder {
	b := NewBoulder(active)
	g.boulders[pos] = b

	for _, player := range g.players {
		g.playerCans[player].canPassBoulder[b] = false
		g.playerCans[player].canPush[b] = false
	}

	return b
}

func (g *Game) SetPlate(pos MapPosition) *Plate {
	p := NewPlate()
	g.plates[pos] = p

	for _, player := range g.players {
		g.playerCans[player].canPressure[p] = false
	}

	return p
}

func (g *Game) SetDoor(pos MapPosition) *Door {
	d := NewDoor()
	g.doors[pos] = d

	for _, player := range g.players {
		g.playerCans[player].canPassDoor[d] = false
		g.playerVis[player].visDoor[d] = false
	}

	return d
}

func (g *Game) SetBannWall(pos MapPosition, bannWallType int) *BannWall {
	bw := NewBannWall(bannWallType)
	g.bannWalls[pos] = bw

	for _, player := range g.players {
		g.playerCans[player].canPassBannWall[bw] = false
	}

	return bw
}

func (g *Game) IsHuman(player Player) bool {
	if int(player) == 0 {
		return true
	} else {
		return false
	}
}

func (g *Game) IsGhost(player Player) bool {
	return !g.IsHuman(player)
}

func (g *Game) NewPlayer(id int) Player {
	player := Player(id)
	g.playerCans[player] = NewPlayerCans()
	g.playerVis[player] = NewPlayerVis()

	g.players = append(g.players, player)

	log.Println(id)
	log.Println(GlobalConfig)
	startPos := GlobalConfig.playerStartPos[id]

	g.playerState[player] = NewPlayerState(startPos,
		GlobalConfig.playerStartLook[id])

	// make all cells invisible
	for y := 0; y < g.Height(); y++ {
		for x := 0; x < g.Width(); x++ {
			g.playerVis[player].visCell[NewMapPosition(x, y)] = false
		}
	}

	for _, room := range g.rooms {
		if room.isVisible {
			for _, cellPos := range room.cells {
				g.playerVis[player].visCell[cellPos] = true
			}
		}
	}

	for _, door := range g.doors {
		g.playerVis[player].visDoor[door] = true
		g.playerCans[player].canPassDoor[door] = false
	}

	for _, bw := range g.bannWalls {
		if g.IsHuman(player) {
			g.playerVis[player].visBannWall[bw] = false
		} else if g.IsGhost(player) {
			g.playerVis[player].visBannWall[bw] = false
		}
	}

	for _, boulder := range g.boulders {
		if g.IsHuman(player) {
			if !boulder.IsActive() {
				g.playerCans[player].canPassBoulder[boulder] = true
				g.playerCans[player].canPush[boulder] = false
				g.playerVis[player].visBoulder[boulder] = false
			} else {
				g.playerCans[player].canPassBoulder[boulder] = false
				g.playerCans[player].canPush[boulder] = true
				g.playerVis[player].visBoulder[boulder] = true
			}
		} else if g.IsGhost(player) {
			g.playerCans[player].canPassBoulder[boulder] = true
			g.playerCans[player].canPush[boulder] = false

			if !boulder.IsActive() {
				g.playerVis[player].visBoulder[boulder] = false
			} else {
				g.playerVis[player].visBoulder[boulder] = true
			}
		}
	}

	for _, trigger := range g.triggers {
		if trigger.canTrigger == GhostOrHuman || player == trigger.canTrigger {
			g.playerCans[player].canTrigger[trigger] = true
		} else {
			g.playerCans[player].canTrigger[trigger] = false
		}

		if trigger.canVis == GhostOrHuman || player == trigger.canVis {
			g.playerVis[player].visTrigger[trigger] = true
		} else {
			g.playerVis[player].visTrigger[trigger] = false
		}
	}

	for _, plate := range g.plates {
		if g.IsHuman(player) {
			g.playerVis[player].visPlate[plate] = false
		} else {
			g.playerVis[player].visPlate[plate] = true
		}
	}

	return player
}

func (g *Game) SetPlayerCanSeeDoor(player Player, door *Door) {
	g.playerVis[player].visDoor[door] = true
}

func (g *Game) SetPlayerCanSeeTrigger(player Player, t *Trigger) {
	g.playerVis[player].visTrigger[t] = true
}

func (g *Game) SetPlayerCanPassDoor(player Player, door *Door) {
	g.playerCans[player].canPassDoor[door] = true
}

func (g *Game) SetPlayerCanPassBoulder(player Player, b *Boulder) {
	g.playerCans[player].canPassBoulder[b] = true
}

func (g *Game) SetPlayerCanPassBannWall(player Player, bw *BannWall) {
	g.playerCans[player].canPassBannWall[bw] = true
}

func (g *Game) SetPlayerCanPushBoulder(player Player, b *Boulder) {
	g.playerCans[player].canPush[b] = true
}

func (g *Game) SetPlayerCanTrigger(player Player, t *Trigger) {
	g.playerCans[player].canTrigger[t] = true
}
func (g *Game) PlayerCanPassDoor(player Player, d *Door) bool {
	return g.playerCans[player].canPassDoor[d]
}

func (g *Game) PlayerCanPassBoulder(player Player, t *Boulder) bool {
	return g.playerCans[player].canPassBoulder[t]
}

func (g *Game) PlayerCanPassBannWall(player Player, t *BannWall) bool {
	return g.playerCans[player].canPassBannWall[t]
}

func (g *Game) PlayerCanPush(player Player, t *Boulder) bool {
	return g.playerCans[player].canPush[t]
}

func (g *Game) PlayerCanTrigger(player Player, t *Trigger) bool {
	return g.playerCans[player].canTrigger[t]
}

func (g *Game) MakePlayerToGhost(player Player) {
	for _, door := range g.doors {
		g.SetPlayerCanPassDoor(player, door)
		g.SetPlayerCanSeeDoor(player, door)
	}

	for _, bannWall := range g.bannWalls {
		g.playerCans[player].canPassBannWall[bannWall] = false
	}

	for _, boulder := range g.boulders {
		g.SetPlayerCanPassBoulder(player, boulder)
		g.playerCans[player].canPush[boulder] = false
	}

	for _, trigger := range g.triggers {
		g.SetPlayerCanSeeTrigger(player, trigger)
	}
}

func (g *Game) BoulderPos(wantedBoulder *Boulder) (MapPosition, error) {
	for pos, b := range g.boulders {
		if wantedBoulder == b {
			return pos, nil
		}
	}
	return MapPosition{-1, -1}, errors.New("No pos")
}

func NewGame() (*Game, error) {
	width, height := GlobalConfig.mapWidth, GlobalConfig.mapHeight
	r := &Game{
		players:     make([]Player, 0, 2),
		gameMap:     NewMap(width, height),
		playerState: make(map[Player]*PlayerState, 2),

		doors:     make(map[MapPosition]*Door),
		boulders:  make(map[MapPosition]*Boulder),
		bannWalls: make(map[MapPosition]*BannWall),
		plates:    make(map[MapPosition]*Plate),
		triggers:  make(map[MapPosition]*Trigger),

		playerCans: make(map[Player]*PlayerCans),
		playerVis:  make(map[Player]*PlayerVis),

		visDelay:               make(map[Player]*VisDelay),
		playerVisTransition:    make(map[Player]*VisStateTransition),
		playerMoveTransition:   make(map[Player]MoveableTransition),
		playerActionTransition: make(map[Player]AnimTransition),
		boulderTransition:      make(map[*Boulder]MoveableTransition),

		triggerTransition: make(map[*Trigger]*TriggerTransition),
		doorTransition:    make(map[*Door]*DoorTransition),

		running: false,
		music:   nil,
	}

	BuildGame(r)

	// if r.music = mixer.LoadMUS("data/music.ogg"); r.music == nil {
	// 	return nil, errors.New(sdl.GetError())
	// }

	// if r.font = ttf.OpenFont("data/font.otf", 32); r.font == nil {
	// return nil, errors.New(sdl.GetError())
	// }

	// textWaiting := ttf.RenderUTF8_Blended(r.font, "Please start")
	// r.spriteWaiting = NewSpriteFromSurface(textWaiting)

	return r, nil
}

func (r *Game) Join(player *Player) {
	// if len(r.entities) == 0 {
	// mixer.ResumeMusic()
	// r.music.PlayMusic(-1)
	// }
	// entity := NewEntity(player)
	// entity.position.x = x
	// entity.position.y = y
	// r.entities = append(r.entities, entity)
}

func (r *Game) Leave(player *Player) {
	// for i := range r.entities {
	// 	if r.entities[i].owner == player {
	// 		r.entities[i] = r.entities[len(r.entities)-1]
	// 		r.entities = r.entities[:len(r.entities)-1]
	// 		break
	// 	}
	// }
	// if len(r.entities) == 0 {
	// 	r.running = false
	// 	mixer.PauseMusic()
	// }
}

func (r *Game) KeyPressed(input sdl.Keysym) {
	if input.Sym == sdl.K_SPACE {
		r.running = true
	}
}

type Entity struct {
	owner Player
}

func (e *Entity) Update(time time.Duration) {
	// t := float32(time)
}

func (e *Entity) Draw() {
}

func NewEntity(owner Player) *Entity {
	return &Entity{
		owner: owner,
	}
}

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

type CellType int

const (
	CellTypeEmpty CellType = iota
	CellTypeWall
	CellTypeDoor
)

var GlobalConfig *MapConfig = NewMapConfig()

type MapConfig struct {
	playerStartPos  []MapPosition
	playerStartLook []Direction
}

func NewMapConfig() *MapConfig {
	return &MapConfig{
		playerStartPos: []MapPosition{
			MapPosition{5, 5},
			MapPosition{4, 4},
		},
		playerStartLook: []Direction{
			DirEast,
			DirWest,
		},
	}
}

type Player int

type Door struct {
	isOpen bool
}

func NewDoor() *Door {
	return &Door{
		isOpen: false,
	}
}

type Trigger struct {
	linkedDoor *Door
}

func NewTrigger() *Trigger {
	return &Trigger{
		linkedDoor: nil,
	}
}

type Boulder struct {
}

func NewBoulder() *Boulder {
	return &Boulder{}
}

type Plate struct {
	isActive       bool
	linkedBannWall *BannWall
}

func NewPlate() *Plate {
	return &Plate{
		isActive:       false,
		linkedBannWall: nil,
	}
}

type BannWall struct {
	isActive bool
}

func NewBannWall() *BannWall {
	return &BannWall{
		isActive: true,
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
		return MapPosition{mp.x, mp.y + 1}
	case DirWest:
		return MapPosition{mp.x + 1, mp.y}
	case DirSouth:
		return MapPosition{mp.x, mp.y - 1}
	case DirEast:
		return MapPosition{mp.x - 1, mp.y}
	}
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

func (g *Game) SetTrigger(pos MapPosition, direction Direction) {
	t := NewTrigger()
	g.triggers = append(g.triggers, t)
	cell := g.gameMap.Cell(pos)
	cell.accessibleTriggers[direction] = t
}

type Direction int

const (
	DirNorth Direction = iota
	DirWest
	DirSouth
	DirEast
)

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

type MoveableTransition interface {
	Transition
	OriginPos() MapPosition
	TargetPos() MapPosition
	InterpPos() Position
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
		if g.PosEmptyInFuture(targetPos) {
		} else if g.IsDoor(targetPos) && g.PlayerCanPassDoor(player, g.doors[targetPos]) {
			// block door close
		} else if g.IsBannWall(targetPos) && g.PlayerCanPassBannWall(player, g.bannWalls[targetPos]) {
		} else if g.IsBoulder(targetPos) && g.PlayerCanPassBoulder(player, g.boulders[targetPos]) {
		} else {
			return errors.New("Is not empty and will not be empty")
		}
	}

	g.playerMoveTransition[player] =
		NewPlayerMoveTransition(playerPos, targetPos)

	return nil
}

func (g *Game) PlayerAction(player Player, direction Direction) error {
	if _, ok := g.playerMoveTransition[player]; ok {
		return errors.New("Player in action")
	}

	playerPos := g.playerState[player].mapPos
	playerCell := g.gameMap.Cell(playerPos)

	if trigger, ok := playerCell.accessibleTriggers[direction]; ok {
		if _, ok := g.triggerTransition[trigger]; ok {
			// trigger in transition
		}

		if _, ok := g.doorTransition[trigger.linkedDoor]; ok {
			// door in transition
		}

		if g.PlayerCanTrigger(player, trigger) {
			// trigger
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
		targetPos := playerPos.Neighbor(direction).Neighbor(direction)

		if _, ok := g.boulderTransition[boulder]; ok {
			// boulder in transition
		}

		// is empty
		if !g.IsEmpty(targetPos) {
			// but something moves away from it
			if g.PosEmptyInFuture(targetPos) {
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
}

type PlayerVis struct {
	visDoor    map[*Door]bool
	visTrigger map[*Trigger]bool
}

func NewPlayerVis() *PlayerVis {
	return &PlayerVis{
		visDoor:    make(map[*Door]bool),
		visTrigger: make(map[*Trigger]bool),
	}
}

func NewPlayerCans() *PlayerCans {
	return &PlayerCans{
		canPassDoor:     make(map[*Door]bool),
		canPassBannWall: make(map[*BannWall]bool),
		canPush:         make(map[*Boulder]bool),
		canTrigger:      make(map[*Trigger]bool),
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

func (pmt *PlayerMoveTransition) IsFinished() bool {
	return pmt.dtime > 10 // Magic ten
}

func (dt *DoorTransition) IsFinished() bool {
	return dt.dtime > 10
}

func (bt *BoulderTransition) IsFinished() bool {
	return bt.dtime > 10
}

func (bwt *BannWallTransition) IsFinished() bool {
	return bwt.dtime > 10
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
	g.playerState[pmt.player].mapPos = pmt.TargetPos()
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

func NewPlayerMoveTransition(from, to MapPosition) *PlayerMoveTransition {
	return &PlayerMoveTransition{
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

func (pmt *BoulderTransition) OriginPos() MapPosition {
	return pmt.fromPos
}

func (pmt *BoulderTransition) TargetPos() MapPosition {
	return pmt.toPos
}

func (pmt *BoulderTransition) InterpPos() Position {
	maxTime := 500 * time.Millisecond
	x, y := pmt.TargetPos().X(), pmt.TargetPos().Y()
	ox, oy := pmt.OriginPos().X(), pmt.OriginPos().Y()
	return Position{
		x: float32(x-ox) * float32(pmt.dtime/maxTime),
		y: float32(y-oy) * float32(pmt.dtime/maxTime),
	}
}

func (pmt *PlayerMoveTransition) OriginPos() MapPosition {
	return pmt.fromPos
}

func (pmt *PlayerMoveTransition) TargetPos() MapPosition {
	return pmt.toPos
}

func (pmt *PlayerMoveTransition) InterpPos() Position {
	maxTime := 500 * time.Millisecond
	x, y := pmt.TargetPos().X(), pmt.TargetPos().Y()
	ox, oy := pmt.OriginPos().X(), pmt.OriginPos().Y()
	return Position{
		x: float32(x-ox) * float32(pmt.dtime/maxTime),
		y: float32(y-oy) * float32(pmt.dtime/maxTime),
	}
}

func (g *Game) Update(t time.Duration) {
	for player, moveTransition := range g.playerMoveTransition {
		moveTransition.Update(t)

		if moveTransition.IsFinished() {
			moveTransition.UpdateGameState(g)
		}
		delete(g.playerMoveTransition, player)
	}

	for door, doorTransition := range g.doorTransition {
		doorTransition.Update(t)

		if doorTransition.IsFinished() {
			doorTransition.UpdateGameState(g)
		}

		delete(g.doorTransition, door)
	}

	for boulder, boulderTransition := range g.boulderTransition {
		boulderTransition.Update(t)

		if boulderTransition.IsFinished() {
			boulderTransition.UpdateGameState(g)
		}

		delete(g.boulderTransition, boulder)
	}
}

func (g *Game) PosEmptyInFuture(pos MapPosition) bool {
	for _, playerMoveTrans := range g.playerMoveTransition {
		if playerMoveTrans.TargetPos() == pos {
			return false
		}
	}

	for _, boulderTransition := range g.boulderTransition {
		if boulderTransition.TargetPos() == pos {
			return false
		}
	}

	return true
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

type Game struct {
	players []Player
	gameMap *Map

	playerState map[Player]*PlayerState

	doors     map[MapPosition]*Door
	boulders  map[MapPosition]*Boulder
	bannWalls map[MapPosition]*BannWall

	triggers []*Trigger

	playerCans map[Player]*PlayerCans
	playerVis  map[Player]*PlayerVis

	playerMoveTransition map[Player]MoveableTransition
	boulderTransition    map[*Boulder]MoveableTransition

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
		return true
	}
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

func (g *Game) SetBoulder(pos MapPosition) {
	b := NewBoulder()
	g.boulders[pos] = b

	for _, player := range g.players {
		g.playerCans[player].canPassBoulder[b] = false
		g.playerCans[player].canPush[b] = false
	}
}

func (g *Game) SetDoor(pos MapPosition) {
	d := NewDoor()
	g.doors[pos] = d

	for _, player := range g.players {
		g.playerCans[player].canPassDoor[d] = false
		g.playerVis[player].visDoor[d] = false
	}
}

func (g *Game) SetBannWall(pos MapPosition, startLookIn Direction) {
	bw := NewBannWall()
	g.bannWalls[pos] = bw

	for _, player := range g.players {
		g.playerCans[player].canPassBannWall[bw] = false
	}
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
	width, height := 10, 10
	r := &Game{
		players:     make([]Player, 0, 2),
		gameMap:     NewMap(width, height),
		playerState: make(map[Player]*PlayerState, 2),

		doors:     make(map[MapPosition]*Door),
		boulders:  make(map[MapPosition]*Boulder),
		bannWalls: make(map[MapPosition]*BannWall),

		playerCans: make(map[Player]*PlayerCans),
		playerVis:  make(map[Player]*PlayerVis),

		playerMoveTransition: make(map[Player]MoveableTransition),
		boulderTransition:    make(map[*Boulder]MoveableTransition),

		triggerTransition: make(map[*Trigger]*TriggerTransition),
		doorTransition:    make(map[*Door]*DoorTransition),

		running: false,
		music:   nil,
	}

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

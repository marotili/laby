package main

import (
	"fmt"
	"laby/game"
	"log"
	"math"
)

type GhostSprites struct {
	animWalk   map[game.Direction][]*Sprite
	animAction map[game.Direction][]*Sprite
}

type HumanSprites struct {
	animWalk   map[game.Direction][]*Sprite
	animAction map[game.Direction][]*Sprite
}

type WallSprites struct {
	walls []*Sprite
}

type FloorSprites struct {
	floor   *Sprite
	boulder *Sprite
	door    *Sprite
	ban     []*Sprite
}

type ToolSprites struct {
	triggersA *Sprite
	triggersB *Sprite
}

func LoadToolSprites() *ToolSprites {
	triggerSpriteA, err := NewSprite("data/trigger0.png", 128, 128)
	if err != nil {
		log.Fatal("Could not open trigger file")
	}

	triggerSpriteB, err := NewSprite("data/trigger1.png", 128, 128)
	if err != nil {
		log.Fatal("Could not open trigger file")
	}

	return &ToolSprites{
		triggersA: triggerSpriteA,
		triggersB: triggerSpriteB,
	}
}

func LoadFloorSprites() *FloorSprites {
	floorSprite, err := NewSprite("data/floor/floor.png", 64, 64)
	if err != nil {
		log.Fatal("Could not open floor tile")
	}

	boulderSprite, err := NewSprite("data/floor/boulder.png", 64, 64)
	if err != nil {
		log.Fatal("Could not open boulder tile")
	}

	doorSprite, err := NewSprite("data/floor/door.png", 64, 64)
	if err != nil {
		log.Fatal("Could not open boulder tile")
	}

	banSprites := make([]*Sprite, 4)
	for i := 1; i < 5; i++ {
		banSprites[i-1], err = NewSprite(fmt.Sprintf("data/ban/%d.png", i), 64, 64)
		if err != nil {
			log.Fatal("Could not load ban")
		}
	}

	return &FloorSprites{
		floor:   floorSprite,
		boulder: boulderSprite,
		door:    doorSprite,
		ban:     banSprites,
	}
}

func LoadGhostSprites() *GhostSprites {
	baseDirectory := "data/Geist/"

	dirDirs := []string{
		"hinten",
		"links",
		"vorne",
		"rechts",
	}

	prefix := []string{
		"",
		// "schieben",
	}

	const Walk = 0
	const Action = 1

	// actionDirectory := "data/ghost/action"
	numWalkFrames := 4
	// numActionFrames := 4
	// numStandFrames := 4
	width, height := float32(256), float32(256)

	walkSprites := make(map[game.Direction][]*Sprite, 4)
	// actionSprites := make(map[game.Direction][]*Sprite, 4)

	for _, dir := range game.Dirs() {
		walkSprites[dir] = make([]*Sprite, numWalkFrames)

		for i := 0; i < numWalkFrames; i++ {
			fileName := fmt.Sprintf("%s/%s/%s%d.png", baseDirectory, dirDirs[int(dir)], prefix[Walk], i+1)
			sprite, err := NewSprite(fileName, width, height)
			if err != nil {
				// log.Fatal("Could not open ghost walk tile", fileName)
			}
			walkSprites[dir][i] = sprite
		}
	}

	// for _, dir := range game.Dirs() {
	// 	actionSprites[dir] = make([]*Sprite, numActionFrames)

	// 	for i := 0; i < numActionFrames; i++ {
	// 		fileName := fmt.Sprintf("%s/%s/%s%d.png", baseDirectory, dirDirs[int(dir)], prefix[Action], i+1)
	// 		sprite, err := NewSprite(fileName, width, height)
	// 		if err != nil {
	// 			// log.Fatal("Could not open ghost walk tile", fileName)
	// 		}
	// 		actionSprites[dir][i] = sprite
	// 	}
	// }

	return &GhostSprites{
		animWalk:   walkSprites,
		animAction: nil,
	}
}

func LoadHumanSprites() *HumanSprites {
	baseDirectory := "data/Neira/"

	dirDirs := []string{
		"hinten",
		"links",
		"vorne",
		"rechts",
	}

	prefix := []string{
		"",
		"schieben",
	}

	const Walk = 0
	const Action = 1

	// actionDirectory := "data/ghost/action"
	numWalkFrames := 4
	numActionFrames := 4
	// numStandFrames := 4
	width, height := float32(256), float32(256)

	walkSprites := make(map[game.Direction][]*Sprite, 4)
	actionSprites := make(map[game.Direction][]*Sprite, 4)

	for _, dir := range game.Dirs() {
		walkSprites[dir] = make([]*Sprite, numWalkFrames)

		for i := 0; i < numWalkFrames; i++ {
			fileName := fmt.Sprintf("%s/%s/%s%d.png", baseDirectory, dirDirs[int(dir)], prefix[Walk], i+1)
			sprite, err := NewSprite(fileName, width, height)
			if err != nil {
				// log.Fatal("Could not open ghost walk tile", fileName)
			}
			walkSprites[dir][i] = sprite
		}
	}

	for _, dir := range game.Dirs() {
		actionSprites[dir] = make([]*Sprite, numActionFrames)

		for i := 0; i < numActionFrames; i++ {
			fileName := fmt.Sprintf("%s/%s/%s%d.png", baseDirectory, dirDirs[int(dir)], prefix[Action], i+1)
			sprite, err := NewSprite(fileName, width, height)
			if err != nil {
				// log.Fatal("Could not open ghost walk tile", fileName)
			}
			actionSprites[dir][i] = sprite
		}
	}

	return &HumanSprites{
		animWalk:   walkSprites,
		animAction: actionSprites,
	}
}

func LoadWallSprites() *WallSprites {
	wallSprites := make([]*Sprite, 8)
	width, height := float32(64), float32(64)
	for i := 0; i < 8; i++ {
		sprite, _ := NewSprite("data/floor/wall.png", width, height)
		wallSprites[i] = sprite
	}
	return &WallSprites{
		walls: wallSprites,
	}
}

type RenderData struct {
	wallSprites  *WallSprites
	ghostSprites *GhostSprites
	floorSprites *FloorSprites
	humanSprites *HumanSprites
	toolSprites  *ToolSprites
}

func LoadRenderData() *RenderData {
	return &RenderData{
		wallSprites:  LoadWallSprites(),
		ghostSprites: LoadGhostSprites(),
		humanSprites: LoadHumanSprites(),
		floorSprites: LoadFloorSprites(),
		toolSprites:  LoadToolSprites(),
	}
}

func ToWorldCoord(pos game.MapPosition) (float32, float32) {
	worldCellSize := tileSize
	return float32(pos.X()) * worldCellSize, float32(pos.Y()) * worldCellSize
}

func FloatPosToWorldCoord(pos game.Position) (float32, float32) {
	worldCellSize := tileSize
	return float32(pos.X()) * worldCellSize, float32(pos.Y()) * worldCellSize
}

var tileSize float32 = 0.8 * 64

func RenderMap(player game.Player, renderData *RenderData, g *game.Game) {
	wall := renderData.wallSprites.walls[0]
	floor := renderData.floorSprites.floor
	ghostStand := renderData.ghostSprites.animWalk
	// ghostAction := renderData.ghostSprites.animAction
	humanStand := renderData.humanSprites.animWalk
	humanAction := renderData.humanSprites.animAction
	boulderSprite := renderData.floorSprites.boulder
	bannWallSprite := renderData.floorSprites.ban
	doorSprite := renderData.floorSprites.door
	triggerA := renderData.toolSprites.triggersA
	triggerB := renderData.toolSprites.triggersB

	offset := float32(tileSize / 2.0)
	baseSize := float32(64)
	scaleMod := tileSize / baseSize

	for y := 0; y < g.Height(); y++ {
		for x := 0; x < g.Width(); x++ {
			pos := game.NewMapPosition(x, y)
			wx, wy := ToWorldCoord(pos)
			cell := g.Cell(pos)

			if !g.PlayerCanSeeCell(player, pos) {
				continue
			}

			if cell.IsWall() {
				wall.Draw(wx+offset, wy+offset, 0, scaleMod*1, false)
			} else {
				floor.Draw(wx+offset, wy+offset, 0, scaleMod*1, false)
			}
		}
	}

	for pos, bannWall := range g.BannWalls() {
		if !g.PlayerCanSeeBannWall(player, bannWall) {
			continue
		}
		wx, wy := ToWorldCoord(pos)
		bannWallSprite[bannWall.Type()].Draw(wx+offset, wy+offset, 0, scaleMod*1, true)
	}

	for pos, boulder := range g.Boulders() {
		if !g.PlayerCanSeeCell(player, pos) {
			continue
		}
		renderPos := g.BoulderRenderPos(boulder)
		wx, wy := FloatPosToWorldCoord(renderPos)

		boulderSprite.Draw(wx+offset, wy+offset, 0, scaleMod*1, true)
	}

	for pos, door := range g.Doors() {
		if !g.PlayerCanSeeDoor(player, door) {
			continue
		}
		if door.IsClosed() {
			wx, wy := ToWorldCoord(pos)
			doorSprite.Draw(wx+offset, wy+offset, 0, scaleMod*1, true)
		}
	}

	for pos, trigger := range g.Triggers() {
		if !g.PlayerCanSeeTrigger(player, trigger) {
			continue
		}
		cell := g.Cell(pos)
		if dir, err := cell.DirOfTrigger(trigger); err == nil {
			wx, wy := ToWorldCoord(pos)

			angle := float32(0)
			switch dir {
			case game.DirNorth:
				angle = math.Pi / 2.0
			case game.DirEast:
				angle = math.Pi
			case game.DirSouth:
				angle = 3 * math.Pi / 2.0
			case game.DirWest:
				angle = 0
			}

			if trigger.IsActive() {
				triggerA.Draw(wx+offset, wy+offset, angle, scaleMod*0.5, true)
			} else {
				triggerB.Draw(wx+offset, wy+offset, angle, scaleMod*0.5, true)
			}
		}
	}

	for _, otherPlayer := range g.Players() {
		if otherPlayer != player && !g.PlayerCanSeeOtherPlayer(player) {
			continue
		}
		renderPos := g.PlayerRenderPos(otherPlayer)
		wx, wy := FloatPosToWorldCoord(renderPos)
		if g.IsHuman(otherPlayer) {
			direction := g.PlayerDirection(otherPlayer)
			if g.PlayerIsWalking(otherPlayer) {
				// log.Println(g.PlayerWalkFrame(player))
				humanStand[direction][g.PlayerWalkFrame(otherPlayer)].Draw(wx+offset, wy+offset, 0, scaleMod*64/256.0, true)
			} else if g.PlayerDoesAction(otherPlayer) {
				humanAction[direction][g.PlayerActionFrame(otherPlayer)].Draw(wx+offset, wy+offset, 0, scaleMod*64/256.0, true)
			} else {
				humanStand[direction][0].Draw(wx+offset, wy+offset, 0, scaleMod*64/256.0, true)
			}
		} else if g.IsGhost(otherPlayer) {
			direction := g.PlayerDirection(otherPlayer)
			if g.PlayerIsWalking(otherPlayer) {
				// log.Println(g.PlayerWalkFrame(player))
				ghostStand[direction][g.PlayerWalkFrame(otherPlayer)].Draw(wx+offset, wy+offset, 0, scaleMod*64/256.0, true)
				// } else if g.PlayerDoesAction(player) {
				// 	ghostAction[direction][g.PlayerActionFrame(player)].Draw(wx+32, wy+32, 0, 64/256.0, true)
			} else {
				ghostStand[direction][0].Draw(wx+offset, wy+offset, 0, scaleMod*64/256.0, true)
			}
		} else {
			log.Fatal("To many players")
		}
	}
}

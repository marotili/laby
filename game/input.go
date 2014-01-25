package game

import (
	"github.com/banthar/Go-SDL/sdl"
	"time"
)

type Key int

const (
	KeyA Key = iota
	KeyW
	KeyD
	KeyS
	KeySpace
	KeyEnter
)

type InputState struct {
	keysDown map[Key]bool
	actions  []Action
	game     *Game
	player   Player
}

type ActionType int

const (
	ActionNoAction ActionType = iota
	ActionLookNorth
	ActionLookWest
	ActionLookSouth
	ActionLookEast
	ActionMoveNorth
	ActionMoveWest
	ActionMoveSouth
	ActionMoveEast
	ActionAction // press, push
	ActionToggleVisibility
	ActionPlayerReady
)

func NewInputState(game *Game, player Player) *InputState {
	return &InputState{
		keysDown: make(map[Key]bool, 6),
		actions:  make([]Action, 0, 10),
		game:     game,
		player:   player,
	}
}

type Action interface {
	Update(*InputState, time.Duration) (ActionType, bool, Action) // Return type: Action ok, continue?, next action
}

func (is *InputState) AddAction(action Action) {
	is.actions = append(is.actions, action)
}

func (is *InputState) SetKeyDown(k Key) {
	is.keysDown[k] = true
}

func (is *InputState) SetKeyUp(k Key) {
	delete(is.keysDown, k)
}

func (is *InputState) KeyDown(k Key) bool {
	if _, ok := is.keysDown[k]; ok {
		return true
	}
	return false
}

func (is *InputState) KeyUp(k Key) bool {
	return !is.KeyDown(k)
}

type KeyShortAction struct {
	dtime time.Duration // time since key down
	key   Key
}

type KeyLongAction struct {
	dtime time.Duration
	key   Key
}

type SpaceAction struct {
}

type EnterAction struct {
}

func NewKeyShortAction(key Key) *KeyShortAction {
	return &KeyShortAction{
		key:   key,
		dtime: 0,
	}
}

func NewKeyLongAction(key Key, dtime time.Duration) *KeyLongAction {
	return &KeyLongAction{
		key:   key,
		dtime: dtime,
	}
}

func NewSpaceAction() *SpaceAction {
	return &SpaceAction{}
}

func NewEnterAction() *EnterAction {
	return &EnterAction{}
}

const keyShortTimeDown time.Duration = 200 * time.Millisecond // ms
const keyLongTimeDown time.Duration = 50 * time.Millisecond   // ~ the time needed for the player to move a tile

func (ksa *KeyShortAction) Update(inputS *InputState, dt time.Duration) (ActionType, bool, Action) {
	if inputS.KeyUp(ksa.key) {
		switch ksa.key {
		case KeyA:
			return ActionLookWest, false, nil // action ok
		case KeyW:
			return ActionLookNorth, false, nil // action ok
		case KeyD:
			return ActionLookEast, false, nil // action ok
		case KeyS:
			return ActionLookSouth, false, nil // action ok
		}
	}

	ksa.dtime += dt

	if ksa.dtime >= keyShortTimeDown {
		var action ActionType = ActionNoAction
		switch ksa.key {
		case KeyA:
			action = ActionLookWest
		case KeyW:
			action = ActionLookNorth
		case KeyD:
			action = ActionLookEast
		case KeyS:
			action = ActionLookSouth
		}
		return action, true, NewKeyLongAction(ksa.key, keyLongTimeDown)
	}

	return ActionNoAction, true, ksa // keep action
}

func (kla *KeyLongAction) Update(inputS *InputState, dt time.Duration) (ActionType, bool, Action) {
	if inputS.KeyUp(kla.key) {
		return ActionNoAction, false, nil
	}

	kla.dtime += dt
	if kla.dtime >= keyLongTimeDown && !inputS.game.PlayerIsWalking(inputS.player) {
		var dir Direction
		switch kla.key {
		case KeyA:
			dir = DirWest
		case KeyW:
			dir = DirNorth
		case KeyD:
			dir = DirEast
		case KeyS:
			dir = DirSouth
		}

		next := NewKeyLongAction(kla.key, 0) // restart action - pulses every keyLongTimeDown ms

		target := inputS.game.playerState[inputS.player].mapPos.Neighbor(dir)
		if inputS.game.IsBoulder(target) {
			return ActionAction, true, next
		}

		switch kla.key {
		case KeyA:
			return ActionMoveWest, true, next // action ok
		case KeyW:
			return ActionMoveNorth, true, next // action ok
		case KeyD:
			return ActionMoveEast, true, next // action ok
		case KeyS:
			return ActionMoveSouth, true, next // action ok
		}
	}

	return ActionNoAction, true, kla // keep action
}

func (sa *SpaceAction) Update(inputS *InputState, dt time.Duration) (ActionType, bool, Action) {
	if inputS.KeyUp(KeySpace) {
		return ActionAction, false, nil
	}
	return ActionNoAction, true, sa
}

func (ea *EnterAction) Update(inputS *InputState, dt time.Duration) (ActionType, bool, Action) {
	if inputS.KeyUp(KeyEnter) {
		return ActionToggleVisibility, false, nil
	}

	return ActionNoAction, true, ea
}

func (is *InputState) HandleEvent(e *sdl.KeyboardEvent) {
	if e.Type == sdl.KEYDOWN {
		switch e.Keysym.Sym {
		case sdl.K_a:
			is.SetKeyDown(KeyA)
			is.AddAction(NewKeyShortAction(KeyA))
		case sdl.K_w:
			is.SetKeyDown(KeyW)
			is.AddAction(NewKeyShortAction(KeyW))
		case sdl.K_d:
			is.SetKeyDown(KeyD)
			is.AddAction(NewKeyShortAction(KeyD))
		case sdl.K_s:
			is.SetKeyDown(KeyS)
			is.AddAction(NewKeyShortAction(KeyS))
		case sdl.K_SPACE:
			is.SetKeyDown(KeySpace)
			is.AddAction(NewSpaceAction())
		case sdl.K_RETURN:
			is.SetKeyDown(KeyEnter)
			is.AddAction(NewEnterAction())
		}
	} else if e.Type == sdl.KEYUP {
		switch e.Keysym.Sym {
		case sdl.K_a:
			is.SetKeyUp(KeyA)
		case sdl.K_w:
			is.SetKeyUp(KeyW)
		case sdl.K_d:
			is.SetKeyUp(KeyD)
		case sdl.K_s:
			is.SetKeyUp(KeyS)
		case sdl.K_SPACE:
			is.SetKeyUp(KeySpace)
		case sdl.K_RETURN:
			is.SetKeyUp(KeyEnter)
		}
	}
}

func (is *InputState) StepActions(dt time.Duration) []ActionType {
	actions := make([]ActionType, 0, len(is.actions))

	nextActions := make([]Action, 0, len(is.actions))

	for _, action := range is.actions {
		actionType, _, nextAction := action.Update(is, dt)
		if actionType != ActionNoAction {
			actions = append(actions, actionType)
		}

		if nextAction != nil {
			nextActions = append(nextActions, nextAction)
		}
	}

	is.actions = nextActions

	return actions
}

// File scene_manager.go defines scene lifecycle primitives—Scene, SceneManager,
// and State—and implements a simple cross-fade transition between scenes.
package asteroids

import "github.com/hajimehoshi/ebiten/v2"

var (
	// transitionFrom is a scratch buffer for rendering the current scene
	// during transitions. Preallocated to avoid per-frame allocations.
	transitionFrom = ebiten.NewImage(ScreenWidth, ScreenHeight)

	// transiionTo is a scratch buffer for rendering the next scene
	// during transitions. Name preserved to match existing code.
	transiionTo = ebiten.NewImage(ScreenWidth, ScreenHeight)
)

// transitionMaxCount controls the duration (in frames) of the cross-fade.
const transitionMaxCount = 25

// Scene is the minimal contract for any drawable/updatable screen of the game.
type Scene interface {
	// Update advances scene logic by one tick.
	// Implementations may return a non-nil error to abort the game loop.
	Update(state *State) error

	// Draw renders the scene to the provided target image.
	Draw(screen *ebiten.Image)
}

// State bundles ambient runtime dependencies passed to scenes during Update.
type State struct {
	SceneManager *SceneManager // Enables a scene to initiate transitions.
	Input        *Input        // Per-frame input snapshot (may be nil if not provided).
}

// SceneManager owns the active scene and handles cross-fade transitions.
type SceneManager struct {
	current         Scene // Currently visible/active scene.
	next            Scene // Pending scene to transition into (if any).
	transitionCount int   // Frames remaining in the current transition, 0 when idle.
}

// Draw renders either the current scene alone or a cross-fade between
// the current scene (as the background) and the next scene (as the overlay).
func (s *SceneManager) Draw(r *ebiten.Image) {
	// If no transition is in progress, draw the current scene directly.
	if s.transitionCount == 0 {
		s.current.Draw(r)
		return
	}

	// During a transition, first draw both scenes into offscreen buffers.
	transitionFrom.Clear()
	s.current.Draw(transitionFrom)

	transiionTo.Clear()
	s.next.Draw(transiionTo)

	// Compose the transition: start with the "from" scene at full opacity.
	r.DrawImage(transitionFrom, nil)

	// Alpha increases from 0 -> 1 as transitionCount decreases from Max -> 0.
	alpha := 1 - float32(s.transitionCount)/float32(transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(alpha)

	// Draw the "to" scene on top, scaled by the computed alpha.
	r.DrawImage(transiionTo, op)
}

// Update advances the transition and delegates logic to the active scene.
//
// Behavior:
//   - When not transitioning, forwards Update to the current scene.
//   - While transitioning, decrements the transition timer until it reaches 0,
//     then swaps next into current.
//
// Note: Input is currently ignored here; scenes receive State without Input.
// Hook it up by populating the State.Input field from the caller when ready.
func (s *SceneManager) Update(_ *Input) error {
	// No transition: update the active scene.
	if s.transitionCount == 0 {
		return s.current.Update(&State{
			SceneManager: s,
			// Input: nil (not wired here); fill from caller when available.
		})
	}

	// Transition in progress: count down one frame.
	s.transitionCount--
	if s.transitionCount > 0 {
		return nil
	}

	// Transition finished: commit the scene swap and clear "next".
	s.current = s.next
	s.next = nil
	return nil
}

// GoToScene initiates a scene change.
//
// If there is no active scene yet, switches immediately without transition.
// Otherwise, schedules a cross-fade into the provided scene.
func (s *SceneManager) GoToScene(scene Scene) {
	if s.current == nil {
		// First scene: enter immediately.
		s.current = scene
	} else {
		// Defer switch via timed transition.
		s.next = scene
		s.transitionCount = transitionMaxCount
	}
}

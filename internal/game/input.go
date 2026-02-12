package game

// InputState contains player intent for one simulation frame.
type InputState struct {
	TurnLeft     bool
	TurnRight    bool
	Accelerating bool
	Firing       bool
}

package robot

/*
	Notes regarding gaits

	1) Yes, these patterns could be created algorithmically, but it would make the code harder to read.
	   It also makes my head hurt.
	2) Currently supported gaits:
		- Hexapods: 									tripod, ripple and wave gait.
		- <Odd number>pods (pentapods/heptapods):	    wave gait.
		- <Even number>pods (hexapod): 	                wave gait, ripple gait.
		  (gait type is still an argument for constructing these patterns, since someone just might come
		   up with some clever new gait)
	3) Yes, a symmetrical centipede with metachronal gait would be nice. Unfortunately I've run out of dynamixels
*/

import "fmt"

type GaitType int

const (
	TRIPOD GaitType = 0
	WAVE   GaitType = 1
	RIPPLE GaitType = 2
)

type GaitPattern [][]int

type Gait struct {
	Pattern                 *GaitPattern
	StanceReturnSpeedFactor float64
	Name                    string
	NumIndicesInPattern     int
}

func NewHeptapodGait(GaitType GaitType) (*Gait, error) {
	if GaitType != WAVE {
		return nil, fmt.Errorf("Nope. Not doing that. Give it time and you'll figure out why (...)")
	}
	p := make(GaitPattern, 49)

	copy(p, [][]int{
		{1, 0, 0, 0, 0, 0, 0}, // 1 == swing phase, 0 == stance phase
		{0, 1, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 1}})
	return &Gait{
		Pattern:                 &p,
		StanceReturnSpeedFactor: 0.16,
		Name:                    "Wave gait",
		NumIndicesInPattern:     7,
	}, nil

}

func NewPentapodGait(GaitType GaitType) (*Gait, error) {
	if GaitType != WAVE {
		return nil, fmt.Errorf("Nope. Not doing that. Give it time and you'll figure out why (...)")
	}

	p := make(GaitPattern, 25)

	copy(p, [][]int{
		{1, 0, 0, 0, 0}, // 1 == swing phase, 0 == stance phase
		{0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0},
		{0, 0, 0, 0, 1}})
	return &Gait{
		Pattern:                 &p,
		StanceReturnSpeedFactor: 0.2,
		Name:                    "Wave gait",
		NumIndicesInPattern:     5,
	}, nil
}

func NewHexapodGait(GaitType GaitType) (*Gait, error) {
	p := make(GaitPattern, 36)

	if GaitType == WAVE {
		copy(p, [][]int{{0, 0, 1, 0, 0, 0}, // 1 == swing phase, 0 == stance phase
			{0, 1, 0, 0, 0, 0},
			{1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1},
			{0, 0, 0, 0, 1, 0},
			{0, 0, 0, 1, 0, 0}})

		return &Gait{
			Pattern:                 &p,
			StanceReturnSpeedFactor: 0.19, // Yeah, there is a rounding error somewhere. Sue me.
			Name:                    "Wave gait",
			NumIndicesInPattern:     6,
		}, nil
	} else if GaitType == RIPPLE {
		copy(p, [][]int{{0, 0, 1, 0, 0, 1}, // 1 == swing phase, 0 == stance phase
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
			{1, 0, 0},
			{0, 1, 0}})

		return &Gait{
			Pattern:                 &p,
			StanceReturnSpeedFactor: 0.4,
			Name:                    "Ripple gait",
			NumIndicesInPattern:     3,
		}, nil
	}

	// default to TRIPOD
	copy(p, [][]int{
		{0, 1}, // 1 == swing phase, 0 == stance phase
		{1, 0},
		{0, 1},
		{1, 0},
		{0, 1},
		{1, 0},
		{0, 1},
		{1, 0}})
	return &Gait{
		Pattern:                 &p,
		StanceReturnSpeedFactor: 1,
		Name:                    "Tripod gait",
		NumIndicesInPattern:     2,
	}, nil
}

func NewGait(NumLegs int, GaitType GaitType) (*Gait, error) {

	switch NumLegs {
	case 6:
		return NewHexapodGait(GaitType)
	case 5:
		return NewPentapodGait(GaitType)
	}

	return nil, fmt.Errorf("missing gait definition for pod with %d legs. Please update gaits.go.", NumLegs)
}

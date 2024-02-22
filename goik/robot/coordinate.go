package robot

import "fmt"

type Coordinate struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
	Z float64 `json:"Z"`
}

func (c *Coordinate) String() string {
	return fmt.Sprintf("[%2.2f, %2.2f, %2.2f]", c.X, c.Y, c.Z)
}

func NewCoordinate(x float64, y float64, z float64) Coordinate {
	return Coordinate{X: x, Y: y, Z: z}
}

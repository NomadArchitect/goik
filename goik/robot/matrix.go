package robot

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

func HomogeneousTransformationMatrix(ProjectionMatrix *mat.Dense, Theta float64, Length float64) *mat.Dense {

	// Z Rotation matrix
	R_Z := mat.NewDense(3, 3, []float64{
		math.Cos(Theta), -math.Sin(Theta), 0,
		math.Sin(Theta), math.Cos(Theta), 0,
		0, 0, 1,
	})

	// Final rotation matrix
	var R mat.Dense
	R.Mul(R_Z, ProjectionMatrix)

	// Displacement vector
	D := mat.NewDense(3, 1, []float64{Length * math.Cos(Theta), Length * math.Sin(Theta), 0})

	// Homogeneous transformation Matrix (Composed from R & D)
	H := mat.NewDense(4, 4, []float64{
		R.At(0, 0), R.At(0, 1), R.At(0, 2), D.At(0, 0),
		R.At(1, 0), R.At(1, 1), R.At(1, 2), D.At(1, 0),
		R.At(2, 0), R.At(2, 1), R.At(2, 2), D.At(2, 0),
		0, 0, 0, 1,
	})

	return H
}

package robot

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// The number of joints in a leg (Coxa, Femur, Tibia, End effector)
const NUM_JOINTS = 4

// Offsets in leg joint array
const COXA_ORIGIN_INDEX = 0
const FEMUR_ORIGIN_INDEX = 1
const TIBIA_ORIGIN_INDEX = 2
const EFFECTOR_ORIGIN_INDEX = 3

// Instead of just attempting to move a servo from one location directly to another
// we create a series of intermediate steps. This also allows us to calculate an arc
// and also lift the end effector when the target position is in the same plane as
// the starting position
// type IntermediateAngles []ServoAngles

type IntermediateAngles struct {
	JointAngle [NUM_JOINTS - 1][INTERPOLATION_STEPS]float64
}

type IntermediateEffectorCoordinates [INTERPOLATION_STEPS]Coordinate

// SegmentLengths defines the lengths of each segment in the robot leg
// The length is calculated from the origin of one reference frame to
// the origin of the next reference frame in the kinematic chain.
type SegmentLengths struct {
	Coxa  float64 `json:"C"`
	Femur float64 `json:"F"`
	Tibia float64 `json:"T"`
}

type Leg struct {
	// Leg index is displayed in the simualtor views and
	// is also used to calculate the servo ID representing
	// a specific joint (please refer to the README.md file
	// for the numbering convention used).
	Index int
	// The robot legs are arranged around the body of the robot
	// Sepration angle defines the separation in degrees
	// between one leg and the next.
	CoxaSeparationAngle float64
	// The offset transformation matrix is used for calculating
	// the positioning of the Coxa reference frame origin in
	// the robot body base reference frame.
	OffsetTransformationMatrix *mat.Dense
	// All robot legs are not necessarily created equal
	// Most hexapods have identical leg topologies, but it
	// never hurts to prepare for other form factors :)
	SegmentLengths SegmentLengths
	// The Joints array contain the location of the reference
	// frame origin for each joint in the base reference frame
	// coordinate system.
	Joints [NUM_JOINTS]Coordinate
	// EffectorTarget is the target location for the leg's end effector
	// when it is moving.
	EffectorTarget Coordinate
	// Intermediate angles contain the set of target angles required
	// for performing a full step
	IntermediateAngles IntermediateAngles
	// Intermediate effector coordinates the set of target locations
	// for the end effector required for performing a full step
	// These do not necessarily describe a straigth line
	IntermediateEffectorCoordinates IntermediateEffectorCoordinates
	// ServoAngles represent the angles for the current state of the leg
	// (moving or stationary)
	ServoAngles ServoAngles
	// Swing interpolation index is the current index in the interpolation
	// table for the leg while it is in the swing phase
	swingInterpolationIndex int
	// Stance interpolation index is the current index in the interpolation
	// table for the leg while it is in the stance phase
	stanceInterpolationIndex float64
	// NeutralEffectorCoordinate defines the end effector's coordinate in
	// the base reference frame when the leg is in a neutral / rest position
	NeutralEffectorCoordinate Coordinate
	// IsReverting is indicating if the leg is reverting to a nutral/rest position
	IsReverting bool
	// RevertInterpolationIndex is the current index in the interpolation table for
	// moving the leg back to a neutral / rest position
	RevertInterpolationIndex int
	// Debug messages sent to the debug channel and will appear in the chat ui
	debugChannel chan string
}

// GetJointOrigin is a helper function for extracting coordinates from
// a homogeneus transformation matrix (after multiplication).
func (l *Leg) GetJointOrigin(H *mat.Dense) Coordinate {
	return Coordinate{X: H.At(0, 3), Y: H.At(1, 3), Z: H.At(2, 3)}
}

// RecalculateForwardKinematics is necessary for viewing the robot representation
// The robot center defines the base reference frame and each joint (including the end effector)
// has a reference frame.
// All rotations are done around the Z-axis of the reference frame and each reference frame
// is displaced from the previous frame.
// By starting with the base reference frame we define a displacement vector and a rotation
// matrix for determining the origin of the next reference frame in the kinematic chain
// (center of robot -> coxa origin -> femur origin -> tibia origin -> end effector)
// We then create a homogeneous tranformation matrix containing the rotation matrix and
// the displacement vector for each new reference frame in the chain.
// By multiplying these matrices together, we can extract the coordinates for each frame (coxa, femur, tibia, end effector)
// from the last column in the resulting 4x4 matrix. This gives us the data we need to represent the
// robot in a 2/3D view.
func (l *Leg) RecalculateForwardKinematics(angles ServoAngles) {
	l.ServoAngles = angles

	P_Femur := mat.NewDense(3, 3, []float64{1, 0, 0, 0, 1, 0, 0, 0, 1}) // Identity
	P_Coxa := mat.NewDense(3, 3, []float64{1, 0, 0, 0, 0, -1, 0, 1, 0})
	P_Tibia := mat.NewDense(3, 3, []float64{1, 0, 0, 0, 1, 0, 0, 0, 1}) // Identity

	H_Coxa := HomogeneousTransformationMatrix(P_Coxa, angles.Coxa*math.Pi/180.0, l.SegmentLengths.Coxa)
	H_Femur := HomogeneousTransformationMatrix(P_Femur, angles.Femur*math.Pi/180.0, l.SegmentLengths.Femur)
	H_Tibia := HomogeneousTransformationMatrix(P_Tibia, angles.Tibia*math.Pi/180.0, l.SegmentLengths.Tibia)

	var H0_1 mat.Dense
	var H1_2 mat.Dense
	var H2_3 mat.Dense

	H0_1.Mul(l.OffsetTransformationMatrix, H_Coxa)
	H1_2.Mul(&H0_1, H_Femur)
	H2_3.Mul(&H1_2, H_Tibia)

	l.Joints[COXA_ORIGIN_INDEX] = l.GetJointOrigin(l.OffsetTransformationMatrix)
	l.Joints[FEMUR_ORIGIN_INDEX] = l.GetJointOrigin(&H0_1)
	l.Joints[TIBIA_ORIGIN_INDEX] = l.GetJointOrigin(&H1_2)
	l.Joints[EFFECTOR_ORIGIN_INDEX] = l.GetJointOrigin(&H2_3)
	l.EffectorTarget = l.Joints[EFFECTOR_ORIGIN_INDEX]
}

// UpdateSwing moves the interpolation index to the next element
// in the interpolation table for the swing phase until it reaches
// the end. It then wraps around and starts from 0 again.
// Return values:
// 1: End of swing
// 0: Interpolating swing
func (l *Leg) UpdateSwing(direction Direction) bool {
	angles := NewServoAngles(l.IntermediateAngles.JointAngle[0][l.swingInterpolationIndex], l.IntermediateAngles.JointAngle[1][l.swingInterpolationIndex], l.IntermediateAngles.JointAngle[2][l.swingInterpolationIndex])
	l.RecalculateForwardKinematics(angles)

	// We need to lift the legs in the swing phase, so we will modify Z target slightly
	// when in the swing phase and then find a new solution where the leg is not touching the ground
	phase_step := math.Pi / (INTERPOLATION_STEPS - 1)
	z := l.Joints[EFFECTOR_ORIGIN_INDEX].Z

	// Phase should swing from 0 to pi
	phase := phase_step * (float64(l.swingInterpolationIndex))
	zNew := Z_LIFT * math.Sin(phase)
	l.ServoAngles, _ = SolveEffectorIK(l, NewCoordinate(l.Joints[EFFECTOR_ORIGIN_INDEX].X, l.Joints[EFFECTOR_ORIGIN_INDEX].Y, z-zNew), l.debugChannel)

	l.RecalculateForwardKinematics(l.ServoAngles)

	l.swingInterpolationIndex += int(direction)

	// Each leg has a dedicated interpolation table for each joint.
	if direction == Forward {
		if l.swingInterpolationIndex > INTERPOLATION_STEPS-1 {
			l.swingInterpolationIndex = 0
			l.stanceInterpolationIndex = INTERPOLATION_STEPS - 1
			return true
		}
	} else if direction == Reverse {
		if l.swingInterpolationIndex < 0 {
			l.swingInterpolationIndex = INTERPOLATION_STEPS - 1
			l.stanceInterpolationIndex = 0
			return true
		}
	}
	return false
}

// UpdateStance moves the interpolation index to the previous element
// in the interpolation table for the stance phase until it reaches
// the start. It then wraps around and starts from the end again.
func (l *Leg) UpdateStance(direction Direction, stanceReturnFactor float64) {
	angles := NewServoAngles(l.IntermediateAngles.JointAngle[0][int(l.stanceInterpolationIndex)], l.IntermediateAngles.JointAngle[1][int(l.stanceInterpolationIndex)], l.IntermediateAngles.JointAngle[2][int(l.stanceInterpolationIndex)])
	l.ServoAngles = angles
	l.RecalculateForwardKinematics(angles)

	l.stanceInterpolationIndex -= float64(direction) * stanceReturnFactor

	// Each leg has a dedicated interpolation table for each joint. The pod has the "master clock"
	// for indexing these tables
	if direction == Forward {
		if l.stanceInterpolationIndex < 0 {
			l.stanceInterpolationIndex = INTERPOLATION_STEPS - 1
			l.swingInterpolationIndex = 0
		}

	} else if direction == Reverse {
		if l.stanceInterpolationIndex > INTERPOLATION_STEPS-1 {
			l.stanceInterpolationIndex = 0
			l.swingInterpolationIndex = INTERPOLATION_STEPS - 1
		}
	}
}

// NewLeg returns a new leg
func NewLeg(
	// 0 - NUM_LEGS-1
	Index int,
	// Separation angle (in degrees) from previous leg
	CoxaSeparationAngle float64,
	// This defines the offset and rotation in relation to the center of the robot.
	// This also represent the start of the kinematic chain for the leg
	OffsetTransformationMatrix *mat.Dense,
	// Servo angles for the leg in rest / neutral position
	// (If all angles are zero, the kinematic chain forms a straight line)
	ServoAngles ServoAngles,
	// The distance between reference frames (coxa == distance from coxa reference frame origin to femur reference frame origin)
	SegmentLengths SegmentLengths,
	// Each servo has it's own unique id. Refer to the README.md file for the numbering scheme used
	ServoIds [NUM_JOINTS - 1]int, // 0: Coxa, 1: Femur, 2: Tibia
	// output channel for debug messages
	debugChannel chan string) *Leg {
	l := Leg{
		OffsetTransformationMatrix: OffsetTransformationMatrix,
		SegmentLengths:             SegmentLengths,
		Index:                      Index,
		ServoAngles:                ServoAngles,
		CoxaSeparationAngle:        CoxaSeparationAngle,
		debugChannel:               debugChannel,
	}

	l.RecalculateForwardKinematics(ServoAngles)

	l.NeutralEffectorCoordinate = l.Joints[EFFECTOR_ORIGIN_INDEX]

	return &l
}

// RevertToNutral updates the interpolation table with the steps
// necessary for moving the leg back from the current position to
// the nwutral / rest position
func (l *Leg) RevertToNutral() {
	targetCoordinate := l.NeutralEffectorCoordinate
	startCoordinate := l.Joints[EFFECTOR_ORIGIN_INDEX]

	stepX := (targetCoordinate.X - startCoordinate.X) / (INTERPOLATION_STEPS - 1)
	stepY := (targetCoordinate.Y - startCoordinate.Y) / (INTERPOLATION_STEPS - 1)
	stepZ := (targetCoordinate.Z - startCoordinate.Z) / (INTERPOLATION_STEPS - 1)

	for i := 0; i < INTERPOLATION_STEPS; i++ {
		swing := NewCoordinate(startCoordinate.X+stepX*float64(i), startCoordinate.Y+stepY*float64(i), startCoordinate.Z+stepZ*float64(i))

		angles, _ := SolveEffectorIK(l, swing, l.debugChannel)

		l.IntermediateAngles.JointAngle[0][i] = angles.Coxa
		l.IntermediateAngles.JointAngle[1][i] = angles.Femur
		l.IntermediateAngles.JointAngle[2][i] = angles.Tibia
	}
}

// UpdateRevert grounds all legs, so that the pod doesn't tip over
func (l *Leg) UpdateRevertPhase0() {

	angles := NewServoAngles(l.IntermediateAngles.JointAngle[0][l.RevertInterpolationIndex], l.IntermediateAngles.JointAngle[1][l.RevertInterpolationIndex], l.IntermediateAngles.JointAngle[2][l.RevertInterpolationIndex])
	l.RecalculateForwardKinematics(angles)
	l.ServoAngles, _ = SolveEffectorIK(l, NewCoordinate(l.Joints[EFFECTOR_ORIGIN_INDEX].X, l.Joints[EFFECTOR_ORIGIN_INDEX].Y, POD_Z_HEIGHT), l.debugChannel)
	l.RecalculateForwardKinematics(l.ServoAngles)
}

// UpdateRevert recalculates forward kinematic for the current interpolation step
// and increments the interpolation index.
func (l *Leg) UpdateRevertPhase1() int {

	angles := NewServoAngles(l.IntermediateAngles.JointAngle[0][l.RevertInterpolationIndex], l.IntermediateAngles.JointAngle[1][l.RevertInterpolationIndex], l.IntermediateAngles.JointAngle[2][l.RevertInterpolationIndex])

	l.RecalculateForwardKinematics(angles)

	l.RevertInterpolationIndex += 1

	if l.RevertInterpolationIndex >= INTERPOLATION_STEPS {
		l.IsReverting = false
		l.RevertInterpolationIndex = 0
		return l.Index + 1
	}

	// We need to lift the legs in the swing phase, so we will modify Z target slightly
	// when in the swing phase and then find a new solution where the leg is not touching the ground
	phase_step := math.Pi / (INTERPOLATION_STEPS - 1)
	z := l.Joints[EFFECTOR_ORIGIN_INDEX].Z

	// Phase should swing from 0 to pi
	phase := phase_step * (float64(l.RevertInterpolationIndex))
	zNew := REVERT_LIFT * math.Sin(phase)
	l.ServoAngles, _ = SolveEffectorIK(l, NewCoordinate(l.Joints[EFFECTOR_ORIGIN_INDEX].X, l.Joints[EFFECTOR_ORIGIN_INDEX].Y, z-zNew), l.debugChannel)

	l.RecalculateForwardKinematics(l.ServoAngles)

	return l.Index
}

// ResetInterpolator reset both the swing and stance interpolation indices
// for all ltegs to the center of their respective interpolation tables
func (l *Leg) ResetInterpolator() {
	l.swingInterpolationIndex = (INTERPOLATION_STEPS - 1) / 2
	l.stanceInterpolationIndex = float64(l.swingInterpolationIndex)
}

// Zero resets all servo angles in the leg to 0 degrees. This should result in the pod having all legs stretched
// out and each leg forming a straight line away from the robot body
// If it does not, you will have to mechanically adjust the robot body so that this condition is
// satisfied.
// This is a prerequisit for the FK/IK math to make sense in meat space ;)
func (l *Leg) Zero() {
	l.ServoAngles.Coxa = 0
	l.ServoAngles.Femur = 0
	l.ServoAngles.Tibia = 0
	l.RecalculateForwardKinematics(l.ServoAngles)
}

// Ground anchors the end effector to a given Z-height
func (l *Leg) Ground(height float64) error {
	var err error
	l.ServoAngles, err = SolveEffectorIK(l, NewCoordinate(l.Joints[EFFECTOR_ORIGIN_INDEX].X, l.Joints[EFFECTOR_ORIGIN_INDEX].Y, height), l.debugChannel)
	if err != nil {
		return err
	}
	l.RecalculateForwardKinematics(l.ServoAngles)
	return nil
}

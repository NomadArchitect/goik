package robot

// Note:
//
//	All servos are not created equal.
//	Some servos will have a center position at 0 degrees, others
//	may have the center position at a different angle
type ServoAngles struct {
	Coxa  float64 `json:"Coxa"`
	Femur float64 `json:"Femur"`
	Tibia float64 `json:"Tibia"`
}

func NewServoAngles(Coxa float64, Femur float64, Tibia float64) ServoAngles {
	return ServoAngles{Coxa: Coxa, Femur: Femur, Tibia: Tibia}
}

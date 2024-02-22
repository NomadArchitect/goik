package robot

import (
	"os"
)

// A MotionPrimitive consists of a set of joint motion sequences
type MotionPrimitive struct {
	rawAngles        []ServoAngles
	normalizedAngles []byte
}

func NewMotionPrimitive() *MotionPrimitive {
	return &MotionPrimitive{}
}

func (m *MotionPrimitive) Add(angles ServoAngles) {
	m.rawAngles = append(m.rawAngles, angles)
}

func (m *MotionPrimitive) Size() int {
	return len(m.rawAngles)
}

func (m *MotionPrimitive) Clear() {
	m.rawAngles = nil
	m.normalizedAngles = nil
}

func (m *MotionPrimitive) normalize(servoRange int, invertedCoxa bool, invertedFemur bool, invertedTibia bool) {
	for _, a := range m.rawAngles {

		var joint uint16
		if invertedCoxa {
			joint = 1024 - uint16((a.Coxa/float64(servoRange))*1024+512)
		} else {
			joint = uint16((a.Coxa/float64(servoRange))*1024 + 512)
		}
		m.normalizedAngles = append(m.normalizedAngles, uint8(joint&0xFF))
		m.normalizedAngles = append(m.normalizedAngles, uint8((joint&0xFF00)>>8))

		if invertedFemur {
			joint = 1024 - uint16((a.Femur/float64(servoRange))*1024+512)
		} else {
			joint = uint16((a.Femur/float64(servoRange))*1024 + 512)
		}
		m.normalizedAngles = append(m.normalizedAngles, uint8(joint&0xFF))
		m.normalizedAngles = append(m.normalizedAngles, uint8((joint&0xFF00)>>8))

		if invertedTibia {
			joint = 1024 - uint16((a.Tibia/float64(servoRange))*1024+512)
		} else {
			joint = uint16((a.Tibia/float64(servoRange))*1024 + 512)
		}
		m.normalizedAngles = append(m.normalizedAngles, uint8(joint&0xFF))
		m.normalizedAngles = append(m.normalizedAngles, uint8((joint&0xFF00)>>8))
	}
}

func (m *MotionPrimitive) createFile(path string) error {
	return os.WriteFile(path, m.normalizedAngles, 0644)
}

func (m *MotionPrimitive) Export(path string, servoRange int, invertedCoxa bool, invertedFemur bool, invertedTibia bool) error {
	m.normalize(servoRange, invertedCoxa, invertedFemur, invertedTibia)
	return m.createFile(path)
}

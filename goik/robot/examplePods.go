package robot

// Example pod with 6 legs. Uneven separation between legs
func NewExampleHexapod0() *BodyDefinition {
	gait, _ := NewHexapodGait(TRIPOD)
	b := &BodyDefinition{
		NumLegs:    6,
		CoxaAngles: []float64{0, 0, 180, 180, 180, 0},
		Gait:       gait,
	}
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{40, 0, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{40, 80, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-40, 80, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-40, 0, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-40, -80, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{40, -80, 0})

	for i := 0; i < b.NumLegs; i++ {
		b.RestAngles = append(b.RestAngles, ServoAngles{Coxa: 0, Femur: -50, Tibia: 100})
		b.Segments = append(b.Segments, SegmentLengths{Coxa: 30, Femur: 70, Tibia: 120})
	}
	return b
}

// Example pod with 6 legs. This can use tripod, ripple and wave gait
func NewExampleHexapod1() *BodyDefinition {
	gait, _ := NewHexapodGait(TRIPOD)
	b := &BodyDefinition{
		NumLegs:    6,
		CoxaAngles: []float64{0, 60, 120, 180, 240, 300},
		Gait:       gait,
	}
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{40, 0, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{20, 34.64, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-20, 34.64, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-40, 0, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-20, -34.64, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{20, -34.64, 0})

	for i := 0; i < b.NumLegs; i++ {
		b.RestAngles = append(b.RestAngles, ServoAngles{Coxa: 0, Femur: -50, Tibia: 100})
		b.Segments = append(b.Segments, SegmentLengths{Coxa: 30, Femur: 70, Tibia: 120})
	}

	return b
}

// Example pod with 6 legs, a relatively small body and short tibias
// This can use tripod, ripple and wave gait
func NewExampleHexapod2() *BodyDefinition {
	gait, _ := NewHexapodGait(TRIPOD)
	b := &BodyDefinition{
		NumLegs:    6,
		CoxaAngles: []float64{0, 60, 120, 180, 240, 300},
		Gait:       gait,
	}
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{40, 0, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{20, 34.64, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-20, 34.64, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-40, 0, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-20, -34.64, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{20, -34.64, 0})

	for i := 0; i < b.NumLegs; i++ {
		b.RestAngles = append(b.RestAngles, ServoAngles{Coxa: 0, Femur: 45, Tibia: 45})
		b.Segments = append(b.Segments, SegmentLengths{Coxa: 53.85, Femur: 48, Tibia: 61.7})
	}

	return b
}

// Example pod with 5 legs. The only valid gate is wave gait
func NewExamplePentapod() *BodyDefinition {
	gait, _ := NewPentapodGait(WAVE)
	b := &BodyDefinition{
		NumLegs:    5,
		CoxaAngles: []float64{0, 72, 144, 216, 288, 360},
		Gait:       gait,
	}

	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{40.00, 0.00, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{12.36, 38.04, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-32.36, 23.51, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-32.36, -23.51, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{12.36, -38.04, 0})

	for i := 0; i < b.NumLegs; i++ {
		b.RestAngles = append(b.RestAngles, ServoAngles{Coxa: 0, Femur: -50, Tibia: 100})
		b.Segments = append(b.Segments, SegmentLengths{Coxa: 40, Femur: 60, Tibia: 150})
	}

	return b
}

// Example pod with 7 legs. The only valid gate is wave gait
func NewHeptapod() *BodyDefinition {
	gait, _ := NewHeptapodGait(WAVE)
	b := &BodyDefinition{
		NumLegs:    7,
		CoxaAngles: []float64{51.43, 102.86, 154.29, 205.72, 257.15, 308.58, 360},

		Gait: gait,
	}

	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{24.94, 31.27, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-8.90, 39.00, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-36.04, 17.36, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-36.04, -17.36, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-8.90, -39.00, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{24.94, -31.27, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{40.00, 0.00, 0})

	for i := 0; i < b.NumLegs; i++ {
		b.RestAngles = append(b.RestAngles, ServoAngles{Coxa: 0, Femur: -50, Tibia: 100})
		b.Segments = append(b.Segments, SegmentLengths{Coxa: 40, Femur: 60, Tibia: 150})
	}

	return b
}

// Example (spider like) pod with 8 legs and varying segment lengths
func NewSpider() *BodyDefinition {
	gait, _ := NewHexapodGait(TRIPOD)
	b := &BodyDefinition{
		NumLegs:    8,
		CoxaAngles: []float64{20, 70, 110, 150, 210, 250, 290, 330},
		Gait:       gait,
	}

	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{10, 10, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{15, 40, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-15, 40, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-10, 10, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-10, -10, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{-10, -30, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{10, -30, 0})
	b.CoxaCoordinates = append(b.CoxaCoordinates, Coordinate{10, -10, 0})

	b.Segments = append(b.Segments, SegmentLengths{Coxa: 20, Femur: 40, Tibia: 60})
	b.Segments = append(b.Segments, SegmentLengths{Coxa: 40, Femur: 70, Tibia: 120})
	b.Segments = append(b.Segments, SegmentLengths{Coxa: 40, Femur: 70, Tibia: 120})
	b.Segments = append(b.Segments, SegmentLengths{Coxa: 20, Femur: 40, Tibia: 60})
	b.Segments = append(b.Segments, SegmentLengths{Coxa: 20, Femur: 40, Tibia: 60})
	b.Segments = append(b.Segments, SegmentLengths{Coxa: 40, Femur: 70, Tibia: 120})
	b.Segments = append(b.Segments, SegmentLengths{Coxa: 40, Femur: 70, Tibia: 120})
	b.Segments = append(b.Segments, SegmentLengths{Coxa: 20, Femur: 40, Tibia: 60})

	for i := 0; i < b.NumLegs; i++ {
		b.RestAngles = append(b.RestAngles, ServoAngles{Coxa: 0, Femur: -50, Tibia: 100})
	}

	return b
}

package views

import (
	"GOIK/robot"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type XyView struct {
	x            float32
	y            float32
	size         float32
	legendOffset float32
}

func NewXyView(x float32, y float32, size float32) *XyView {
	return &XyView{x: x, y: y, size: size, legendOffset: 10}
}

func (v *XyView) TranslateX(x float64) int {
	return int(v.x) + int(v.size)/2 - int(v.legendOffset) + int(x)
}

func (v *XyView) TranslateY(y float64) int {
	return int(v.y) + int(v.size)/2 - int(v.legendOffset) + int(y)
}

var effectorTargetUpdateThreshold = 20
var effectorTargetUpdateTicker = 0

func (v *XyView) Render(screen *ebiten.Image, p *robot.Pod) {

	DrawFrame(screen, "XY View ", v.size, v.x, v.y, v.legendOffset)

	// Draw body frame
	for l := 0; l < p.BodyDefinition.NumLegs-1; l++ {
		vector.StrokeLine(screen,
			float32(v.TranslateX(p.Legs[l].Joints[0].X)),
			float32(v.TranslateY(p.Legs[l].Joints[0].Y)),
			float32(v.TranslateX(p.Legs[l+1].Joints[0].X)),
			float32(v.TranslateY(p.Legs[l+1].Joints[0].Y)),
			5,
			White(),
			true)
	}
	vector.StrokeLine(screen,
		float32(v.TranslateX(p.Legs[p.BodyDefinition.NumLegs-1].Joints[0].X)),
		float32(v.TranslateY(p.Legs[p.BodyDefinition.NumLegs-1].Joints[0].Y)),
		float32(v.TranslateX(p.Legs[0].Joints[0].X)),
		float32(v.TranslateY(p.Legs[0].Joints[0].Y)),
		5,
		White(),
		true)

	// Draw Coxa, Femur and Tibia
	for j := 0; j < robot.NUM_JOINTS-1; j++ {
		for l := 0; l < p.BodyDefinition.NumLegs; l++ {
			col := White()
			width := 3
			if p.IsSwingPhase(l) {
				col = Blue()
			}
			vector.StrokeLine(screen,
				float32(v.TranslateX(p.Legs[l].Joints[j].X)),
				float32(v.TranslateY(p.Legs[l].Joints[j].Y)),
				float32(v.TranslateX(p.Legs[l].Joints[j+1].X)),
				float32(v.TranslateY(p.Legs[l].Joints[j+1].Y)),
				float32(width),
				col,
				true)
		}
	}

	// Draw joints
	for _, l := range p.Legs {
		for _, j := range l.Joints {
			vector.DrawFilledCircle(screen, float32(v.TranslateX(j.X)), float32(v.TranslateY(j.Y)), 5, Red(), true)
		}
	}

	// Annotate legs with lex indexes
	for l := 0; l < p.BodyDefinition.NumLegs; l++ {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", l), v.TranslateX(p.Legs[l].Joints[3].X+15), v.TranslateY(p.Legs[l].Joints[3].Y+5))
	}

	// Draw IK effector targets
	targetSize := 5
	effectorTargetUpdateTicker++
	if effectorTargetUpdateTicker > effectorTargetUpdateThreshold {
		if effectorTargetUpdateTicker > 2*effectorTargetUpdateThreshold {
			effectorTargetUpdateTicker = 0
		}
		targetSize = 10
	}
	for i, l := range p.Legs {
		t := l.EffectorTarget
		vector.DrawFilledCircle(screen, float32(v.TranslateX(t.X)), float32(v.TranslateY(t.Y)), float32(targetSize), Blue(), true)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), v.TranslateX(t.X)+10, v.TranslateY(t.Y))
	}

	if p.HasDefinedStride {
		for l := range p.Legs {
			for i := range p.Legs[l].IntermediateEffectorCoordinates {
				c := p.Legs[l].IntermediateEffectorCoordinates[i]
				vector.DrawFilledCircle(screen, float32(v.TranslateX(c.X)), float32(v.TranslateY(c.Y)), float32(1), White(), false)
			}
		}
	}
}

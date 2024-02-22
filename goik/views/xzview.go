package views

import (
	"GOIK/robot"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type XzView struct {
	x            float32
	y            float32
	size         float32
	legendOffset float32
}

func NewXzView(x float32, y float32, size float32) *XzView {
	return &XzView{x: x, y: y, size: size, legendOffset: 10}
}

func (v *XzView) TranslateX(x float64) int {
	return int(v.x) + int(v.size)/2 - int(v.legendOffset) + int(x)
}

func (v *XzView) TranslateY(y float64) int {
	return int(v.size)/2 - int(v.legendOffset) + int(y)
}

func (v *XzView) Render(screen *ebiten.Image, p *robot.Pod) {
	DrawFrame(screen, "XZ View", v.size, v.x, v.y, v.legendOffset)

	// Draw body frame
	for l := 0; l < p.BodyDefinition.NumLegs-1; l++ {
		vector.StrokeLine(screen,
			float32(v.TranslateX(p.Legs[l].Joints[0].X)),
			float32(v.TranslateY(p.Legs[l].Joints[0].Z)),
			float32(v.TranslateX(p.Legs[l+1].Joints[0].X)),
			float32(v.TranslateY(p.Legs[l+1].Joints[0].Z)),
			5,
			White(),
			true)
	}
	vector.StrokeLine(screen,
		float32(v.TranslateX(p.Legs[p.BodyDefinition.NumLegs-1].Joints[0].X)),
		float32(v.TranslateY(p.Legs[p.BodyDefinition.NumLegs-1].Joints[0].Z)),
		float32(v.TranslateX(p.Legs[0].Joints[0].X)),
		float32(v.TranslateY(p.Legs[0].Joints[0].Z)),
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
				float32(v.TranslateY(p.Legs[l].Joints[j].Z)),
				float32(v.TranslateX(p.Legs[l].Joints[j+1].X)),
				float32(v.TranslateY(p.Legs[l].Joints[j+1].Z)),
				float32(width),
				col,
				true)
		}
	}

	for _, l := range p.Legs {
		for _, j := range l.Joints {
			vector.DrawFilledCircle(screen, float32(v.TranslateX(j.X)), float32(v.TranslateY(j.Z)), 5, Red(), true)
		}
	}

}

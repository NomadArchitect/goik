package views

import (
	"GOIK/robot"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func Red() color.Color {
	return color.RGBA{255, 0, 0, 1}
}

func Blue() color.Color {
	return color.RGBA{0, 0, 255, 1}
}

func White() color.Color {
	return color.RGBA{255, 255, 255, 1}
}

func SwingPassiveClr() color.Color {
	return color.RGBA{200, 200, 200, 1}
}

func SwingActiveClr() color.Color {
	return color.RGBA{255, 255, 255, 1}
}

func StancePassiveClr() color.Color {
	return color.RGBA{32, 32, 32, 1}
}

func StanceActiveClr() color.Color {
	return color.RGBA{64, 64, 64, 1}
}

type View interface {
	Render(screen *ebiten.Image, p *robot.Pod)
}

type RenderViews []View

func DrawFrame(screen *ebiten.Image, title string, size float32, x float32, y float32, titleOffset float32) {
	framecolor := color.RGBA{64, 64, 64, 1}
	vector.StrokeRect(screen, x, y, size, size, 2, framecolor, false)
	ebitenutil.DebugPrintAt(screen, title, int(x+titleOffset), int(y)+int(titleOffset))
}

func DrawAxis(screen *ebiten.Image, horizontalLegend string, verticalLegend string, size float32, x float32, y float32, legendOffset float32) {
	// Axis
	vector.StrokeLine(screen, x+legendOffset, y+size/2-legendOffset, x+size-legendOffset, y+size/2-legendOffset, 1, color.White, false)
	vector.StrokeLine(screen, x+size-legendOffset, y+legendOffset, x+size-legendOffset, y+size-legendOffset, 1, color.White, false)

	// Legend
	ebitenutil.DebugPrintAt(screen, horizontalLegend, int(x+legendOffset), int(y+size/2-legendOffset-20))
	ebitenutil.DebugPrintAt(screen, verticalLegend, int(x+size-legendOffset-20), int(y+legendOffset))
}

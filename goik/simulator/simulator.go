// Copyright 2025 Hans Jørgen Grimstad
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package simulator

import (
	"GOIK/comms"
	"GOIK/robot"
	"GOIK/views"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const window_size = 1024

var networkcontroller comms.NetworkController

type Game struct {
	views views.RenderViews
	Shell *Shell
}

func NewGame(s *Shell) *Game {
	g := Game{Shell: s}

	g.views = append(g.views, views.NewXzView(0, 0, window_size/2))
	g.views = append(g.views, views.NewXyView(window_size/2, 0, window_size/2))
	g.views = append(g.views, views.NewIsoView(window_size/2, window_size/2, window_size/2))
	g.views = append(g.views, views.NewGaitView(0, window_size/2, window_size/2))

	return &g
}

var DELAY_COUNTER int = 10
var counter int = 0

func (g *Game) Update() error {

	counter++
	if counter >= 3*DELAY_COUNTER {
		g.Shell.Pod.Update()
		networkcontroller.Update()
		counter = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, v := range g.views {
		v.Render(screen, g.Shell.Pod)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 1024
}

func Run() {

	// Create the pod body and define a default gait
	pod := robot.NewPod(robot.NewExampleHexapod2())

	// Create the command shell
	shell := NewShell(pod)
	go shell.Run()
	g := NewGame(shell)

	// Create a robot network controller
	networkcontroller = *comms.NewNetworkController(1, pod, shell.outputCh)

	// Create main window and start the simulation
	ebiten.SetWindowSize(window_size, window_size)
	ebiten.SetWindowTitle("Pod simulator")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

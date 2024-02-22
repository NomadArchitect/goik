package main

import (
	"GOIK/simulator"
)

func main() {
	// go tool pprof cpu.pprof
	// defer profile.Start(profile.ProfilePath(".")).Stop()
	simulator.Run()
}

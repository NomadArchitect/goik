package simulator

import (
	"GOIK/robot"
	"log"
	"strings"

	"github.com/borud/chatui"
)

type dispatchFunc func(args []string) error

type Shell struct {
	Pod         *robot.Pod
	outputCh    chan string
	commandCh   chan string
	dispatchMap map[string]dispatchFunc
}

func NewShell(pod *robot.Pod) *Shell {
	s := Shell{Pod: pod, outputCh: make(chan string, 10), commandCh: make(chan string)}

	s.Pod.SetDebugChannel(s.outputCh)

	s.dispatchMap = map[string]dispatchFunc{
		"help":             s.executeHelpCmd,
		"effectors":        s.executeEffectorsCmd,
		"set_coxa_length":  s.executeSetCoxaLengthCmd,
		"set_femur_length": s.executeSetFemurLengthCmd,
		"set_tibia_length": s.executeSetTibiaLengthCmd,
		"set_coxa_angle":   s.executeSetCoxaAngleCmd,
		"set_femur_angle":  s.executeSetFemurAngleCmd,
		"set_tibia_angle":  s.executeSetTibiaAngleCmd,
		"stride_vector":    s.executeStrideVectorCmd,
		"stride_angle":     s.executeStrideAngleCmd,
		"start":            s.executeStartCmd,
		"stop":             s.executeStopCmd,
		"reset":            s.executeResetCmd,
		"speed":            s.executeSpeedCmd,
		"zlift":            s.executeZLiftCmd,
		"gait":             s.executeGaitCmd,
		"open":             s.executeOpenServoPortCmd,
		"close":            s.executeCloseServoPortCmd,
		"save":             s.executeSaveCmd,
		"load":             s.executeLoadCmd,
		"zero":             s.executeZeroCmd,
		"reverse":          s.executeReverseCmd,
		"revert":           s.executeRevertCmd,
		"record":           s.executeRecordCmd,
		"export":           s.executeExportCmd,
		"debug":            s.executeDebugCmd,
		"step":             s.executeStepCycleCmd,
		"pitch":            s.executePitchCmd,
		"yaw":              s.executeYawCmd,
		"roll":             s.executeRollCmd,
		"up":               s.executeUpCmd,
		"down":             s.executeDownCmd,
		"ground":           s.executeGroundCmd,
	}

	return &s
}

func (s *Shell) Run() {

	chatui := chatui.New(chatui.Config{
		OutputCh:     s.outputCh,
		CommandCh:    s.commandCh,
		DynamicColor: false,
		BlockCtrlC:   true,
		HistorySize:  10,
	})

	s.outputCh <- "Pod playground"

	go func() {
		for command := range s.commandCh {
			if strings.ToLower(command) == "/quit" {
				chatui.Stop()
			}
			err := s.Dispatch(command)
			if err != nil {
				s.outputCh <- err.Error()
			}
			chatui.SetStatus("last command was: " + command)
		}
	}()

	go func() {
		// this is done in a goroutine because it will block if the UI is not running.
		chatui.SetStatus("type /quit to exit")

		s.outputCh <- "-----------------------------------"
		s.outputCh <- "Hexapod Designer"
		s.outputCh <- "Copyright Hans Jørgen Grimstad 2024"
		s.outputCh <- "www.TimeExpander.com"
		s.outputCh <- "-----------------------------------"
		s.outputCh <- ""
		s.outputCh <- "Type 'help' for a list of commands"

	}()

	err := chatui.Run()
	if err != nil {
		log.Fatal(err)
	}
}

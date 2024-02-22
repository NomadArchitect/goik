package simulator

import (
	"GOIK/robot"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const POD_FOLDER = "pods"
const PRIMITIVES_FOLDER = "primitives"

func (s *Shell) executeHelpCmd(args []string) error {
	s.outputCh <- "Commands:"
	s.outputCh <- "\teffectors                                  - output current end effector positions."
	s.outputCh <- "\tgait <tripod | ripple | wave>              - select new gait."
	s.outputCh <- "\tset_coxa_length <ALL | legNum> <length>"
	s.outputCh <- "\tset_femur_length <ALL | legNum> <length>"
	s.outputCh <- "\tset_tibia_length <ALL | legNum> <length>"
	s.outputCh <- "\tset_coxa_angle <ALL | legNum> <angle>"
	s.outputCh <- "\tset_femur_angle <ALL | legNum> <angle>"
	s.outputCh <- "\tset_tibia_angle <ALL | legNum> <angle>"
	s.outputCh <- "\tground <height>                            - Grounds all end effectors and updates rest angles"
	s.outputCh <- "\tstride_vector <nrepeats> <x> <y>           - Set direction. x & y are relative to current location"
	s.outputCh <- "\tstride_angle <nrepeats> <degrees>          - Rotate around center of gravity"
	s.outputCh <- "\tpitch <degrees>                            - Pitch move"
	s.outputCh <- "\tyaw <degrees>                              - Yaw move"
	s.outputCh <- "\troll <degrees>                             - Roll move"
	s.outputCh <- "\tup <z>                                     - Stand tall"
	s.outputCh <- "\tdown <z>                                   - Low rider"
	s.outputCh <- "\tstart                                      - Start pod"
	s.outputCh <- "\tstop                                       - Stop pod"
	s.outputCh <- "\treset <1|2|3|4|5>                          - Reset to design preset <n>"
	s.outputCh <- "\tspeed                                      - speed <1-10>"
	s.outputCh <- "\tzlift                                      - defines leg lift during swing phase"
	s.outputCh <- "\topen <IP:port>                             - open connection to dynamixel  UDP bridge"
	s.outputCh <- "\tclose                                      - close dynamixel connection"
	s.outputCh <- "\tsave <filename>                            - save pod definition to file"
	s.outputCh <- "\tload <filename>                            - load pod definition from file"
	s.outputCh <- "\tzero                                       - Aligns all servos to zero degrees"
	s.outputCh <- "\treverse                                    - Reverses walking direction"
	s.outputCh <- "\tstep                                       - Performs a single cycle through a gait pattern"
	s.outputCh <- "\trevert                                     - Revert to a neutral position"
	s.outputCh <- "\trecord <on|off>                            - Records next run or stops recording"
	s.outputCh <- "\texport <file> <max deg> <mask>             - Save recording to a file. Servo range: 180-360."
	s.outputCh <- "\t                                             mask is of the format \"100\", where a \"1\""
	s.outputCh <- "\t                                             signifies that the servo horn is pointing in negative Z"
	s.outputCh <- "\t                                             and a \"0\" that it is pointing in positive Z direction"
	s.outputCh <- "\t                                             The bitmask order is coxa, femur, tibia"

	return nil
}

func (s *Shell) executeEffectorsCmd(args []string) error {
	s.outputCh <- "Current end effector positions:"

	positions := s.Pod.GetEndEffectorPositions()
	for _, p := range positions {
		s.outputCh <- p.String()
	}

	return nil
}

func (s *Shell) executeSetCoxaLengthCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 3 {
		return fmt.Errorf("syntax error ('set_coxa_length <legnum> <length>'): %+v", args)
	}

	length, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('set_coxa_length <legnum> <length>'): %+v", args)
	}

	if strings.ToUpper(args[1]) == "ALL" {
		for _, l := range s.Pod.Legs {
			s.Pod.SetCoxaLength(l.Index, length)
			s.outputCh <- fmt.Sprintf("Changing coxa length of leg %d to %2.2f", l.Index, length)
		}
	} else {
		legnum, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			return fmt.Errorf("syntax error ('set_coxa_length <legnum> <length>'): %+v", args)
		}
		if legnum < 0 || legnum >= int64(s.Pod.BodyDefinition.NumLegs) {
			return fmt.Errorf("invalid leg index. (Pod has %d legs. Indexing is 0 based)", s.Pod.BodyDefinition.NumLegs)
		}
		s.outputCh <- fmt.Sprintf("Changing coxa length of leg %d to %2.2f", legnum, length)
		s.Pod.SetCoxaLength(int(legnum), length)
	}
	return nil
}

func (s *Shell) executeSetFemurLengthCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)
	if len(args) != 3 {
		return fmt.Errorf("syntax error ('set_femur_length <legnum> <length>'): %+v", args)
	}

	length, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('set_femur_length <legnum> <length>'): %+v", args)
	}

	if strings.ToUpper(args[1]) == "ALL" {
		for _, l := range s.Pod.Legs {
			s.Pod.SetFemurLength(l.Index, length)
			s.outputCh <- fmt.Sprintf("Changing femur length of leg %d to %2.2f", l.Index, length)
		}
	} else {
		legnum, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			return fmt.Errorf("syntax error ('set_femur_length <legnum> <length>'): %+v", args)
		}
		if legnum < 0 || legnum >= int64(s.Pod.BodyDefinition.NumLegs) {
			return fmt.Errorf("invalid leg index. (Pod has %d legs. Indexing is 0 based)", s.Pod.BodyDefinition.NumLegs)
		}
		s.outputCh <- fmt.Sprintf("Changing femur length of leg %d to %2.2f", legnum, length)
		s.Pod.SetFemurLength(int(legnum), length)
	}

	return nil
}

func (s *Shell) executeSetTibiaLengthCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)
	if len(args) != 3 {
		return fmt.Errorf("syntax error ('set_tibia_length <legnum> <length>'): %+v", args)
	}

	length, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('set_tibia_length <legnum> <length>'): %+v", args)
	}

	if strings.ToUpper(args[1]) == "ALL" {
		for _, l := range s.Pod.Legs {
			s.Pod.SetTibiaLength(l.Index, length)
			s.outputCh <- fmt.Sprintf("Changing tibia length of leg %d to %2.2f", l.Index, length)
		}
	} else {
		legnum, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			return fmt.Errorf("syntax error ('set_tibia_length <legnum> <length>'): %+v", args)
		}
		if legnum < 0 || legnum >= int64(s.Pod.BodyDefinition.NumLegs) {
			return fmt.Errorf("invalid leg index. (Pod has %d legs. Indexing is 0 based)", s.Pod.BodyDefinition.NumLegs)
		}
		s.outputCh <- fmt.Sprintf("Changing tibia length to %2.2f", length)
		s.Pod.SetTibiaLength(int(legnum), length)
	}
	return nil
}

func (s *Shell) executeSetCoxaAngleCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)
	if len(args) != 3 {
		return fmt.Errorf("syntax error ('set_coxa_angle <legnum> <angle>'): %+v", args)
	}

	angle, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('set_coxa_angle  <legnum> <angle>'): %+v", args)
	}

	if strings.ToUpper(args[1]) == "ALL" {
		for _, l := range s.Pod.Legs {
			s.Pod.SetCoxaAngle(l.Index, angle)
			s.outputCh <- fmt.Sprintf("Changing coxa angle of leg %d to %2.2f", l.Index, angle)
		}
	} else {
		legnum, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			return fmt.Errorf("syntax error ('set_coxa_angle  <legnum> <angle>'): %+v", args)
		}
		if legnum < 0 || legnum >= int64(s.Pod.BodyDefinition.NumLegs) {
			return fmt.Errorf("invalid leg index. (Pod has %d legs. Indexing is 0 based)", s.Pod.BodyDefinition.NumLegs)
		}
		s.outputCh <- fmt.Sprintf("Changing coxa angle of leg %d to %2.2f", legnum, angle)
		s.Pod.SetCoxaAngle(int(legnum), angle)
	}
	return nil
}

func (s *Shell) executeSetFemurAngleCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)
	if len(args) != 3 {
		return fmt.Errorf("syntax error ('set_femur_angle  <legnum> <angle>'): %+v", args)
	}

	angle, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('set_femur_angle  <legnum> <angle>'): %+v", args)
	}

	if strings.ToUpper(args[1]) == "ALL" {
		for _, l := range s.Pod.Legs {
			s.Pod.SetFemurAngle(l.Index, angle)
			s.outputCh <- fmt.Sprintf("Changing femur angle of leg %d to %2.2f", l.Index, angle)
		}
	} else {
		legnum, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			return fmt.Errorf("syntax error ('set_femur_angle  <legnum> <angle>'): %+v", args)
		}
		if legnum < 0 || legnum >= int64(s.Pod.BodyDefinition.NumLegs) {
			return fmt.Errorf("invalid leg index. (Pod has %d legs. Indexing is 0 based)", s.Pod.BodyDefinition.NumLegs)
		}
		s.outputCh <- fmt.Sprintf("Changing femur angle of leg %d to %2.2f", angle)
		s.Pod.SetFemurAngle(int(legnum), angle)
	}

	return nil
}

func (s *Shell) executeSetTibiaAngleCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)
	if len(args) != 3 {
		return fmt.Errorf("syntax error ('set_tibia_angle <angle>'): %+v", args)
	}

	angle, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('set_tibia_angle <angle>'): %+v", args)
	}

	if strings.ToUpper(args[1]) == "ALL" {
		for _, l := range s.Pod.Legs {
			s.Pod.SetTibiaAngle(l.Index, angle)
			s.outputCh <- fmt.Sprintf("Changing tibia angle of leg %d to %2.2f", l.Index, angle)
		}
	} else {
		legnum, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			return fmt.Errorf("syntax error ('set_tibia_angle <angle>'): %+v", args)
		}
		if legnum < 0 || legnum >= int64(s.Pod.BodyDefinition.NumLegs) {
			return fmt.Errorf("invalid leg index. (Pod has %d legs. Indexing is 0 based)", s.Pod.BodyDefinition.NumLegs)
		}
		s.outputCh <- fmt.Sprintf("Changing tibia angle og leg %d to %2.2f", legnum, angle)
		s.Pod.SetTibiaAngle(int(legnum), angle)
	}
	return nil
}

func (s *Shell) Dispatch(command string) error {
	command = strings.TrimSpace(command)
	for key, value := range s.dispatchMap {
		if len(command) >= len(key) && command[:len(key)] == key {
			return value(strings.Split(command, " "))
		}
	}

	return fmt.Errorf("unknown command: '%s'", command)
}

func (s *Shell) executeStrideVectorCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 4 {
		return fmt.Errorf("syntax error ('stride_vector <nrepeats> <X> <Y>'): %+v", args)
	}

	repeats, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return fmt.Errorf("syntax error ('stride_vector <nrepeats> <X> <Y>'): %+v", args)
	}

	deltaX, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('stride_vector <nrepeats> <X> <Y>'): %+v", args)
	}

	deltaY, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('stride_vector <nrepeats> <X> <Y>'): %+v", args)
	}

	// if s.Pod.IsWalking {
	// 	return fmt.Errorf("live direction transitions aren't implemented yet. Please stop and reset before switching direction")
	// }

	s.Pod.ResetInterpolator()

	return s.Pod.SetStrideVector(int(repeats), deltaX, deltaY)
}

func (s *Shell) executeStrideAngleCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 3 {
		return fmt.Errorf("syntax error ('stride_angle <nrepeats> <degrees>'): %+v", args)
	}

	repeats, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return fmt.Errorf("syntax error ('stride_angle <nrepeats> <degrees>'): %+v", args)
	}

	degrees, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('stride_angle <nrepeats> <degrees>'): %+v", args)
	}

	// if s.Pod.IsWalking {
	// 	return fmt.Errorf("live direction transitions aren't implemented yet. Please stop and reset before switching direction")
	// }

	s.Pod.ResetInterpolator()

	err = s.Pod.SetRotation(int(repeats), degrees)
	if err != nil {
		return err
	}
	return nil
}

func (s *Shell) executeStartCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 1 {
		return fmt.Errorf("syntax error ('start'): %+v", args)
	}

	if !s.Pod.HasDefinedStride {
		return fmt.Errorf("no defined stride")
	}
	err := s.Pod.Start()
	networkcontroller.Start()

	if err != nil {
		return err
	}

	return nil
}

func (s *Shell) executeStopCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 1 {
		return fmt.Errorf("syntax error ('stop'): %+v", args)
	}

	s.Pod.Stop()
	return nil
}

func (s *Shell) executeResetCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('reset <1|2>'): %+v", args)
	}

	s.Pod.Stop()
	networkcontroller.Disconnect()

	if args[1] == "0" {
		s.Pod = robot.NewPod(robot.NewExampleHexapod0())
	} else if args[1] == "1" {
		s.Pod = robot.NewPod(robot.NewExampleHexapod1())
	} else if args[1] == "2" {
		s.Pod = robot.NewPod(robot.NewExampleHexapod2())
	} else if args[1] == "3" {
		s.Pod = robot.NewPod(robot.NewExamplePentapod())
	} else if args[1] == "4" {
		s.Pod = robot.NewPod(robot.NewHeptapod())
	} else if args[1] == "5" {
		s.Pod = robot.NewPod(robot.NewSpider())
	} else {
		return fmt.Errorf("Unknown example preset")
	}

	s.Pod.Update()
	s.Pod.SetDebugChannel(s.outputCh)

	return nil
}

func (s *Shell) executeSpeedCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('speed <1-10>'): %+v", args)
	}

	speed, err := strconv.ParseFloat(args[1], 64)
	if err != nil || speed < 1 || speed > 10 {
		return fmt.Errorf("syntax error ('speed <1-10>'): %+v", args)
	}

	DELAY_COUNTER = 10 - int(speed)

	return nil
}

func (s *Shell) executeZLiftCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('zlift <0-90>'): %+v", args)
	}

	zlift, err := strconv.ParseFloat(args[1], 64)
	if err != nil || zlift < 0 || zlift > 90 {
		return fmt.Errorf("syntax error ('zlift <0-90>'): %+v", args)
	}

	robot.Z_LIFT = zlift

	return nil
}

func (s *Shell) executeGaitCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('gait <tripod|ripple|wave>'): %+v", args)
	}

	for {
		if s.Pod.IsReverting {
			time.Sleep(time.Millisecond * 20)
		} else {
			break
		}
	}

	var err error
	switch args[1] {
	case "tripod":
		s.Pod.BodyDefinition.Gait, err = robot.NewGait(s.Pod.BodyDefinition.NumLegs, robot.TRIPOD)
	case "ripple":
		s.Pod.BodyDefinition.Gait, err = robot.NewGait(s.Pod.BodyDefinition.NumLegs, robot.RIPPLE)
	case "wave":
		s.Pod.BodyDefinition.Gait, err = robot.NewGait(s.Pod.BodyDefinition.NumLegs, robot.WAVE)
	default:
		return fmt.Errorf("syntax error ('gait <tripod|ripple|wave>'): %+v", args)
	}

	s.Pod.ResetInterpolator()
	// s.Pod.RevertToNutral()
	s.Pod.SetDebugChannel(s.outputCh)
	if err != nil {
		return err
	}

	s.Pod.Start()
	networkcontroller.Start()

	return err
}

func (s *Shell) executeOpenServoPortCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('open_servo_port <IP:port>'): %+v", args)
	}

	if strings.Contains(args[1], ":") {
		s.outputCh <- fmt.Sprintf("Opening connection to %s", args[1])
		err := networkcontroller.Dial(args[1])
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("syntax error ('open_servo_port <IP:port>'): %+v", args)
	}

	return nil
}

func (s *Shell) executeCloseServoPortCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 1 {
		return fmt.Errorf("syntax error ('close_servo_port'): %+v", args)
	}

	networkcontroller.Disconnect()

	return nil
}

func (s *Shell) executeSaveCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	exists, err := folderExists(fmt.Sprintf("./%s", POD_FOLDER))
	if err != nil {
		return err
	}

	if !exists {
		err := os.Mkdir(fmt.Sprintf("%s", POD_FOLDER), 0755)
		if err != nil {
			return err
		}
	}

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('save <filename>'): %+v", args)
	}

	return s.Pod.BodyDefinition.Save(fmt.Sprintf("./%s/%s", POD_FOLDER, args[1]))
}

func (s *Shell) executeLoadCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('load <filename>'): %+v", args)
	}

	definition, err := s.Pod.BodyDefinition.Load(fmt.Sprintf("./%s/%s", POD_FOLDER, args[1]))
	if err != nil {
		return err
	}

	s.Pod.LoadBodyDefinition(definition)
	s.Pod.UpdatePodStructure()

	return nil
}

func (s *Shell) executeZeroCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)
	s.Pod.Zero()

	networkcontroller.Start()

	return nil
}

func (s *Shell) executeReverseCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	s.Pod.ReverseDirection()
	return nil
}

func (s *Shell) executeRevertCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	s.Pod.ResetInterpolator()
	s.Pod.RevertToNutral()

	// TODO: refactor
	if s.Pod.IsRecording {
		for _, l := range s.Pod.Legs {
			s.Pod.MotionPrimitive.Add(l.ServoAngles)
		}
	}

	return nil
}

func (s *Shell) executeRecordCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('record <filename>'): %+v", args)
	}

	if args[1] == "on" {
		s.Pod.IsRecording = true
	} else if args[1] == "off" {
		s.Pod.IsRecording = false
		s.Pod.ResetTicks()
		s.Pod.ClearPrimitives()
	} else {
		return fmt.Errorf("syntax error ('record <on|off>'): %+v", args)
	}
	return nil
}

func folderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *Shell) executeExportCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 4 {
		return fmt.Errorf("syntax error ('export <filename> <servo range> <mask>'): %+v", args)
	}

	servoRange, err := strconv.ParseInt(args[2], 10, 32)
	if err != nil {
		return fmt.Errorf("syntax error ('export <filename> <servo range>' <mask>): %+v", args)
	}

	exists, err := folderExists(fmt.Sprintf("./%s", PRIMITIVES_FOLDER))
	if err != nil {
		return err
	}

	if !exists {
		err := os.Mkdir(fmt.Sprintf("%s", PRIMITIVES_FOLDER), 0755)
		if err != nil {
			return err
		}
	}

	if len(args[3]) != 3 {
		return fmt.Errorf("invalid bit mask: %+v", args[3])
	}

	maskOk := true
	for _, bit := range args[3] {
		if bit != '0' && bit != '1' {
			maskOk = false
		}
	}
	if !maskOk {
		return fmt.Errorf("invalid bit mask: %+v", args[3])
	}

	if !s.Pod.IsRecording {
		return fmt.Errorf("no recording started. nothing to export")
	}

	s.Pod.MotionPrimitive.Export(fmt.Sprintf("./%s/%s", PRIMITIVES_FOLDER, args[1]), int(servoRange), args[3][0] == '1', args[3][1] == '1', args[3][2] == '1')

	s.outputCh <- fmt.Sprintf("Servoangles normalized for %d degrees", servoRange)
	s.outputCh <- "\tMidpoint equals raw value 512"
	s.outputCh <- fmt.Sprintf("\tnegative %d degrees equals raw value 0", servoRange/2)
	s.outputCh <- fmt.Sprintf("\tpositive %d degrees equals raw value 1024", servoRange/2)
	s.outputCh <- fmt.Sprintf("Recording exported to : ./%s/%s", PRIMITIVES_FOLDER, args[1])

	s.Pod.ClearPrimitives()

	return nil
}

func (s *Shell) executeDebugCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 1 {
		return fmt.Errorf("syntax error ('debug'): %+v", args)
	}

	s.Pod.Debug(fmt.Sprintf("Motion set size: %d", s.Pod.MotionPrimitive.Size()))

	return nil
}

func (s *Shell) executeStepCycleCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if !s.Pod.HasDefinedStride {
		return fmt.Errorf("no defined stride")
	}

	s.Pod.AddTargetGaitCycles(1)
	return nil
}

func (s *Shell) executePitchCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	return fmt.Errorf("pitch command is not implemented yet")
}

func (s *Shell) executeYawCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	return fmt.Errorf("yaw command is not implemented yet")
}

func (s *Shell) executeRollCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	return fmt.Errorf("roll command is not implemented yet")
}

func (s *Shell) executeUpCmd(args []string) error {

	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('up <z>'): %+v", args)
	}

	// z, err := strconv.ParseFloat(args[1], 64)
	// if err != nil {
	// 	return fmt.Errorf("syntax error ('stride_vector <nrepeats> <X> <Y>'): %+v", args)
	// }

	// s.Pod.ResetInterpolator()

	// return s.Pod.SetHeight(1, z)

	return fmt.Errorf("up command is not implemented yet")

}

func (s *Shell) executeDownCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('up <z>'): %+v", args)
	}

	z, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('stride_vector <nrepeats> <X> <Y>'): %+v", args)
	}

	for _, l := range s.Pod.Legs {
		l.Joints[0].Z = l.Joints[0].Z + z
		l.Joints[1].Z = l.Joints[1].Z + z
		l.NeutralEffectorCoordinate.Z = l.NeutralEffectorCoordinate.Z + z
	}

	s.Pod.Update()

	// TODO: This is a bit "out of band"... Refactor all recording stuff
	if s.Pod.IsRecording {
		for _, l := range s.Pod.Legs {
			s.Pod.MotionPrimitive.Add(l.ServoAngles)
		}
	}

	return fmt.Errorf("down command is not implemented yet")
}

func (s *Shell) executeGroundCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('ground <height>'): %+v", args)
	}

	height, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return fmt.Errorf("syntax error ('ground <height>'): %+v", args)
	}

	for _, l := range s.Pod.Legs {
		err := l.Ground(height)
		if err != nil {
			return err
		}
	}

	return nil
}

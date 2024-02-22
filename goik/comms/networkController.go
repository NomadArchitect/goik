package comms

import (
	"GOIK/robot"
	"net"
)

type ControlMode int

const (
	Streaming ControlMode = 1
	Primitive ControlMode = 2
)

type NetworkController struct {
	pod          *robot.Pod
	connection   *net.UDPConn
	packet       []byte
	id           uint8
	DebugChannel chan string
	isRunning    bool
	mode         ControlMode
}

func NewNetworkController(id uint8, p *robot.Pod, DebugChannel chan string) *NetworkController {
	return &NetworkController{id: id,
		pod:          p,
		packet:       make([]byte, 2*(robot.NUM_JOINTS-1)*p.BodyDefinition.NumLegs+3),
		DebugChannel: DebugChannel,
		isRunning:    false,
		mode:         Streaming,
	}
}

func (n *NetworkController) Update() {
	if n.connection == nil || !n.isRunning || n.mode != Streaming {
		return
	}

	// Frame format:
	// Byte
	// 0: Robot id
	// 1-2: Coxa 0
	// 3-4: Femur 0
	// 5-6: Tibia 0
	// 7-8: Coxa 1
	// ...
	// 50-51: XOR of id and all positions bytes

	// Note regarding servo orientation:
	// Formula for mapping joint angle to raw units for dynamixel servos with non 360 degree rotation:
	// 		raw_unit = uint16((n.pod.Legs[l].ServoAngles.Joint/<maxium servo rotation in degrees>)*<maximum unit> + <maximum unit/2>)
	// Example for XL-320:
	//	coxa = uint16((n.pod.Legs[l].ServoAngles.Coxa/300)*1024+512)
	//
	// In case a servo is mounted mirrored (coxa servo horn pointing down instead of up etc)
	// use the following formula instead
	//	coxa = 1024 - uint16((n.pod.Legs[l].ServoAngles.Coxa/300)*1024+512)

	var coxa uint16
	var femur uint16
	var tibia uint16
	var checksum uint16
	var i = 0

	n.packet[i] = n.id

	checksum ^= uint16(n.id)
	for l, _ := range n.pod.Legs {

		coxa = 1024 - uint16((n.pod.Legs[l].ServoAngles.Coxa/300)*1024+512)

		checksum ^= coxa
		i += 1
		n.packet[i] = uint8(coxa & 0xFF)
		i += 1
		n.packet[i] = uint8((coxa & 0xFF00) >> 8)

		femur = uint16((n.pod.Legs[l].ServoAngles.Femur/300)*1024 + 512)
		checksum ^= femur
		i += 1
		n.packet[i] = uint8(femur & 0xFF)
		i += 1
		n.packet[i] = uint8((femur & 0xFF00) >> 8)

		tibia = uint16((n.pod.Legs[l].ServoAngles.Tibia/300)*1024 + 512)
		i += 1
		n.packet[i] = uint8(tibia & 0xFF)
		i += 1
		n.packet[i] = uint8((tibia & 0xFF00) >> 8)
		checksum ^= tibia
	}

	i += 1
	n.packet[i] = uint8(checksum << 8)
	i += 1
	n.packet[i] = uint8(checksum >> 8)

	// n.DebugChannel <- fmt.Sprintf("Packet length: %d\n", len(n.packet))
	n.connection.Write(n.packet[:])
}

func (n *NetworkController) Dial(address string) error {
	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		return err
	}

	// Dial to the address with UDP
	n.connection, err = net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		return err
	}

	return nil
}

func (n *NetworkController) Disconnect() error {
	if n.connection == nil {
		return nil
	}

	return n.connection.Close()
}

func (n *NetworkController) Start() {
	n.isRunning = true
}



package serial

import (
	"fmt"
	"regexp"
	"strconv"

	"go.bug.st/serial"
)

type SerialPort struct {
	portName string
	baud     int
	port     serial.Port
	buff     []byte
}

func GetPorts() ([]string, error) {
	portNames, err := serial.GetPortsList()
	if err != nil {
		return []string{}, err
	}
	return portNames, nil
}

func BaudOptions() []string {
	return []string{
		"300",
		"600",
		"750",
		"1200",
		"2400",
		"4800",
		"9600",
		"19200",
		"31250",
		"38400",
		"57600",
		"74880",
		"115200",
		"230400",
		"250000",
		"460800",
		"500000",
		"921600",
		"1000000",
		"2000000",
	}
}

func (s *SerialPort) OpenPort() error {
	mode := &serial.Mode{
		BaudRate: s.baud,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	fmt.Println("Opening port...")
	fmt.Println("name: ", s.portName)
	fmt.Println("baud: ", s.baud)
	port, err := serial.Open(s.portName, mode)
	if err != nil {
		fmt.Println("Error opening port!!!")
		return err
	}
	s.port = port
	return nil
}

func (s *SerialPort) SetPortName(portName string) {
	s.portName = portName
}

func (s *SerialPort) SetBaud(baud int) {
	s.baud = baud
}

func (s *SerialPort) readPort(data []byte) (int, error) {
	n, err := s.port.Read(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func parseData(raw string) (float32, error) {
	expression := regexp.MustCompile(`\s*(?P<name>[^:]+):\s*(?P<value>[\d\.]+)`)
	match := expression.FindStringSubmatch(raw)
	result := make(map[string]string)
	for i, name := range expression.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	datum, err := strconv.ParseFloat(result["value"], 32)
	if err != nil {
		return 0, err
	}
	return float32(datum), nil

}

func New(portName string, baud int) *SerialPort {
	s := &SerialPort{
		portName: portName,
		baud:     baud,
		buff:     make([]byte, 255),
	}
	return s
}

func (s *SerialPort) ReadSource() (float32, error) {
	// TODO: Figure out how to stop whatever is going on here
	bytesRead, err := s.readPort(s.buff)
	if err != nil {
		fmt.Println("failed to read port", err)
		return 0, err
	}
	// fmt.Println("Buff: ", string(buff[:bytesRead]))
	datum, err := parseData(string(s.buff[:bytesRead]))
	if err != nil {
		fmt.Println("failed to parse data", err)
		return 0, err
	}
	fmt.Println("Datum: ", datum)
	return datum, err
}

func (s *SerialPort) Close() error {
	err := s.port.Close()
	if err != nil {
		fmt.Println("failed to close serial port", err)
	}
	return err
}

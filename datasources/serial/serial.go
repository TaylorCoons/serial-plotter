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

func (s *SerialPort) openPort() error {
	mode := &serial.Mode{
		BaudRate: s.baud,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(s.portName, mode)
	if err != nil {
		return err
	}
	s.port = port
	return nil
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

func New(portName string, baud int) (*SerialPort, error) {
	s := &SerialPort{
		portName: portName,
		baud:     baud,
		buff:     make([]byte, 255),
	}
	err := s.openPort()
	if err != nil {
		fmt.Println("Failed to open serial port: ", err)
		return nil, err
	}
	return s, nil
}

func (s *SerialPort) ReadSource() (float32, error) {
	// TODO: Figure out how to stop whatever is going on here
	bytesRead, err := s.readPort(s.buff)
	if err != nil {
		fmt.Println("failed to read port")
		fmt.Println("TODO: Figure out how to handle this error")
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

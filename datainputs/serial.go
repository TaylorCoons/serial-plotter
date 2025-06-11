package datainputs

import (
	"fmt"

	"go.bug.st/serial"
)

func GetPorts() ([]string, error) {
	portNames, err := serial.GetPortsList()
	if err != nil {
		return []string{}, err
	}
	return portNames, nil
}

func OpenPort(portName string, baud int) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: baud,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}
	return port, nil
}

func ReadPort(port serial.Port, data []byte) (int, error) {
	buff := make([]byte, 255) // TODO: Make configurable
	for {
		n, err := port.Read(buff)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			// EOF
			break
		}
		fmt.Printf("%s", string(buff[:n]))
	}
	return 0, nil
}

package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/taylorcoons/serial-plotter/gui"
)

func SetupLogger() (*log.Logger, error) {
	logDir := "./"
	logFile := path.Join(logDir, "serial-plotter.log")
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %v", err)
	}
	logger := log.New(f, "serial-plotter: ", log.Ldate|log.Lshortfile|log.Ltime)
	return logger, nil
}

func main() {

	// logger, err := SetupLogger()
	// if err != nil {
	// 	fmt.Fprint(os.Stderr, err)
	// }
	// ports, err := datainputs.GetPorts()
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// for _, port := range ports {
	// 	fmt.Printf("Found port: %v\n", port)
	// }
	// port, err := datainputs.OpenPort("/dev/ttyUSB0", 9600)
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// buff := make([]byte, 255)
	// read, err := datainputs.ReadPort(port, buff)
	// logger.Print(read)
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	gui.Main()

}

package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/taylorcoons/serial-plotter/datasources/serial"
	"github.com/taylorcoons/serial-plotter/gui/graph"
)

func Main() {
	serialSource, err := serial.New("/dev/ttyUSB0", 9600)
	if err != nil {
		fmt.Println("Failed to create serialSource", err)
	}

	dataChannel := make(chan float32)

	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")

	myWindow.Resize(fyne.NewSize(800, 800))
	ports, err := serial.GetPorts()
	if err != nil {
		fmt.Println("failed to get ports", err)
	}
	portSelect := widget.NewSelect(ports, func(value string) {
		fmt.Println("Port set to ", value)
	})
	portSelect.PlaceHolder = "Serial Port"

	baudSelect := widget.NewSelect([]string{"4800", "9600"}, func(value string) {
		fmt.Println("Baud set to ", value)
	})
	stop := make(chan int)
	stopButton := widget.NewButton("Stop", func() {
		stop <- 0
	})
	startButton := widget.NewButton("Start", func() {
		// tick := time.NewTicker(time.Millisecond * 250)
		// go func() {
		// 	counter := 0
		// 	for {
		// 		select {
		// 		case <-tick.C:
		// 			fmt.Println("Tick")
		// 			counter++
		// 			value := 10 * float32(math.Sin(float64(counter%100)*math.Pi/100*2))
		// 			dataChannel <- float32(value)
		// 		case <-stop:
		// 			fmt.Println("Stopping")
		// 			return
		// 		}
		// 	}
		// }()
		go func() {
			for {
				datum, err := serialSource.ReadSource()
				select {
				case <-stop:
					fmt.Println("Stopping data collection")
					return
				default:
					if err != nil {
						fmt.Println("failed to read source", err)
						return
					}
					dataChannel <- datum
				}
			}
		}()
		fmt.Println("Start pressed")
	})
	baudSelect.PlaceHolder = "Baud Rate"
	serialOptions := container.NewVBox(portSelect, baudSelect)
	graphControls := container.NewVBox(startButton, stopButton)
	options := container.NewHBox(serialOptions, graphControls)
	graphContainer := container.NewWithoutLayout()
	content := container.NewBorder(options, nil, nil, nil, graphContainer)

	myWindow.SetContent(content)
	data := []float32{}
	graphStruct := graph.GraphStruct{}
	graphStruct.Show(graphContainer)
	go func() {
		for {
			value, ok := <-dataChannel
			if ok {
				fmt.Println("Appending data")
				data = append(data, value)
				graphStruct.Update(graphContainer, data)
				fyne.Do(func() {
					graphContainer.Refresh()
				})
			}
		}
	}()
	myWindow.ShowAndRun()

}

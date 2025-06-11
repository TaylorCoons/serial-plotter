package gui

import (
	"fmt"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/taylorcoons/serial-plotter/datainputs"
	"github.com/taylorcoons/serial-plotter/gui/graph"
)

func Main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")

	myWindow.Resize(fyne.NewSize(800, 800))
	ports, err := datainputs.GetPorts()
	if err != nil {
		fmt.Println("Couldn't get ports")
	}

	portSelect := widget.NewSelect(ports, func(value string) {
		fmt.Println("Port set to ", value)
	})
	portSelect.PlaceHolder = "Serial Port"

	baudSelect := widget.NewSelect([]string{"4800", "9600"}, func(value string) {
		fmt.Println("Baud set to ", value)
	})
	dataChannel := make(chan float32)
	stop := make(chan int)
	stopButton := widget.NewButton("Stop", func() {
		stop <- 0
	})
	startButton := widget.NewButton("Start", func() {
		tick := time.NewTicker(time.Millisecond * 250)
		go func() {
			counter := 0
			for {
				select {
				case <-tick.C:
					fmt.Println("Tick")
					counter++
					value := 10 * float32(math.Sin(float64(counter%100)*math.Pi/100*2))
					dataChannel <- float32(value)
				case <-stop:
					fmt.Println("Stopping")
					return
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

package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/taylorcoons/serial-plotter/datasources"
	"github.com/taylorcoons/serial-plotter/datasources/pseudo"
	"github.com/taylorcoons/serial-plotter/datasources/serial"
	"github.com/taylorcoons/serial-plotter/gui/graph"
)

func Main() {
	serialSource, err := serial.New("/dev/ttyUSB0", 9600)
	if err != nil {
		fmt.Println("Failed to create serialSource", err)
	}

	transformMap := map[string]pseudo.Transform{
		"Sine":     pseudo.SinTransform,
		"Square":   pseudo.SquareTransform,
		"Sawtooth": pseudo.SawtoothTransform,
	}
	pseudoSource := pseudo.New(time.Millisecond*250, pseudo.SawtoothTransform)

	dataChannel := make(chan float32)

	myApp := app.New()
	myWindow := myApp.NewWindow("Serial Plotter")

	myWindow.Resize(fyne.NewSize(800, 800))
	ports, err := serial.GetPorts()
	if err != nil {
		fmt.Println("failed to get ports", err)
	}
	dataSourcesList := []string{"Serial", "Dummy"}
	dataSourcesSelect := widget.NewSelect(dataSourcesList, func(value string) {
		fmt.Println("data sources set to: ", value)
	})
	dataSourcesSelect.SetSelectedIndex(0)
	portSelect := widget.NewSelect(ports, func(value string) {
		fmt.Println("port set to ", value)
	})
	portSelect.PlaceHolder = "Serial Port"

	baudSelect := widget.NewSelect([]string{"4800", "9600"}, func(value string) {
		fmt.Println("baud set to ", value)
	})
	stop := make(chan int)
	stopButton := widget.NewButton("Stop", func() {
		stop <- 0
	})
	startButton := widget.NewButton("Start", func() {
		go func() {
			// TODO: Select off UI
			var dataSource datasources.DataSourcer
			// For now to change the source put the one you want to use last
			dataSource = serialSource
			dataSource = pseudoSource
			for {
				datum, err := dataSource.ReadSource()
				select {
				case <-stop:
					fmt.Println("stopping data collection")
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
	transformKeys := []string{}
	for k := range transformMap {
		transformKeys = append(transformKeys, k)
	}
	transformSelect := widget.NewSelect(transformKeys, func(value string) {
		fmt.Println("changed transform to: ", value)
	})
	transformSelect.SetSelectedIndex(0)
	dummyOptions := container.NewVBox(transformSelect)
	graphControls := container.NewVBox(startButton, stopButton)
	options := container.NewGridWithColumns(4, dataSourcesSelect, serialOptions, dummyOptions, graphControls)
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

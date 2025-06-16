package gui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/taylorcoons/serial-plotter/datasources"
	"github.com/taylorcoons/serial-plotter/datasources/pseudo"
	"github.com/taylorcoons/serial-plotter/datasources/serial"
	"github.com/taylorcoons/serial-plotter/gui/graph"
	"github.com/taylorcoons/serial-plotter/transformers"
	"github.com/taylorcoons/serial-plotter/transformers/passthrough"
	"github.com/taylorcoons/serial-plotter/transformers/sma"
)

type appState struct {
	dataSourceType string
	serialSource   *serial.SerialPort
	dummySource    *pseudo.Pseudo
}

func (a *appState) DataSourcesPanel(serialSourceContainer *fyne.Container, dummySourceContainer *fyne.Container) *fyne.Container {
	dataSourcesList := []string{"Serial", "Dummy"}
	dataSourcesSelect := widget.NewSelect(dataSourcesList, func(value string) {
		serialSourceContainer.Hide()
		dummySourceContainer.Hide()
		switch value {
		case "Serial":
			serialSourceContainer.Show()
		case "Dummy":
			dummySourceContainer.Show()
		}
		a.dataSourceType = value
		fmt.Println("datasourcetype set to: ", value)
	})
	dataSourcesDefaultIndex := 0
	dataSourcesSelect.SetSelectedIndex(dataSourcesDefaultIndex)
	a.dataSourceType = dataSourcesList[dataSourcesDefaultIndex]
	dataSourcesContainer := container.NewVBox(dataSourcesSelect)
	return dataSourcesContainer
}

func (a *appState) SerialSourceOptions() (*fyne.Container, error) {
	ports, err := serial.GetPorts()
	if err != nil {
		fmt.Println("failed to get ports", err)
	}
	portName := ""
	if len(ports) > 0 {
		portName = ports[0]
	}
	defaultBaudIndex := 1
	baudOptions := []string{"4800", "9600"}
	defaultBaudValue, err := strconv.Atoi(baudOptions[defaultBaudIndex])
	a.serialSource = serial.New(portName, defaultBaudValue)
	portSelect := widget.NewSelect(ports, func(value string) {
		a.serialSource.SetPortName(value)
		fmt.Println("port set to ", value)
	})
	portSelect.PlaceHolder = "Serial Port"
	if err != nil {
		fmt.Println("failed to parse baud option", err)
		return nil, err
	}
	baudSelect := widget.NewSelect(baudOptions, func(value string) {
		baudValue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("failed to parse baud option", err)
			return
		}
		a.serialSource.SetBaud(baudValue)
		fmt.Println("baud set to ", value)
	})
	baudSelect.SetSelectedIndex(defaultBaudIndex)
	serialOptions := container.NewVBox(portSelect, baudSelect)
	return serialOptions, nil
}

func (a *appState) DummySourceOptions() *fyne.Container {
	transformMap := map[string]pseudo.Transform{
		"Sine":     pseudo.SinTransform,
		"Square":   pseudo.SquareTransform,
		"Sawtooth": pseudo.SawtoothTransform,
	}
	defaultTransformIndex := 0
	transformKeys := []string{}
	for k := range transformMap {
		transformKeys = append(transformKeys, k)
	}
	defaultTransform := transformMap[transformKeys[defaultTransformIndex]]
	a.dummySource = pseudo.New(time.Millisecond*250, defaultTransform)
	transformSelect := widget.NewSelect(transformKeys, func(value string) {
		a.dummySource.SetTransform(transformMap[value])
		fmt.Println("changed transform to: ", value)
	})
	transformSelect.SetSelectedIndex(defaultTransformIndex)
	return container.NewVBox(transformSelect)
}

func (a *appState) controlsPanel(dataChannel chan float32) *fyne.Container {
	stop := make(chan int)
	stopButton := widget.NewButton("Stop", func() {
		fmt.Println("Stop pressed")
		stop <- 0
	})
	startButton := widget.NewButton("Start", func() {
		go func() {
			var dataSource datasources.DataSourcer
			switch a.dataSourceType {
			case "Dummy":
				dataSource = a.dummySource
			case "Serial":
				a.serialSource.OpenPort()
				dataSource = a.serialSource
			}
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
	return container.NewVBox(startButton, stopButton)
}

func Main() {
	dataChannel := make(chan float32)

	appState := &appState{}

	myApp := app.New()
	myWindow := myApp.NewWindow("Serial Plotter")

	myWindow.Resize(fyne.NewSize(800, 800))

	serialOptions, err := appState.SerialSourceOptions()
	if err != nil {
		fmt.Println("failed to create serial source options")
	}
	dummyOptions := appState.DummySourceOptions()
	controlsPanel := appState.controlsPanel(dataChannel)
	dataSourcesPanel := appState.DataSourcesPanel(serialOptions, dummyOptions)
	options := container.NewGridWithColumns(3, dataSourcesPanel, serialOptions, dummyOptions, controlsPanel)
	graphContainer := container.NewWithoutLayout()
	content := container.NewBorder(options, nil, nil, nil, graphContainer)

	myWindow.SetContent(content)
	data := []float32{}
	graphStruct := graph.GraphStruct{}
	graphStruct.Show(graphContainer)
	go func() {
		sma := sma.New(3)
		passthrough := passthrough.New()
		// TODO: Select off of UI
		var transformer transformers.Transformer
		transformer = sma
		transformer = passthrough
		for {
			value, ok := <-dataChannel
			if ok {
				fmt.Println("Appending data")
				data = append(data, transformer.Compute(data, value))
				graphStruct.Update(graphContainer, data)
				fyne.Do(func() {
					graphContainer.Refresh()
				})
			}
		}
	}()
	myWindow.ShowAndRun()
}

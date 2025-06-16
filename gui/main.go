package gui

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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
	transform      transformers.Transformer
	window         fyne.Window
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
	functionMap := map[string]pseudo.Function{
		"Sine":     pseudo.SinFunction,
		"Square":   pseudo.SquareFunction,
		"Sawtooth": pseudo.SawtoothFunction,
	}
	defaultFunctionIndex := 0
	functionKeys := []string{}
	for k := range functionMap {
		functionKeys = append(functionKeys, k)
	}
	defaultFunction := functionMap[functionKeys[defaultFunctionIndex]]
	a.dummySource = pseudo.New(time.Millisecond*250, defaultFunction)
	functionSelect := widget.NewSelect(functionKeys, func(value string) {
		a.dummySource.SetFunction(functionMap[value])
		fmt.Println("changed function to: ", value)
	})
	functionSelect.SetSelectedIndex(defaultFunctionIndex)
	return container.NewVBox(functionSelect)
}

func (a *appState) TransformOptions() *fyne.Container {
	transformMap := map[string]transformers.Transformer{
		"None":                  passthrough.New(),
		"Simple Moving Average": sma.New(3),
	}
	defaultTransformIndex := 0
	transformKeys := []string{}
	for k := range transformMap {
		transformKeys = append(transformKeys, k)
	}
	transformSelect := widget.NewSelect(transformKeys, func(value string) {
		a.transform = transformMap[value]
	})
	transformSelect.SetSelectedIndex(defaultTransformIndex)
	return container.NewVBox(transformSelect)

}

func ErrorModal(message string, window fyne.Window) {
	text := canvas.NewText(message, color.Black)
	var popUp *widget.PopUp
	closeButton := widget.NewButton("Close", func() {
		popUp.Hide()
	})
	container := container.NewVBox(text, closeButton)
	popUp = widget.NewModalPopUp(container, window.Canvas())
	fyne.Do(func() {
		popUp.Show()
	})
}

func (a *appState) ControlsPanel(dataChannel chan float32, window fyne.Window) *fyne.Container {
	stop := make(chan int)
	stopButton := widget.NewButton("Stop", func() {
		fmt.Println("Stop pressed")
		switch a.dataSourceType {
		case "Serial":
			err := a.serialSource.Close()
			if err != nil {
				ErrorModal(fmt.Sprintf("Error closing port %s", err), a.window)
			}
		}
		stop <- 0
	})
	startButton := widget.NewButton("Start", func() {
		go func() {
			var dataSource datasources.DataSourcer
			switch a.dataSourceType {
			case "Dummy":
				dataSource = a.dummySource
			case "Serial":
				fmt.Println("Opening serial port")
				err := a.serialSource.OpenPort()
				if err != nil {
					fmt.Println("error opening port ", err)
					ErrorModal(fmt.Sprintf("Error opening port %s", err), a.window)
					return
				}
				dataSource = a.serialSource
			}
			for {
				datum, err := dataSource.ReadSource()
				if err != nil {
					fmt.Println("failed to read source", err)
					ErrorModal(fmt.Sprintf("Failed to read data source %s", err), a.window)
					return
				}
				select {
				case <-stop:
					fmt.Println("stopping data collection")
					return
				default:
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
	appState.window = myWindow

	myWindow.Resize(fyne.NewSize(800, 800))

	serialOptions, err := appState.SerialSourceOptions()
	if err != nil {
		fmt.Println("failed to create serial source options")
	}
	dummyOptions := appState.DummySourceOptions()
	controlsPanel := appState.ControlsPanel(dataChannel, myWindow)
	dataSourcesPanel := appState.DataSourcesPanel(serialOptions, dummyOptions)
	transformOptions := appState.TransformOptions()
	options := container.NewGridWithColumns(4, dataSourcesPanel, serialOptions, dummyOptions, transformOptions, controlsPanel)
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
				data = append(data, appState.transform.Compute(data, value))
				graphStruct.Update(graphContainer, data)
				fyne.Do(func() {
					graphContainer.Refresh()
				})
			}
		}
	}()
	myWindow.ShowAndRun()
}

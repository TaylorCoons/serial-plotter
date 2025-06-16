package gui

import (
	"fmt"
	"image/color"
	"slices"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/taylorcoons/serial-plotter/datasources"
	"github.com/taylorcoons/serial-plotter/datasources/dummy"
	"github.com/taylorcoons/serial-plotter/datasources/serial"
	"github.com/taylorcoons/serial-plotter/gui/graph"
	"github.com/taylorcoons/serial-plotter/gui/preference"
	"github.com/taylorcoons/serial-plotter/transformers"
	"github.com/taylorcoons/serial-plotter/transformers/passthrough"
	"github.com/taylorcoons/serial-plotter/transformers/sma"
)

type appState struct {
	dataSourceType string
	serialSource   *serial.SerialPort
	dummySource    *dummy.Dummy
	transform      transformers.Transformer
	window         fyne.Window
	data           []float32
	app            fyne.App
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
		a.app.Preferences().SetString(preference.DataSource.String(), value)
		a.dataSourceType = value
		fmt.Println("datasourcetype set to: ", value)
	})
	selected := a.app.Preferences().StringWithFallback(preference.DataSource.String(), "Dummy")
	dataSourcesSelect.SetSelected(selected)
	a.dataSourceType = selected
	dataSourcesContainer := container.NewVBox(dataSourcesSelect)
	return dataSourcesContainer
}

func (a *appState) SerialSourceOptions() (*fyne.Container, error) {
	ports, err := serial.GetPorts()
	if err != nil {
		fmt.Println("failed to get ports", err)
	}
	defaultPort := a.app.Preferences().StringWithFallback(preference.PortName.String(), "")
	if !slices.Contains(ports, defaultPort) {
		defaultPort = ""
	}
	baudOptions := serial.BaudOptions()
	defaultBaud := a.app.Preferences().StringWithFallback(preference.Baud.String(), "9600")
	defaultBaudValue, err := strconv.Atoi(defaultBaud)
	if err != nil {
		fmt.Println("failed to parse baud option", err)
		return nil, err
	}
	a.serialSource = serial.New(defaultPort, defaultBaudValue)
	portSelect := widget.NewSelect(ports, func(value string) {
		a.serialSource.SetPortName(value)
		a.app.Preferences().SetString(preference.PortName.String(), value)
		fmt.Println("port set to ", value)
	})
	if defaultPort != "" {
		portSelect.SetSelected(defaultPort)
	} else {
		portSelect.PlaceHolder = "Serial Port"
	}
	baudSelect := widget.NewSelect(baudOptions, func(value string) {
		baudValue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("failed to parse baud option", err)
			return
		}
		a.serialSource.SetBaud(baudValue)
		a.app.Preferences().SetString(preference.Baud.String(), value)
		fmt.Println("baud set to ", value)
	})
	baudSelect.SetSelected(defaultBaud)
	serialOptions := container.NewVBox(portSelect, baudSelect)
	return serialOptions, nil
}

func (a *appState) DummySourceOptions() *fyne.Container {
	functionMap := map[string]dummy.Function{
		"Sine":     dummy.SinFunction,
		"Square":   dummy.SquareFunction,
		"Sawtooth": dummy.SawtoothFunction,
		"Constant": dummy.ConstantFunction,
	}
	functionKeys := []string{}
	for k := range functionMap {
		functionKeys = append(functionKeys, k)
	}
	selectedFunction := a.app.Preferences().StringWithFallback(preference.Function.String(), "Sine")
	a.dummySource = dummy.New(time.Millisecond*250, functionMap[selectedFunction])
	functionSelect := widget.NewSelect(functionKeys, func(value string) {
		a.dummySource.SetFunction(functionMap[value])
		a.app.Preferences().SetString(preference.Function.String(), value)
		fmt.Println("changed function to: ", value)
	})
	functionSelect.SetSelected(selectedFunction)
	return container.NewVBox(functionSelect)
}

func (a *appState) TransformOptions() *fyne.Container {
	transformMap := map[string]transformers.Transformer{
		"None":                  passthrough.New(),
		"Simple Moving Average": sma.New(3),
	}
	transformKeys := []string{}
	for k := range transformMap {
		transformKeys = append(transformKeys, k)
	}
	selected := a.app.Preferences().StringWithFallback(preference.Transform.String(), "None")
	transformSelect := widget.NewSelect(transformKeys, func(value string) {
		a.transform = transformMap[value]
		a.app.Preferences().SetString(preference.Transform.String(), value)
	})
	transformSelect.SetSelected(selected)
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

func (a *appState) InitializeSource() (datasources.DataSourcer, error) {
	switch a.dataSourceType {
	case "Dummy":
		return a.dummySource, nil
	case "Serial":
		fmt.Println("Opening serial port")
		err := a.serialSource.OpenPort()
		if err != nil {
			fmt.Println("error opening port ", err)
			ErrorModal(fmt.Sprintf("Error opening port %s", err), a.window)
			return nil, err
		}
		return a.serialSource, nil
	}
	return nil, fmt.Errorf("unknown data source selected")
}

func (a *appState) CloseDataSource() error {
	switch a.dataSourceType {
	case "Serial":
		err := a.serialSource.Close()
		if err != nil {
			ErrorModal(fmt.Sprintf("Error closing port %s", err), a.window)
			return err
		}
	}
	return nil
}

func (a *appState) ControlsPanel(dataChannel chan float32, clearChannel chan int, window fyne.Window) *fyne.Container {
	stop := make(chan int)
	var startButtonContainer *fyne.Container
	var stopButtonContainer *fyne.Container
	stopButton := widget.NewButton("Stop", func() {
		fmt.Println("Stop pressed")
		stop <- 0
	})
	startButton := widget.NewButton("Start", func() {
		startButtonContainer.Hide()
		stopButtonContainer.Show()
		go func() {
			defer fyne.Do(func() {
				startButtonContainer.Show()
				stopButtonContainer.Hide()
			})
			dataSource, err := a.InitializeSource()
			if err != nil {
				fmt.Println("Failed to initialize data source", err)
				return
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
					err := a.CloseDataSource()
					if err != nil {
						fmt.Println("failed to close data source", err)
					}
					return
				default:
					dataChannel <- datum
				}
			}
		}()
		fmt.Println("Start pressed")
	})
	clearButton := widget.NewButton("Clear", func() {
		clearChannel <- 0
	})
	startButton.Importance = widget.LowImportance
	stopButton.Importance = widget.LowImportance
	startButtonContainer = container.NewStack(canvas.NewRectangle(color.RGBA{0, 255, 0, 127}), startButton)
	stopButtonContainer = container.NewStack(canvas.NewRectangle(color.RGBA{255, 0, 0, 127}), stopButton)
	stopButtonContainer.Hide()

	return container.NewVBox(startButtonContainer, stopButtonContainer, clearButton)
}

func Main() {
	dataChannel := make(chan float32)
	clearChannel := make(chan int)

	app := app.New()
	appState := &appState{app: app}
	window := app.NewWindow("Serial Plotter")
	appState.window = window

	window.Resize(fyne.NewSize(800, 800))

	serialOptions, err := appState.SerialSourceOptions()
	if err != nil {
		fmt.Println("failed to create serial source options")
	}
	dummyOptions := appState.DummySourceOptions()
	controlsPanel := appState.ControlsPanel(dataChannel, clearChannel, window)
	dataSourcesPanel := appState.DataSourcesPanel(serialOptions, dummyOptions)
	transformOptions := appState.TransformOptions()
	options := container.NewGridWithColumns(4, dataSourcesPanel, serialOptions, dummyOptions, transformOptions, controlsPanel)
	graphContainer := container.NewWithoutLayout()
	content := container.NewBorder(options, nil, nil, nil, graphContainer)

	window.SetContent(content)
	appState.data = []float32{}
	graphStruct := graph.GraphStruct{}
	graphStruct.Show(graphContainer)
	go func() {
		for {
			select {
			case value := <-dataChannel:
				fmt.Println("Appending data")
				appState.data = append(appState.data, appState.transform.Compute(appState.data, value))
			case <-clearChannel:
				fmt.Println("Clearing data")
				appState.data = []float32{}
			}
			graphStruct.Update(graphContainer, appState.data)
			fyne.Do(func() {
				graphContainer.Refresh()
			})
		}
	}()
	window.ShowAndRun()
}

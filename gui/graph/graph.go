package graph

import (
	"image/color"
	"math"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type GraphStruct struct {
	xAxis, yAxis   *canvas.Line
	xTicks, yTicks []*canvas.Line
	lines          []*canvas.Line
}

type axisRange struct {
	min         float32
	max         float32
	realizedMin float32
	realizedMax float32
	tickSize    float32
	numTicks    int
}

func (g *GraphStruct) render(graphContainer *fyne.Container, size fyne.Size, data []float32) {
	g.xAxis = &canvas.Line{}
	g.yAxis = &canvas.Line{}
	g.xAxis.StrokeWidth = 2
	g.xAxis.StrokeColor = color.Black
	g.yAxis.StrokeWidth = 2
	g.yAxis.StrokeColor = color.Black

	yMin := float32(-10)
	yMax := float32(10)
	if len(data) != 0 {
		yMin = slices.Min(data)
		yMax = slices.Max(data)
	}
	yMagnitude := math.Abs(float64(yMax - yMin))
	orderMagnitude := 1
	for yMagnitude/10 > 1 {
		yMagnitude = yMagnitude / 10
		orderMagnitude = orderMagnitude * 10
	}
	realizedMin := yMin - float32(orderMagnitude)
	realizedMax := yMax + float32(orderMagnitude)
	yRange := axisRange{
		min:         yMin,
		max:         yMax,
		realizedMin: realizedMin,
		realizedMax: realizedMax,
		tickSize:    float32(orderMagnitude),
		numTicks:    int(math.Abs(float64(realizedMax-realizedMin)) / float64(orderMagnitude)),
	}

	zeroHeight := linearMap(0, yRange.realizedMin, yRange.realizedMax, size.Height, 0)
	g.xAxis.Position1 = fyne.NewPos(0, zeroHeight)
	g.xAxis.Position2 = fyne.NewPos(size.Width, zeroHeight)
	g.yAxis.Position1 = fyne.NewPos(0, 0)
	g.yAxis.Position2 = fyne.NewPos(0, size.Height)

	g.xTicks = []*canvas.Line{}
	for index := range data {
		xTick := &canvas.Line{}
		// TODO: Make tick length relative
		xTick.Position1 = fyne.NewPos(float32(index)*size.Width/float32(len(data)), zeroHeight+5)
		xTick.Position2 = fyne.NewPos(float32(index)*size.Width/float32(len(data)), zeroHeight-5)
		xTick.StrokeColor = color.Black
		xTick.StrokeWidth = 2
		g.xTicks = append(g.xTicks, xTick)
	}
	g.yTicks = []*canvas.Line{}
	for index := 0; index < yRange.numTicks; index++ {
		yTick := &canvas.Line{}
		// TODO: Make this tick length relative
		yTick.Position1 = fyne.NewPos(0, linearMap(float32(index), 0, float32(yRange.numTicks), size.Height, 0))
		yTick.Position2 = fyne.NewPos(5, linearMap(float32(index), 0, float32(yRange.numTicks), size.Height, 0))
		yTick.StrokeColor = color.Black
		yTick.StrokeWidth = 2
		g.yTicks = append(g.yTicks, yTick)
	}
	g.lines = []*canvas.Line{}
	for index := range data {
		if index == 0 {
			continue
		}
		line := &canvas.Line{}
		line.Position1 = fyne.NewPos(linearMap(float32(index-1), 0, float32(len(data)), 0, size.Width), linearMap(data[index-1], yRange.realizedMin, yRange.realizedMax, size.Height, 0))
		line.Position2 = fyne.NewPos(linearMap(float32(index), 0, float32(len(data)), 0, size.Width), linearMap(data[index], yRange.realizedMin, yRange.realizedMax, size.Height, 0))
		line.StrokeColor = color.Black
		line.StrokeWidth = 1
		g.lines = append(g.lines, line)
	}

	graphContainer.RemoveAll()
	// Add graph objects
	for _, xTick := range g.xTicks {
		graphContainer.Objects = append(graphContainer.Objects, xTick)
	}
	for _, yTick := range g.yTicks {
		graphContainer.Objects = append(graphContainer.Objects, yTick)
	}
	for _, line := range g.lines {
		graphContainer.Objects = append(graphContainer.Objects, line)
	}
	graphContainer.Objects = append(graphContainer.Objects, g.xAxis, g.yAxis)
}

func (g *GraphStruct) Show(graphContainer *fyne.Container) {
	g.render(graphContainer, graphContainer.Size(), []float32{})
}

func (g *GraphStruct) Update(graphContainer *fyne.Container, data []float32) {
	g.render(graphContainer, graphContainer.Size(), data)
}

func linearMap(value, inputMin, inputMax, outputMin, outputMax float32) float32 {
	return outputMin + (outputMax-outputMin)/(inputMax-inputMin)*(value-inputMin)
}

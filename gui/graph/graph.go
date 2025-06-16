package graph

import (
	"image/color"
	"math"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type GraphStruct struct {
	xAxis, yAxis     *canvas.Line
	xTicks, yTicks   []*canvas.Line
	xLabels, yLabels []*canvas.Text
	lines            []*canvas.Line
}

type axisRange struct {
	min         float32
	max         float32
	realizedMin float32
	realizedMax float32
	tickSize    float32
	numTicks    int
	zeroHeight  float32
}

func (g *GraphStruct) render(graphContainer *fyne.Container, size fyne.Size, data []float32) {
	g.xAxis = &canvas.Line{}
	g.yAxis = &canvas.Line{}
	g.xAxis.StrokeWidth = 2
	g.xAxis.StrokeColor = color.White
	g.yAxis.StrokeWidth = 2
	g.yAxis.StrokeColor = color.White

	yRange := g.createYRange(&size, data)

	g.addAxes(&size, &yRange)

	g.addXTicks(&size, &yRange, data)

	g.addYTicks(&size, &yRange, data)

	g.addLines(&size, &yRange, data)

	g.addGraphObjects(graphContainer)
}

func (g *GraphStruct) createYRange(size *fyne.Size, data []float32) axisRange {
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
	zeroHeight := linearMap(0, realizedMin, realizedMax, size.Height, 0)
	return axisRange{
		min:         yMin,
		max:         yMax,
		realizedMin: realizedMin,
		realizedMax: realizedMax,
		tickSize:    float32(orderMagnitude),
		numTicks:    int(math.Abs(float64(realizedMax-realizedMin)) / float64(orderMagnitude)),
		zeroHeight:  zeroHeight,
	}
}

func (g *GraphStruct) addAxes(size *fyne.Size, yRange *axisRange) {
	g.xAxis.Position1 = fyne.NewPos(0, yRange.zeroHeight)
	g.xAxis.Position2 = fyne.NewPos(size.Width, yRange.zeroHeight)
	g.yAxis.Position1 = fyne.NewPos(0, 0)
	g.yAxis.Position2 = fyne.NewPos(0, size.Height)
}

func positionXLabel(index int, length int, xLabel *canvas.Text, size *fyne.Size, yRange *axisRange) fyne.Position {
	xPos := float32(index)*size.Width/float32(length) + xLabel.Size().Width/2
	yPos := yRange.zeroHeight + 5
	return fyne.NewPos(xPos, yPos)
}

func (g *GraphStruct) addXTicks(size *fyne.Size, yRange *axisRange, data []float32) {
	g.xTicks = []*canvas.Line{}
	g.xLabels = []*canvas.Text{}
	xLabelDensity := func(data []float32, xLabelDivisor int, size *fyne.Size) float32 {
		return float32(len(data)) / float32(xLabelDivisor) / size.Width
	}
	xLabelDivisor := 1
	for xLabelDensity(data, xLabelDivisor, size) > 0.02 {
		xLabelDivisor++
	}
	for index := 0; index < len(data); index += xLabelDivisor {
		xTick := &canvas.Line{}
		xLabel := canvas.NewText(strconv.Itoa(index), color.White)
		xLabel.Alignment = fyne.TextAlignCenter
		xLabel.Move(positionXLabel(index, len(data), xLabel, size, yRange))
		// TODO: Make tick length relative
		xTick.Position1 = fyne.NewPos(float32(index)*size.Width/float32(len(data)), yRange.zeroHeight+5)
		xTick.Position2 = fyne.NewPos(float32(index)*size.Width/float32(len(data)), yRange.zeroHeight-5)
		xTick.StrokeColor = color.White
		xTick.StrokeWidth = 2
		g.xTicks = append(g.xTicks, xTick)
		g.xLabels = append(g.xLabels, xLabel)
	}
}

func (g *GraphStruct) addYTicks(size *fyne.Size, yRange *axisRange, data []float32) {
	g.yTicks = []*canvas.Line{}
	for index := 0; index < yRange.numTicks; index++ {
		yTick := &canvas.Line{}
		// TODO: Make this tick length relative
		yTick.Position1 = fyne.NewPos(0, linearMap(float32(index), 0, float32(yRange.numTicks), size.Height, 0))
		yTick.Position2 = fyne.NewPos(5, linearMap(float32(index), 0, float32(yRange.numTicks), size.Height, 0))
		yTick.StrokeColor = color.White
		yTick.StrokeWidth = 2
		g.yTicks = append(g.yTicks, yTick)
	}
}

func (g *GraphStruct) addLines(size *fyne.Size, yRange *axisRange, data []float32) {
	g.lines = []*canvas.Line{}
	for index := range data {
		if index == 0 {
			continue
		}
		line := &canvas.Line{}
		line.Position1 = fyne.NewPos(linearMap(float32(index-1), 0, float32(len(data)), 0, size.Width), linearMap(data[index-1], yRange.realizedMin, yRange.realizedMax, size.Height, 0))
		line.Position2 = fyne.NewPos(linearMap(float32(index), 0, float32(len(data)), 0, size.Width), linearMap(data[index], yRange.realizedMin, yRange.realizedMax, size.Height, 0))
		line.StrokeColor = color.White
		line.StrokeWidth = 1
		g.lines = append(g.lines, line)
	}
}

func (g *GraphStruct) addGraphObjects(graphContainer *fyne.Container) {
	graphContainer.RemoveAll()
	for _, xTick := range g.xTicks {
		graphContainer.Objects = append(graphContainer.Objects, xTick)
	}
	for _, xLabel := range g.xLabels {
		graphContainer.Objects = append(graphContainer.Objects, xLabel)
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

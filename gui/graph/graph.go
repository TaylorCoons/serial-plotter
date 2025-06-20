package graph

import (
	"image/color"
	"math"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
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
	yAxisOffset float32
	tickMin     float32
	tickMax     float32
	tickLength  float32
}

func foregroundColor() color.Color {
	return theme.DefaultTheme().Color(theme.ColorNameForeground, theme.VariantDark)
}

func backgroundColor() color.Color {
	return theme.DefaultTheme().Color(theme.ColorNameBackground, theme.VariantDark)
}

func primaryColor() color.Color {
	return theme.DefaultTheme().Color(theme.ColorNamePrimary, theme.VariantDark)
}

func (g *GraphStruct) render(graphContainer *fyne.Container, size fyne.Size, data []float32) {
	g.xAxis = &canvas.Line{}
	g.yAxis = &canvas.Line{}
	g.xAxis.StrokeWidth = 2
	g.xAxis.StrokeColor = foregroundColor()
	g.yAxis.StrokeWidth = 2
	g.yAxis.StrokeColor = foregroundColor()

	axisRange := g.createAxisRange(&size, data)

	g.addAxes(&size, &axisRange)

	g.addXTicks(&size, &axisRange, data)

	g.addYTicks(&size, &axisRange)

	g.addLines(&size, &axisRange, data)

	g.addGraphObjects(graphContainer)
}

func calcMaxTextWidth(a, b float32) float32 {
	// TODO: Make this function generic to take variadic arguments of float32
	minText := canvas.NewText(strconv.Itoa(int(a)), foregroundColor())
	maxText := canvas.NewText(strconv.Itoa(int(b)), foregroundColor())
	return float32(math.Max(float64(minText.MinSize().Width), float64(maxText.MinSize().Width)))
}

func (g *GraphStruct) createAxisRange(size *fyne.Size, data []float32) axisRange {
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
	realizedMin := yMin - float32(orderMagnitude)/2
	realizedMax := yMax + float32(orderMagnitude)/2
	zeroHeight := linearMap(0, realizedMin, realizedMax, size.Height, 0)

	tickSize := float32(orderMagnitude)
	tickMin := float32(math.Round(float64(math.Round(float64(yMin/tickSize)))) * float64(tickSize))
	tickMax := float32(math.Round(float64(math.Round(float64(yMax/tickSize)))) * float64(tickSize))
	numTicks := int(math.Round(math.Abs(float64(yMax-yMin))/float64(orderMagnitude))) + 1
	yAxisOffset := calcMaxTextWidth(tickMin, tickMax)
	tickLength := float32(0.0125 * math.Max(float64(size.Width), float64(size.Height)))
	return axisRange{
		min:         yMin,
		max:         yMax,
		realizedMin: realizedMin,
		realizedMax: realizedMax,
		tickSize:    tickSize,
		numTicks:    numTicks,
		zeroHeight:  zeroHeight,
		yAxisOffset: yAxisOffset + 10,
		tickMin:     tickMin,
		tickMax:     tickMax,
		tickLength:  tickLength,
	}
}

func (g *GraphStruct) addAxes(size *fyne.Size, axisRange *axisRange) {
	g.xAxis.Position1 = fyne.NewPos(axisRange.yAxisOffset, axisRange.zeroHeight)
	g.xAxis.Position2 = fyne.NewPos(size.Width, axisRange.zeroHeight)
	g.yAxis.Position1 = fyne.NewPos(axisRange.yAxisOffset, 0)
	g.yAxis.Position2 = fyne.NewPos(axisRange.yAxisOffset, size.Height)
}

func positionXLabel(index int, length int, xLabel *canvas.Text, size *fyne.Size, axisRange *axisRange) fyne.Position {
	xPos := axisRange.yAxisOffset + float32(index)*size.Width/float32(length) + xLabel.Size().Width/2
	yPos := axisRange.zeroHeight + axisRange.tickLength/2 + 3
	return fyne.NewPos(xPos, yPos)
}

func (g *GraphStruct) addXTicks(size *fyne.Size, axisRange *axisRange, data []float32) {
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
		xLabel.Move(positionXLabel(index, len(data), xLabel, size, axisRange))
		xTick.Position1 = fyne.NewPos(axisRange.yAxisOffset+float32(index)*size.Width/float32(len(data)), axisRange.zeroHeight+(axisRange.tickLength/2))
		xTick.Position2 = fyne.NewPos(axisRange.yAxisOffset+float32(index)*size.Width/float32(len(data)), axisRange.zeroHeight-(axisRange.tickLength/2))
		xTick.StrokeColor = foregroundColor()
		xTick.StrokeWidth = 2
		// Skip 0 label
		if index != 0 {
			g.xLabels = append(g.xLabels, xLabel)
		}
		g.xTicks = append(g.xTicks, xTick)
	}
}

func (g *GraphStruct) addYTicks(size *fyne.Size, axisRange *axisRange) {
	g.yLabels = []*canvas.Text{}
	g.yTicks = []*canvas.Line{}
	for index := 0; index < axisRange.numTicks; index++ {
		yTick := &canvas.Line{}
		tickMin := float32(math.Round(float64(math.Round(float64(axisRange.min/axisRange.tickSize)))) * float64(axisRange.tickSize))
		tickMax := float32(math.Round(float64(math.Round(float64(axisRange.max/axisRange.tickSize)))) * float64(axisRange.tickSize))
		yValue := linearMap(float32(index), 0, float32(axisRange.numTicks)-1, tickMin, tickMax)
		yLabel := canvas.NewText(strconv.Itoa(int(math.Round(float64(yValue)))), color.White)
		tickHeight := linearMap(yValue, axisRange.realizedMin, axisRange.realizedMax, size.Height, 0)
		yLabel.Move(fyne.NewPos(0, tickHeight-yLabel.MinSize().Height/2))
		yTick.Position1 = fyne.NewPos(axisRange.yAxisOffset+axisRange.tickLength, tickHeight)
		yTick.Position2 = fyne.NewPos(axisRange.yAxisOffset, tickHeight)
		yTick.StrokeColor = foregroundColor()
		yTick.StrokeWidth = 2
		g.yTicks = append(g.yTicks, yTick)
		g.yLabels = append(g.yLabels, yLabel)
	}
}

func (g *GraphStruct) addLines(size *fyne.Size, axisRange *axisRange, data []float32) {
	g.lines = []*canvas.Line{}
	for index := range data {
		if index == 0 {
			continue
		}
		line := &canvas.Line{}
		line.Position1 = fyne.NewPos(axisRange.yAxisOffset+linearMap(float32(index-1), 0, float32(len(data)), 0, size.Width), linearMap(data[index-1], axisRange.realizedMin, axisRange.realizedMax, size.Height, 0))
		line.Position2 = fyne.NewPos(axisRange.yAxisOffset+linearMap(float32(index), 0, float32(len(data)), 0, size.Width), linearMap(data[index], axisRange.realizedMin, axisRange.realizedMax, size.Height, 0))
		line.StrokeColor = primaryColor()
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
	for _, yLabel := range g.yLabels {
		graphContainer.Objects = append(graphContainer.Objects, yLabel)
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

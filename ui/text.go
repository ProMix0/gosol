package ui

import (
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Text is a widget that displays a a multiline of text
type Text struct {
	WidgetBase
	text       string
	lines      []string
	lineHeight int
}

func (w *Text) createImg() *ebiten.Image {
	if w.lines == nil {
		log.Panic("widget Text.createImg with no lines")
	}
	dc := gg.NewContext(w.width, w.height)

	dc.SetRGBA(1, 1, 1, 1)
	// nota bene - text is drawn with y as a baseline
	dc.SetFontFace(schriftbank.RobotoRegular14)
	y := w.lineHeight
	for _, str := range w.lines {
		dc.DrawString(str, 0, float64(y-8)) // move up a little to stop descenders being clipped on last line
		y += w.lineHeight
	}
	// uncomment this line to visualize text box
	// dc.DrawLine(0, 0, float64(w.width), float64(w.height))
	// dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

func (w *Text) calcHeights() {
	dc := gg.NewContext(w.width, 48)
	dc.SetFontFace(schriftbank.RobotoRegular14)
	// MeasureString says this text, requested to be 48 high, is 14 high
	w.lines = dc.WordWrap(w.text, float64(w.width-48)) // 24 padding left and right
	w.lineHeight = 24
	w.height = w.lineHeight * len(w.lines) // + w.lineHeight
}

// NewText creates a new Text
func NewText(parent Container, text string) *Text {
	width, _ := parent.Size()
	// widget x, y will be set by LayoutWidgets
	// widget height will be set when wordwrapping in createImg
	w := &Text{
		WidgetBase: WidgetBase{parent: parent, img: nil, width: width},
		text:       text}
	w.calcHeights()
	return w
}

// Activate tells the input we need notifications
func (w *Text) Activate() {
	w.disabled = false
	w.img = w.createImg()
	// w.input.Add(w)
}

// Deactivate tells the input we no longer need notofications
func (w *Text) Deactivate() {
	w.disabled = true
	w.img = w.createImg()
	// w.input.Remove(w)
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (w *Text) NotifyCallback(v input.StrokeEvent) {
	if w.disabled {
		return
	}
	// switch v.Event {
	// case "tap":
	// 	if util.InRect(v.X, v.Y, w.OffsetRect) {
	// 		println("Text NotifyCallback", w.text)
	// 	}
	// }
}

// Update the state of this widget
func (w *Text) Update() {
}

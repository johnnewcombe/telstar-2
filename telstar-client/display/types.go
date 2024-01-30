package display

import (
	"fyne.io/fyne/v2"
	"image/color"
)

type Screen struct {
	disabledRows     []bool
	cursorCol        int
	cursorRow        int
	cursorOn         bool
	textSize         float32
	graphics         bool
	attributes       []Attributes
	escape           bool
	keyShift         bool
	keyControl       bool
	rootCanvas       fyne.Canvas
	rootContainer    *fyne.Container
	flashHidePeriod  bool
	alphaCharMap     map[int]int
	colourMap        map[int]color.Color
	revealMode       bool
	holdGraphicsMode bool
	doubleHeightMode bool
	heldGraphic      HeldGraphic//int
	debugBuffer      []byte
	Debug            bool
	debugTmp         int
}

type HeldGraphic struct
{
	character int
	attributes Attributes
}

/*
type Character struct {
	translatedValue   int
	doubleHeightUpper int
	doubleHeightLower int
}
*/

type Attributes struct {
	controlValue  byte
	graphics      bool
	doubleHeight  bool
	foreColour    color.Color
	backColour    color.Color
	flash         bool
	conceal       bool
	nonContiguous bool
	cursor        bool
}

// copyTo copies the attributes to the specified destination, note that
// this function does not copy the control value, nor the cursor value
func (a *Attributes) copyTo(destination *Attributes) {

	destination.controlValue = 0
	destination.graphics = a.graphics
	destination.doubleHeight = a.doubleHeight
	destination.backColour = a.backColour
	destination.foreColour = a.foreColour
	destination.conceal = a.conceal
	destination.flash = a.flash
	destination.nonContiguous = a.nonContiguous
	//destination.cursor = a.cursor

}

func (a *Attributes) clear() {
	a.controlValue = 0
	a.graphics = false
	a.doubleHeight = false
	a.backColour = colourMap[0]
	a.foreColour = colourMap[7]
	a.conceal = false
	a.flash = false
	a.nonContiguous = false
	a.cursor = false
}

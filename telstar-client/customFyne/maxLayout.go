package customFyne

// Package layout defines the various layouts available to Fyne apps

import "C"
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*maxLayout)(nil)

type maxLayout struct {
}

// NewMaxLayout creates a new MaxLayout instance
func NewMaxLayout() fyne.Layout {
	return &maxLayout{}
}

// Layout is called to pack all child objects into a specified size.
// For MaxLayout this sets all children to the full size passed.
func (m *maxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {

	//fmt.Printf("New Height:%f, New Width:%f\r\n", size.Height, size.Width)

	topLeft := fyne.NewPos(0, 0)

	/* This was suggested at stack overflow
	// set the text size based on the new size
	textSize := size.Height * 0.9
	// resize the child (txt canvas)
	txt.TextSize = textSize
	txt.Refresh()


	*/
	// this is from the bud=gs example
	//textMin := fyne.MeasureText(" ", textSize, fyne.TextStyle{Monospace: true})

	for _, child := range objects { //

		// something special needed if the child is a text canvas
		if _, ok := child.(*canvas.Text); ok {

			// get the text canvas
			txt := child.(*canvas.Text)

			textSize := getTextSize(size, txt.TextSize, 10.0)

			//FIXME could this be done at block layout level and passed in a property? That would save repeating
			// the task for each character? althoght each char could be a different height.
			/*
				Essentially you would have to calculate a new font size based on the input fyne.Size
				into your Layout function. Once calculated you can set it as
				    myText.TextSize = newFontSize;
				    myText.Refresh().

				This will work instead of the Resize which merely sets the space available.

				You can see an example in a Fyne demo repo https://github.com/fyne-io/examples/blob/06b751e39187df843fd23e1215f55fb41845f5a5/bugs/button.go#L36
			*/

			// resize the child (txt canvas)
			if txt.TextSize != textSize {
				txt.TextSize = textSize
				txt.Refresh()
			}

			//cSize := getCanvasSize()
			//tSize :=getTextSize(cSize, txt.TextSize)
			//fmt.Printf("Canvas Height:%f, Text Size:%f\r\n",cSize.Height, tSize)

			//txt.TextSize = textSize
			//txt.Resize(fyne.NewSize(size.Width, textMin.Height))
			//txt.Move(fyne.NewPos(0, (size.Height-textMin.Height)/2))

		} else {
			child.Resize(size)
		}
		child.Move(topLeft)

	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For MaxLayout this is determined simply as the MinSize of the largest child.
func (m *maxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	return minSize
}

// getCanvasSize: gets the minimum canvas size for the current text size
func getCanvasSize(textSize float32) fyne.Size {

	//var textSize float32 = 14.0 //canvasSize.Height * .67
	var text string = " " // monotext font so all chars are the same size including space.

	// RenderedTextSize returns the size required to render the given string of specified
	// font size and style. It also returns the height to text baseline, measured from the top.
	s, _ := fyne.CurrentApp().Driver().RenderedTextSize(text, textSize, fyne.TextStyle{Monospace: true})
	//fmt.Printf("%v %f ", s, b)
	return s
}

// getTextSize returns the maximum text size that can be set for the given canvas size,
// this is the opposite to getCanvasSize()
func getTextSize(newCanvasSize fyne.Size, currentTextSize float32, step float32) float32 {

	const(
		AllowableError float32 = 0.002
	)

	var (
		newTextSize float32
	)

	// ignore if height is zero, this occurs during initialisation
	if newCanvasSize.Height == 0 {
		return currentTextSize
	}

	// TODO insist on step being positive?
	/*
		We are trying to establish the correct text size for a specific canvas size
		we know the current canvas size
	*/

	//get the size of the canvas for the current text size
	textCanvasSize := getCanvasSize(currentTextSize)

	//fmt.Printf("%f\r\n", currentTextCanvasSize.Height)
	//fmt.Printf("Text Height: %f, Text Width: %f, Text Size: %f\r\n", currentTextCanvasSize.Height, currentTextCanvasSize.Width, currentTextSize)

	complete := func() bool {
		return textCanvasSize.Height <= newCanvasSize.Height &&
			textCanvasSize.Height > (newCanvasSize.Height-AllowableError)
	}

	// if the height of the current text canvas is the same as the new text canvas height the we have calculated it within limits
	if complete() {
		// all done
		return currentTextSize
	}

	//print("Recalculate\r\n")
	newTextSize = currentTextSize

	// make sure the text is smaller than the space
	for textCanvasSize.Height > newCanvasSize.Height {
		newTextSize = newTextSize * .67
		textCanvasSize = getCanvasSize(newTextSize)
	}

	i:=0 // safety measure
	// at this point the newTextSize is still too small so increase it by the step size
	for !complete(){

		newTextSize = newTextSize + step
		textCanvasSize = getCanvasSize(newTextSize)

		//fmt.Printf("index: %d step:%f newTextSize:%f textCanvasSize.Height:%f newCanvasSize.Height:%f\r\n",i, step, newTextSize, textCanvasSize.Height , newCanvasSize.Height )

		if textCanvasSize.Height > newCanvasSize.Height  {
			//print("Step Change\r\n")
			// the step was too big so use previous text size
			newTextSize = newTextSize - step
			step = step/10
		}
		i++
		if i>1000 {
			break
		}
	}
	//print("Complete!\r\n")
	return newTextSize
}

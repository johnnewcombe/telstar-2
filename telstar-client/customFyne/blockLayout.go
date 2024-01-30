package customFyne

import (
	"fyne.io/fyne/v2"
	"math"
)

type CharType int

const (
	NonContiguousBorder = 1.5
)

const (
	CharTypeNone CharType = iota
	CharTypeDoubleHeightUpper
	CharTypeDoubleHeightLower
)

type mosaicLayout struct {
	Cols                    int
	Separated, HoldGraphics bool
	CharType                CharType
	vertical, adapt         bool
}

// Layout is called to pack all child objects into a specified size.
// For a Mosaic Layout this will pack objects into a table format with three rows of two
// columns.
func (m *mosaicLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {

	var (
		xp, yp, xr, yr float32
		x1, y1, x2, y2 float32
	)

	// keep things limited to 6 objects
	if len(objects) != 6 {
		return
	}

	rows := m.countRows(objects)

	cellWidth := float64(size.Width) / float64(m.Cols)
	cellHeightUnit := float64(size.Height) / 10.0

	if !m.horizontal() {
		cellWidth = float64(size.Width) / float64(rows)
	}

	row, col := 0, 0

	for i := 0; i < 6; i++ {

		child := objects[i]

		if !child.Visible() {
			continue
		}

		/*
			Double height uses two sets of sixels.

				|---|---|
				| 0 | 1 |
				|---|---|
				| 2 | 3 |
				|---|---|
				| 4 | 5 |
				|---|---|

				|---|---|
				| 0 | 1 |
				|---|---|
				| 2 | 3 |
				|---|---|
				| 4 | 5 |
				|---|---|

		*/

		// gets the bounds of the sixel
		x1, y1, x2, y2 = m.getBounds(row, col, cellHeightUnit, cellWidth)

		// Calculates the correct size and position of the sixel based on where it sits in the block, caters for
		// upper and lower double height.
		if m.Separated && !m.HoldGraphics {

			// FIXME If we are using a Hold Graphic, we ignore separated unless the hold graphic is separated.

			// x,y for size for normal height separated graphics
			xr = x2 - x1 - (2 * NonContiguousBorder)
			yr = y2 - y1 - (2 * NonContiguousBorder)

			// x,y for position for normal height separated graphics
			xp = x1 + NonContiguousBorder
			yp = y1 + NonContiguousBorder

			// make any adjustments necessary for double height
			if m.CharType == CharTypeDoubleHeightUpper {

				if i == 2 || i == 3 {
					yp = y1
				} else {
					yp = y1 + (2 * NonContiguousBorder)
				}

			} else if m.CharType == CharTypeDoubleHeightLower {

				if i == 2 || i == 3 {
					yp = y1 + (2 * NonContiguousBorder)
				} else {
					yp = y1
				}
			}

			child.Move(fyne.NewPos(xp, yp))
			child.Resize(fyne.NewSize(xr, yr))

		} else {
			child.Move(fyne.NewPos(x1, y1))
			child.Resize(fyne.NewSize(x2-x1, y2-y1))
		}

		if m.horizontal() {
			if (i+1)%m.Cols == 0 {
				row++
				col = 0
			} else {
				col++
			}
		} else {
			if (i+1)%m.Cols == 0 {
				col++
				row = 0
			} else {
				row++
			}
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a GridLayout this is the size of the largest child object multiplied by
// the required number of columns and rows.
func (m *mosaicLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	rows := m.countRows(objects)
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	if m.horizontal() {
		minContentSize := fyne.NewSize(minSize.Width*float32(m.Cols), minSize.Height*float32(rows))
		return minContentSize.Add(fyne.NewSize(0, 0))
	}

	minContentSize := fyne.NewSize(minSize.Width*float32(rows), minSize.Height*float32(m.Cols))
	return minContentSize.Add(fyne.NewSize(0, 0))
}

func (m *mosaicLayout) horizontal() bool {
	if m.adapt {
		return fyne.IsHorizontal(fyne.CurrentDevice().Orientation())
	}

	return !m.vertical
}

func (m *mosaicLayout) countRows(objects []fyne.CanvasObject) int {
	count := 0
	for _, child := range objects {
		if child.Visible() {
			count++
		}
	}

	return int(math.Ceil(float64(count) / float64(m.Cols)))
}

// Get the leading (top or left) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func (m *mosaicLayout) getLeading(size float64, offset int) float32 {
	ret := size * float64(offset)
	return float32(math.Round(ret))
}

// Get the trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func (m *mosaicLayout) getTrailing(size float64, offset int) float32 {
	return m.getLeading(size, offset+1)
}

func (m *mosaicLayout) getBounds(row int, col int, heightUnit float64, width float64) (float32, float32, float32, float32) {

	var x1, y1, x2, y2 float64

	x1 = float64(m.getLeading(width, col))
	x2 = float64(m.getTrailing(width, col))

	//
	switch row {
	//row 0

	case 0:
		switch m.CharType {
		case CharTypeDoubleHeightUpper:
			y1 = 0.0
			y2 = y1 + heightUnit*3.0 // height is 3
		case CharTypeDoubleHeightLower:
			y1 = 0.0
			y2 = y1 + heightUnit*4.0 // height is 4
		default:
			y1 = 0.0
			y2 = y1 + heightUnit*3.0 // height is 3
		}

	case 1:
		//row 1
		switch m.CharType {
		case CharTypeDoubleHeightUpper:
			y1 = heightUnit * 3.0    // previous row height was 3
			y2 = y1 + heightUnit*3.0 // height is 3
		case CharTypeDoubleHeightLower:
			y1 = heightUnit * 4.0    // previous row height was 4
			y2 = y1 + heightUnit*3.0 // height is 3
		default:
			y1 = heightUnit * 3.0    // previous row height was 3
			y2 = y1 + heightUnit*4.0 // height is 4
		}

	case 2:
		//row 2
		switch m.CharType {
		case CharTypeDoubleHeightUpper:
			y1 = heightUnit * 6.0    // previous row height was 3 + 3
			y2 = y1 + heightUnit*4.0 // height is 4
		case CharTypeDoubleHeightLower:
			y1 = heightUnit * 7.0    // previous row height was 4 + 4
			y2 = y1 + heightUnit*3.0 // height is 3
		default:
			y1 = heightUnit * 7.0    // previous row height was 3 + 4
			y2 = y1 + heightUnit*3.0 // height is 3
		}
	}

	return float32(x1), float32(y1), float32(x2), float32(y2)
}

func NewBlockLayout() fyne.Layout {
	return &mosaicLayout{Cols: 2}
}

package display

import (
	"bitbucket.org/johnnewcombe/telstar-client/customFyne"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/johnnewcombe/telstar-client/constants"
)

func (s *Screen) isLastRow(row int) bool {
	return row == constants.ROWS-1
}

func (s *Screen) isLastColumn(col int) bool {
	return col == constants.COLS-1
}

func (s *Screen) isControlCode(controlCode int) bool {
	return controlCode >= 0x40 && controlCode <= 0x7f
}

func (s *Screen) isBlastThrough(char int) bool {

	// includes double height upper and lower characters
	return (char >= 0x40 && char <= 0x5f) ||
		(char >= 0xe040 && char <= 0xe05f) ||
		(char >= 0xe140 && char <= 0xe15f)
}

func (s *Screen) isRowDisabled(row int) bool {

	// ignore cant write to the lower row of a pair where the upper row has a double height
	if row > 0 {
		// check the row above for a DH
		for col := 0; col < constants.COLS; col++ {

			if s.attributes[s.getIndex(col, row-1)].controlValue == constants.DOUBLE_HEIGHT {
				// row above has a DH so this row should not be written to
				// TODO place any chars in a buffer so that they can be recovered should the DH
				//  be deleted or overwritten
				return true
			}
		}
	}
	return false

}

func (s *Screen) isDoubleHeightRow(row int) bool {

	for col := 0; col < constants.COLS; col++ {
		if s.attributes[s.getIndex(col, row)].controlValue == constants.DOUBLE_HEIGHT {
			return true
		}
	}
	return false
}

func (s *Screen) isDoubleHeightLowerRow(row int) bool {

	if row > 0 && s.isDoubleHeightRow(row-1) {
		return true
	}
	return false
}

func (s *Screen) clearToEndOfLine(cursorCol int, cursorRow int) {
	for col := cursorCol; col < constants.COLS; col++ {
		s.setText(0x20, col, cursorRow, false, customFyne.CharTypeNone)
	}
}

func (s *Screen) getIndex(cursorCol int, cursorRow int) int {
	return (cursorRow * constants.COLS) + cursorCol
}

// getRowIndex return index of the first position of the specified row.
func (s *Screen) getRowIndex(row int) int {
	return row * constants.COLS
}

// getColRow returns the column and row from the index
func (s *Screen) getColRow(index int) (int, int) {
	rows := index / constants.COLS // integer division, decimals are truncated
	cols := index % constants.COLS
	return cols, rows
}

// getRowEndIndex return index of the last position of the specified row.
func (s *Screen) getRowEndIndex(row int) int {
	return row*constants.COLS + constants.COLS - 1
}

func (s *Screen) setCursor(index int) {

	// a cursor can appear at every char position so
	// before setting one, clear the whole display
	for i := 0; i < constants.ROWS*constants.COLS; i++ {
		cursor := s.getCursorCanvas(i)
		cursor.Hidden = true
		s.attributes[i].cursor = false
	}
	// set the current cursor on screen and in the attributes
	// attributes are used to determine which items to flash
	cursor := s.getCursorCanvas(index)
	s.attributes[index].cursor = true
	cursor.Hidden = !s.cursorOn
}

func (s *Screen) getBackgroundCanvas(index int) *canvas.Rectangle {

	if index >= constants.COLS*constants.ROWS {
		return nil
	}

	maxContainer := s.rootContainer.Objects[index].(*fyne.Container)
	bg := maxContainer.Objects[0].(*canvas.Rectangle)
	return bg
}

func (s *Screen) getMosaicContainer(index int) *customFyne.Block {

	if index >= constants.COLS*constants.ROWS {
		return nil
	}

	maxContainer := s.rootContainer.Objects[index].(*fyne.Container)
	c := maxContainer.Objects[1].(*fyne.Container)
	block := customFyne.Block{c}
	return &block
}

// func (s *Screen) getTextCanvas(index int) *bugButton {
func (s *Screen) getTextCanvas(index int) *canvas.Text {

	if index >= constants.COLS*constants.ROWS {
		return nil
	}

	maxContainer := s.rootContainer.Objects[index].(*fyne.Container)
	txt := maxContainer.Objects[2].(*canvas.Text)
	txt.TextStyle = fyne.TextStyle{Monospace: true}
	return txt
}

func (s *Screen) getCursorCanvas(index int) *canvas.Rectangle {

	if index >= constants.COLS*constants.ROWS {
		return nil
	}

	maxContainer := s.rootContainer.Objects[index].(*fyne.Container)
	bg := maxContainer.Objects[3].(*canvas.Rectangle)
	return bg
}

func (s *Screen) roll() {
	if s.cursorRow >= constants.ROWS {
		s.cursorRow = s.cursorRow - constants.ROWS
	}
	if s.cursorRow < 0 {
		s.cursorRow = constants.ROWS + s.cursorRow // add the negative number
	}
}

func (s *Screen) wrap() {
	if s.cursorCol >= constants.COLS {
		s.cursorCol = s.cursorCol - constants.COLS
		s.cursorRow++
	}
	if s.cursorCol < 0 {
		s.cursorCol = constants.COLS + s.cursorCol // add the negative number
		s.cursorRow--
	}
}

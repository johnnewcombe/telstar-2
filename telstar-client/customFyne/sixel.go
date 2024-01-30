package customFyne

import (
	"fyne.io/fyne/v2/canvas"
)

type Sixel struct{
	*canvas.Rectangle
	Separated bool
}

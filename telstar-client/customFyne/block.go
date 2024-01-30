package customFyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

type Block struct {
	*fyne.Container
}

func (b *Block) Clear() {
	for s := 0; s < 6; s++ {
		sixel := b.Objects[s].(*canvas.Rectangle)
		sixel.FillColor = nil
	}
}

func (b *Block) SetMosaic(char byte, colour color.Color, separated bool, charType CharType) {

	layout:=b.Layout.(*mosaicLayout)
	layout.Separated = separated
	layout.CharType = charType

	if colour == nil {
		colour = color.White
	}

	if char >=0x20 && char <= 0x3f{
		char -= 0x20
	} else if char >= 0x60 && char <=0x7f{
		char -= 0x40
	} else {
		// not a graphic character
		return
	}

	// set appropriate sixels
	for s := 0; s < 6; s++ {

		// if associated bit of the byte is set, set the sixel to the foreground color
		if char&0b00000001 == 1 {
			b.Objects[s].(*canvas.Rectangle).FillColor = colour
		} else {
			b.Objects[s].(*canvas.Rectangle).FillColor = nil
		}
		char = char >> 1
	}
}

func NewBlock() *Block {

	mosaic := container.New(NewBlockLayout()) // 0x40 represents a block character
	m := Block{mosaic}

	for s := 0; s < 6; s++ {
		m.Objects = append(m.Objects, new(canvas.Rectangle))
	}
	return &m

}

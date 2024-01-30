package keyboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type funcDef func(byte)

type Keyboard struct {
	keyShift   bool
	keyControl bool
	keyMap     map[string]Key
	fOutput    funcDef
}

type Key struct {
	Unmodified byte
	Shift      byte
	Control    byte
}

func (k *Keyboard) keyDown(key *fyne.KeyEvent) {

	var (
		result byte
	)

	switch key.Name {
	case "LeftShift", "RightShift":
		k.keyShift = true
		return
	case "LeftControl", "RightControl":
		k.keyControl = true
		return
	}

	translatedKey := k.keyMap[string(key.Name)]

	if k.keyControl {
		result = translatedKey.Control
		k.keyControl = false
	} else if k.keyShift {
		result = translatedKey.Shift
		k.keyShift = false
	} else {
		result = translatedKey.Unmodified
	}

	k.fOutput(result)
}

func (k *Keyboard) keyUp(key *fyne.KeyEvent) {

	switch key.Name {
	case "LeftShift", "RightShift":
		k.keyShift = false
	case "Control":
		k.keyControl = false
	}
}

func (k *Keyboard) Initialise(w fyne.Window, f funcDef) {

	k.fOutput = f

	// this provides the keypress stuff for desktop only
	// cannot be used on mobile/tablet etc where there is
	// a touch screen not a keyboard
	if deskCanvas, ok := w.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(k.keyDown)
		deskCanvas.SetOnKeyUp(k.keyUp)
	}

	// special cases for key names
	k.keyMap = make(map[string]Key)
	k.keyMap["Return"] = Key{0x5f, 0x0D, 0x5f}
	k.keyMap["Space"] = Key{0x20, 0x20, 0x20}
	k.keyMap["Escape"] = Key{0x1b, 0x1b, 0x1b}
	k.keyMap["Tab"] = Key{0x9, 0x9, 0x9}
	k.keyMap["BackSpace"] = Key{0x8, 0x8, 0x8}
	k.keyMap["Insert"] = Key{0, 0, 0}
	k.keyMap["Delete"] = Key{0x8, 0x8, 0x8}
	k.keyMap["Right"] = Key{0x9, 0x9, 0x9}
	k.keyMap["Left"] = Key{0x8, 0x8, 0x8}
	k.keyMap["Down"] = Key{0xa, 0xa, 0xa}
	k.keyMap["Up"] = Key{0xb, 0xb, 0xb}
	k.keyMap["Prior"] = Key{0, 0, 0}
	k.keyMap["Next"] = Key{0, 0, 0}
	k.keyMap["Home"] = Key{0, 0, 0}
	k.keyMap["End"] = Key{0, 0, 0}
	k.keyMap["F1"] = Key{0xB1, 0xC1, 0xD1}
	k.keyMap["F2"] = Key{0xB2, 0xC2, 0xD2}
	k.keyMap["F3"] = Key{0xB3, 0xC3, 0xD3}
	k.keyMap["F4"] = Key{0xB4, 0xC4, 0xD4}
	k.keyMap["F5"] = Key{0xB5, 0xC5, 0xD5}
	k.keyMap["F6"] = Key{0xB6, 0xC6, 0xD6}
	k.keyMap["F7"] = Key{0xB7, 0xC7, 0xD7}
	k.keyMap["F8"] = Key{0xB8, 0xC8, 0xD8}
	k.keyMap["F9"] = Key{0xB9, 0xC9, 0xD9}
	k.keyMap["F10"] = Key{0xB0, 0xC0, 0xD0}
	k.keyMap["F11"] = Key{0, 0, 0}
	k.keyMap["F12"] = Key{0, 0, 0}
	k.keyMap["KP_Enter"] = Key{0x5f, 0, 0} // keypad
	k.keyMap["0"] = Key{0x30, 0x29, 0}
	k.keyMap["1"] = Key{0x31, 0x21, 0}
	k.keyMap["2"] = Key{0x32, 0x40, 0}
	k.keyMap["3"] = Key{0x33, 0x23, 0}
	k.keyMap["4"] = Key{0x34, 0x24, 0}
	k.keyMap["5"] = Key{0x35, 0x25, 0}
	k.keyMap["6"] = Key{0x36, 0, 0}
	k.keyMap["7"] = Key{0x37, 0x26, 0}
	k.keyMap["8"] = Key{0x38, 0x2a, 0}
	k.keyMap["9"] = Key{0x39, 0x28, 0}
	k.keyMap["A"] = Key{0x61, 0x41, 1}
	k.keyMap["B"] = Key{0x62, 0x42, 2}
	k.keyMap["C"] = Key{0x63, 0x43, 3}
	k.keyMap["D"] = Key{0x64, 0x44, 4}
	k.keyMap["E"] = Key{0x65, 0x45, 5}
	k.keyMap["F"] = Key{0x66, 0x46, 6}
	k.keyMap["G"] = Key{0x67, 0x47, 7}
	k.keyMap["H"] = Key{0x68, 0x48, 8}
	k.keyMap["I"] = Key{0x69, 0x49, 9}
	k.keyMap["J"] = Key{0x6a, 0x4a, 10}
	k.keyMap["K"] = Key{0x6b, 0x4b, 11}
	k.keyMap["L"] = Key{0x6c, 0x4c, 12}
	k.keyMap["M"] = Key{0x6d, 0x4d, 13}
	k.keyMap["N"] = Key{0x6e, 0x4e, 14}
	k.keyMap["O"] = Key{0x6f, 0x4f, 15}
	k.keyMap["P"] = Key{0x70, 0x50, 16}
	k.keyMap["Q"] = Key{0x71, 0x51, 17}
	k.keyMap["R"] = Key{0x72, 0x52, 18}
	k.keyMap["S"] = Key{0x73, 0x53, 19}
	k.keyMap["T"] = Key{0x74, 0x54, 20}
	k.keyMap["U"] = Key{0x75, 0x55, 21}
	k.keyMap["V"] = Key{0x76, 0x56, 22}
	k.keyMap["W"] = Key{0x77, 0x57, 23}
	k.keyMap["X"] = Key{0x78, 0x58, 24}
	k.keyMap["Y"] = Key{0x79, 0x59, 25}
	k.keyMap["Z"] = Key{0x7a, 0x5a, 26}
	k.keyMap["'"] = Key{0x27, 0x22, 0}
	k.keyMap[","] = Key{0x2c, 0x3c, 0}
	k.keyMap["-"] = Key{0x2d, 0x2d, 0}
	k.keyMap["."] = Key{0x2e, 0x3e, 0}
	k.keyMap["/"] = Key{0x2f, 0x3f, 0}
	k.keyMap["\\"] = Key{0xBc, 0, 0}
	k.keyMap["["] = Key{0x5b, 0x5b, 0}
	k.keyMap["]"] = Key{0x5d, 0x5d, 0}
	k.keyMap[";"] = Key{0x3b, 0x3a, 0}
	k.keyMap["="] = Key{0x3d, 0x2b, 0}
	k.keyMap["*"] = Key{0x2a, 0x2a, 0}
	k.keyMap["+"] = Key{0x2b, 0x2b, 0}
	k.keyMap["`"] = Key{0, 0, 0}

}

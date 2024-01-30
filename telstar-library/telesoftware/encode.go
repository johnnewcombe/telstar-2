package telesoftware

import (
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"fmt"
)

const (
	BLOCK_START      = "\x7C\x41" // 'A'
	BLOCK_END        = "\x7C\x5A" // 'Z'
	BLOCK_G          = "\x7C\x47" // G
	BLOCK_I          = "\x7C\x49" // I
	EOL              = "\x7c\x4c" // EOL

	// MAX_CHARS_PER_BLOCK spec not clear whether 859 includes all block/chksum bytes
	// in this implementation all data and block markers and chksum falls within this value.
)

type Block struct {
	Data       []byte
	PageNumber int
	FrameId    rune
}

// need to call Encode then Enblock before creating the frames

// Encode returns data in Telesoftware encoded format.
func Encode(data []byte) ([]byte, error) {

	var (
		err          error
		rawByteCount int
		//encodedData     strings.Builder
		encodedData     []byte
		controlSequence byte
	)

	for i := 0; i < len(data); i++ {

		rawByteCount++

		byt := data[i]

		// note the the order of these statements is important
		if byt == 0x0d {

			if err = appendEncodedData(0x4c, &encodedData, true); err != nil {
				return encodedData, err
			}

		} else if byt == 0x20 {

			if err = updateControlSequence(&controlSequence, 0, &encodedData); err != nil {
				return encodedData, err
			}
			if err = appendEncodedData(0x7d, &encodedData, false); err != nil {
				return encodedData, err
			}

		} else if byt == 0x7c {
			if err = appendEncodedData(0x45, &encodedData, true); err != nil {
				return encodedData, err
			}

		} else if byt == 0x7d {
			if err = appendEncodedData(0x7d, &encodedData, true); err != nil {
				return encodedData, err
			}

		} else if byt == 0x7e {
			if err = updateControlSequence(&controlSequence, 0, &encodedData); err != nil {
				return encodedData, err
			}
			if err = appendEncodedData(byt, &encodedData, false); err != nil {
				return encodedData, err
			}
		} else if byt < 0x20 {
			if err = updateControlSequence(&controlSequence, 1, &encodedData); err != nil {
				return encodedData, err
			}
			byt += 0x40
			if err = appendEncodedData(byt, &encodedData, false); err != nil {
				return encodedData, err
			}
		} else if byt < 0x80 {
			if err = updateControlSequence(&controlSequence, 0, &encodedData); err != nil {
				return encodedData, err
			}
			if err = appendEncodedData(byt, &encodedData, false); err != nil {
				return encodedData, err
			}
		} else if byt < 0xa0 {
			if err = updateControlSequence(&controlSequence, 2, &encodedData); err != nil {
				return encodedData, err
			}
			byt -= 0x40
			if err = appendEncodedData(byt, &encodedData, false); err != nil {
				return encodedData, err
			}
		} else if byt < 0xc0 {
			if err = updateControlSequence(&controlSequence, 3, &encodedData); err != nil {
				return encodedData, err
			}
			byt -= 0x60
			if err = appendEncodedData(byt, &encodedData, false); err != nil {
				return encodedData, err
			}
		} else if byt < 0xe0 {
			if err = updateControlSequence(&controlSequence, 4, &encodedData); err != nil {
				return encodedData, err
			}
			byt -= 0x80
			if err = appendEncodedData(byt, &encodedData, false); err != nil {
				return encodedData, err
			}
		} else {
			if err = updateControlSequence(&controlSequence, 5, &encodedData); err != nil {
				return encodedData, err
			}
			byt -= 0xA0
			if err = appendEncodedData(byt, &encodedData, false); err != nil {
				return encodedData, err
			}
		}
	}
	if err = appendEncodedData(0x46, &encodedData, true); err != nil {
		return encodedData, err
	}
	return encodedData, nil
}

// Enblock This method creates the correct number of blocks based on the specified data.
// The each block represents a frame starting at pageId.
// The data should have been encoded before being passed to this function
func Enblock(encodedData []byte, pageNumber int, frameId rune, name string, maxCharsPerBlock int) ([]Block, error) {

	var (
		blocks         []Block
		block          Block
		charsAvailable int
		checksum       byte
		err        error
	)

	charsAvailable = maxCharsPerBlock

	// create the header block
	// add the header to the result (data to come later)
	blocks = append(blocks, Block{[]byte{}, pageNumber, frameId})

	// create the first program block
	if pageNumber, frameId, err = utils.GetFollowOnPID(pageNumber, frameId); err != nil {
		return blocks, err
	}

	block = Block{[]byte(BLOCK_START + BLOCK_G + string(frameId) + BLOCK_I), pageNumber, frameId}

	for i := 0; i < len(encodedData); i++ {

		// TODO: Tests need to be made in relation to a 7c sequence crossing a boundary, is it allowed?

		byt := encodedData[i]
		if charsAvailable == 0 || (charsAvailable == 1 && byt == 0x7c) {

			// end of a block so tidy up
			block.Data = append(block.Data, BLOCK_END...)
			if checksum, err = calculateChecksum(block.Data); err != nil {
				return blocks, err
			}
			block.Data = append(block.Data, fmt.Sprintf("%03d", checksum)...)

			// add the completed block to the result
			blocks = append(blocks, block)

			// start a new block by getting the next frame number
			// get the next pageId as this will handle zero page routing etc
			if pageNumber, frameId, err = utils.GetFollowOnPID(pageNumber, frameId); err != nil {
				return blocks, err
			}
			block = Block{[]byte{}, pageNumber, frameId}
			block.Data = append(block.Data, BLOCK_START+BLOCK_G+string(frameId)+BLOCK_I...)
			charsAvailable = maxCharsPerBlock

			// add the byte to the new block and update the chars available
			block.Data = append(block.Data, byt)
			charsAvailable--

		} else {
			// add the byte to the block and update the chars available
			block.Data = append(block.Data, byt)
			charsAvailable--
		}
	}
	// complete the last block
	block.Data = append(block.Data, BLOCK_END...)
	if checksum, err = calculateChecksum(block.Data); err != nil {
		return blocks, err
	}
	block.Data = append(block.Data, fmt.Sprintf("%03d", checksum)...)

	// if block is not empty, add it to the list
	if len(block.Data) > 6 {
		blocks = append(blocks, block)
	}

	//create the header data
	//blockCount = byte(len(blocks))
	blockCountS := fmt.Sprintf("%03d", len(blocks)-1)
	blocks[0].Data = []byte(BLOCK_START)
	blocks[0].Data = append(blocks[0].Data, BLOCK_G...)
	blocks[0].Data = append(blocks[0].Data, byte(blocks[0].FrameId))
	blocks[0].Data = append(blocks[0].Data, BLOCK_I...)
	blocks[0].Data = append(blocks[0].Data, name...)
	blocks[0].Data = append(blocks[0].Data, EOL...)
	blocks[0].Data = append(blocks[0].Data, blockCountS...) // needs to be padded to three bytes
	blocks[0].Data = append(blocks[0].Data, BLOCK_END...)

	if checksum, err = calculateChecksum(blocks[0].Data); err != nil {
		return blocks, err
	}
	blocks[0].Data = append(blocks[0].Data, fmt.Sprintf("%03d", checksum)...)

	return blocks, nil
}

func updateControlSequence(currentControlSequence *byte, newControlSequence byte, encodedData *[]byte) error {

	if newControlSequence != *currentControlSequence {
		*currentControlSequence = newControlSequence
		if err := appendEncodedData(newControlSequence+0x30, encodedData, true); err != nil {
			return err
		}
	}
	return nil
}

func appendEncodedData(data byte, encodedData *[]byte, escape bool) error {

	if data > 0x7f {
		return fmt.Errorf("byte is out of range (0-7f), byte was %d", data)
	}
	if escape {
		*encodedData = append(*encodedData, 0x7c)
	}
	*encodedData = append(*encodedData, data)

	return nil
}

func calculateChecksum(data []byte) (byte, error) {

	var (
		checksum byte
		chkData  []byte
	)

	chkData = data[2 : len(data)-2]
	for i := 0; i < len(chkData); i++ {
		byt := chkData[i]
		checksum ^= byt
	}

	return checksum, nil

}

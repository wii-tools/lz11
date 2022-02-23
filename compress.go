package lz11

import (
	"bytes"
	"encoding/binary"
)

var input *bytes.Buffer
var tempBuffer *bytes.Buffer
var output *bytes.Buffer

func Compress(passed []byte) ([]byte, error) {
	input = bytes.NewBuffer(passed)
	tempBuffer = new(bytes.Buffer)
	output = new(bytes.Buffer)

	size := input.Len()

	// First a sanity check
	if size > RawMax {
		return nil, ErrInputTooLarge
	}

	// Write header
	header := uint32(0x11 + (size << 8))
	err := binary.Write(output, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}

	off := 0
	byteVar := 0
	index := 7

	for off < size {
		pos, length := findLongestMatch(*input, off, 65809)

		if pos == nil {
			index -= 1
			err := tempBuffer.WriteByte(input.Bytes()[off])
			if err != nil {
				panic(err)
			}

			off += 1
		} else {
			lzOff := off - *pos - 1
			byteVar |= 1 << index
			index -= 1

			if *length < 0x11 {
				l := *length - 1
				cmp := []byte{
					byte(lzOff>>8) + byte(l<<4),
					byte(lzOff),
				}

				tempBuffer.Write(cmp)
			} else if *length < 0x111 {
				l := *length - 0x11
				cmp := []byte{
					byte(l >> 4),
					byte(lzOff>>8) + byte(l<<4),
					byte(lzOff),
				}

				tempBuffer.Write(cmp)
			} else {
				l := *length - 0x111
				cmp := []byte{
					byte((l >> 12) + 0x10),
					byte(l >> 4),
					byte((lzOff >> 8) + (l << 4)),
					byte(lzOff),
				}

				tempBuffer.Write(cmp)
			}

			off += *length
		}

		if index < 0 {
			output.WriteByte(byte(byteVar))
			output.Write(tempBuffer.Bytes())
			byteVar = 0
			index = 7
			tempBuffer.Reset()
		}
	}

	if tempBuffer.Len() != 0 {
		output.WriteByte(byte(byteVar))
		output.Write(tempBuffer.Bytes())
	}

	output.WriteByte(byte(0xFF))

	return output.Bytes(), nil
}

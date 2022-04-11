package lz11

import (
	"bytes"
	"encoding/binary"
)

var compressed *bytes.Buffer
var decompressed *bytes.Buffer

// Decompress decompresses the passed data using LZ11.
func Decompress(passed []byte) ([]byte, error) {
	compressed = bytes.NewBuffer(passed)
	// Ensure size validity.
	if CompressedMin > compressed.Len() {
		return nil, ErrCompressedTooSmall
	}
	if compressed.Len() > CompressedMax {
		return nil, ErrCompressedTooLarge
	}
	// Ensure the first byte of this data is 0x11, signifying a proper file.
	if compressed.Next(1)[0] != FileMagic {
		return nil, ErrInvalidMagic
	}

	// Obtain the length of the decompressed file.
	// We then drop the highest byte to strip the 0x11 magic.
	header := append(compressed.Next(3), []byte{0}...)
	originalLen := binary.LittleEndian.Uint32(header)

	// Create a buffer for us to write to throughout decompression.
	decompressed = new(bytes.Buffer)

	for decompressed.Len() < int(originalLen) {
		read := compressed.Next(1)[0]

		for i := 7; i != -1; i-- {
			if decompressed.Len() >= int(originalLen) {
				break
			}

			if (read >> i) & 1 == 0 {
				decompressed.WriteByte(compressed.Next(1)[0])
			} else {
				lenmsb := uint32(compressed.Next(1)[0])
				lsb := uint32(compressed.Next(1)[0])

				length := lenmsb >> 4
				disp := ((lenmsb & 15) << 8) + lsb

				if length > 1 {
					length += 1
				} else if length == 0 {
					length = (lenmsb & 15) << 4
					length += lsb >> 4
					length += 0x11
					msb := uint32(compressed.Next(1)[0])
					disp = ((lsb & 15) << 8) + msb

				} else {
					length = (lenmsb & 15) << 12
					length += lsb << 4
					someBytes := compressed.Next(2)
					length += uint32(someBytes[0]) >> 4
					length += 0x111
					disp = ((uint32(someBytes[0]) & 15) << 8) + uint32(someBytes[1])

				}

				start := decompressed.Len() - int(disp) - 1

				for i1 := 0; i1 < int(length); i1++ {
					val := decompressed.Bytes()[start + i1]
					decompressed.WriteByte(val)
				}
			}
		}
	}

	return decompressed.Bytes(), nil
}

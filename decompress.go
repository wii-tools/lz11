package lz11

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	fmt.Println("Source file is claimed to be", originalLen, "bytes")

	for decompressed.Len() < int(originalLen) {
		currentByte, err := compressed.ReadByte()
		if err != nil {
			return nil, ErrTruncated
		}

		for _, flag := range bits(currentByte) {
			if flag == 0 {
				// Copy a byte as-is.
				nextByte, err := compressed.ReadByte()
				if err != nil {
					return nil, err
				}

				err = decompressed.WriteByte(nextByte)
				if err != nil {
					return nil, err
				}
			} else if flag == 1 {
				// Determine how many times to copy a byte.
				nextByte, err := compressed.ReadByte()
				if err != nil {
					return nil, err
				}

				indicator := nextByte >> 4
				var count int
				if indicator == 0 {
					// 8 bit count, 12 bit disp
					count = int(nextByte << 4)

					// Read a further byte for full displacement.
					nextByte, err = compressed.ReadByte()
					if err != nil {
						return nil, err
					}

					count += int(nextByte) >> 4
					count += 0x11
				} else if indicator == 1 {
					// 16 bit count, 12 bit disp
					count = int((nextByte & 0xf) << 12)
					nextByte, err = compressed.ReadByte()
					if err != nil {
						return nil, err
					}
					count += int(nextByte << 4)

					nextByte, err = compressed.ReadByte()
					if err != nil {
						return nil, err
					}
					count += int(nextByte) >> 4
					count += 0x111
				} else {
					// Indicator is count, 12 bit disp
					count = int(indicator)
					count += 1
				}

				// Determine the offset to copy from within decompressed.
				disp := (int(nextByte) & 0xf) << 8
				dispByte, err := compressed.ReadByte()
				if err != nil {
					return nil, err
				}
				disp += int(dispByte)
				disp += 1

				// Copy the offset of bytes from the current decompressed buffer for the specified amount.
				for count != 0 {
					current := decompressed.Bytes()
					err = decompressed.WriteByte(current[len(current)-disp])
					if err != nil {
						return nil, err
					}

					count--
				}
			} else {
				return nil, ErrInvalidData
			}

			if int(originalLen) <= decompressed.Len() {
				break
			}
		}
	}

	return decompressed.Bytes(), nil
}

// bits returns a byte array with the individual bits for a byte.
func bits(passed byte) [8]byte {
	return [8]byte{
		(passed >> 7) & 1,
		(passed >> 6) & 1,
		(passed >> 5) & 1,
		(passed >> 4) & 1,
		(passed >> 3) & 1,
		(passed >> 2) & 1,
		(passed >> 1) & 1,
		passed & 1,
	}
}

package lz11

import (
	"bytes"
	"encoding/binary"
	"math"
)

func Compress(data []byte) ([]byte, error){
	// First ensure we can compress the data
	if len(data) <= RawMin {
		return nil, ErrInputTooSmall
	}

	if len(data) > RawMax {
		return nil, ErrInputTooLarge
	}
	
	buffer := bytes.NewBuffer(nil)
	err := binary.Write(buffer, binary.LittleEndian, uint32(0x11|len(data)<<8))
	if err != nil {
		return nil, ErrFailedBufferWrite
	}

	mask := 0
	flag := 0
	index := 0
	_len := 0
	pos := 0
	lenBest := 0
	posBest := 0
	compPos := 4

	for index < len(data) {
		mask >>= BitShiftCount
		if mask == 0 {
			buffer.WriteByte(0)

			flag = compPos
			compPos++
			mask = DefaultMask
		}

		lenBest = MaxNotEncode

		if index >= MaxOffset {
			pos = MaxOffset
		} else {
			pos = index
		}
		for ; pos > VRAMCompatible; pos-- {
			for _len = 0; _len < MaxCoded3; _len++ {
				if index +_len == len(data) {
					break
				}

				if index +_len >= len(data) {
					break
				}

				if data[index +_len] != data[index +_len-pos] {
					break
				}
			}

			if _len > lenBest {
				posBest = pos
				lenBest = _len

				if lenBest == MaxCoded3 {
					break
				}
			}
		}
		if lenBest > MaxNotEncode {
			index += lenBest
			buffer.Bytes()[flag] |= byte(mask)

			if lenBest > MaxCoded2 {
				lenBest -= MaxCoded2 + 1

				cmp := []byte{
					byte((lenBest >> 12) | 16),
					byte((lenBest >> 4) & math.MaxUint8),
					byte(((lenBest & 15) << 4) | (posBest-1)>>8),
					byte((posBest - 1) & math.MaxUint8),
				}

				buffer.Write(cmp)
				compPos += 4
			} else if lenBest > MaxCoded1 {
				lenBest -= MaxCoded1 + 1

				cmp := []byte{
					byte(lenBest >> 4),
					byte(((lenBest & 15) << 4) | (posBest-1)>>8),
					byte((posBest - 1) & math.MaxUint8),
				}

				buffer.Write(cmp)
				compPos += 3
			} else {
				lenBest--
				cmp := []byte{
					byte(((lenBest & 15) << 4) | (posBest-1)>>8),
					byte((posBest - 1) & math.MaxUint8),
				}

				buffer.Write(cmp)
				compPos += 2
			}
		} else {
			buffer.WriteByte(data[index])
			compPos++
			index++
		}
	}

	return buffer.Bytes(), nil
}

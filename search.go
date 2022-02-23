package lz11

import (
	"bytes"
)

func getNeedleTable(input bytes.Buffer) []int {
	needleTable := make([]int, 256)

	for i := 0; i < len(needleTable); i++ {
		needleTable[i] = input.Len()
	}

	for i, c := range input.Bytes() {
		needleTable[c] = input.Len() - i
	}

	return needleTable
}

func searchOne(haystack, needle bytes.Buffer, needleTable []int) *int {
	cur := 0

	for haystack.Len()-cur >= needle.Len() {
		var output *int = nil
		for i := needle.Len() - 1; i != -1; i-- {
			if haystack.Bytes()[cur+i] == needle.Bytes()[i] {
				output = &cur
				break
			}
		}

		if output != nil {
			return output
		}

		cur += needleTable[haystack.Bytes()[cur+needle.Len()-1]]
	}

	return nil
}

func search(haystack, needle []byte) []int {
	var newHaystack *bytes.Buffer
	var newNeedle *bytes.Buffer

	newHaystack = bytes.NewBuffer(haystack)
	newNeedle = bytes.NewBuffer(needle)

	needleTable := getNeedleTable(*newNeedle)
	cur := 0
	var positions []int

	for cur+newNeedle.Len() < newHaystack.Len() {
		var fixedHaystack *bytes.Buffer
		fixedHaystack = bytes.NewBuffer(newHaystack.Bytes()[cur:])
		foundPos := searchOne(*fixedHaystack, *newNeedle, needleTable)

		if foundPos != nil {
			positions = append(positions, *foundPos)
			cur += *foundPos + newNeedle.Len() + 1
		} else {
			return positions
		}
	}

	return positions
}

func findLongestMatch(data bytes.Buffer, off, max int) (*int, *int) {
	if off < 4 || data.Len()-off < 4 {
		return nil, nil
	}

	longestPos := 0
	longestLen := 0
	start := 0

	if off > 0x1000 {
		start = off - 0x1000
	}

	for _, pos := range search(data.Bytes()[start:off+2], data.Bytes()[off:off+3]) {
		length := 0
		loopPos := 0
		for i := off; i < data.Len(); i++ {
			if length == max {
				returnValue := start + pos
				return &returnValue, &length
			}

			if data.Bytes()[i] != data.Bytes()[start+pos+loopPos] {
				break
			}
			length += 1
		}

		if length > longestLen {
			longestPos = pos
			longestLen = length
		}
	}

	if longestLen < 3 {
		return nil, nil
	}

	returnValue := start + longestPos
	return &returnValue, &longestLen
}

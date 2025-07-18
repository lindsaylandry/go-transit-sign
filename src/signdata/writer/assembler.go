package writer

import (
	"fmt"
)

func CreateVisualString(stop string) ([][]uint8, error) {
	l := getLetters()

	// put boolean letters in a horizontal array
	length := 0
	// first pass - make str dimensions
	for _, r := range stop {
		// return error if rune is not in list of ascii letters
		val, ok := l[r]
		if !ok {
			return [][]uint8{}, fmt.Errorf("The letter %c does not exist in pixel library", r)
		}

		length += len(val.Design[0]) + 1
	}

	str := make([][]uint8, len(l['0'].Design))
	for i := range str {
		str[i] = make([]uint8, length)
	}

	startCol := 0
	// second pass - put letters in matrix
	for _, r := range stop {
		for i, a := range l[r].Design {
			for j, b := range a {
				str[i][startCol+j] = b
			}
		}

		startCol += len(l[r].Design[0]) + 1
	}

	return str, nil
}

func CreateVisualNextArrival(dest string, timeLeft string, maxWidth int) ([][]uint8, int, error) {
	destMatrix, err := CreateVisualString(dest)
	if err != nil {
		return [][]uint8{}, 0, err
	}
	timeLeftMatrix, errStr := CreateVisualString(timeLeft)
	if err != nil {
		return [][]uint8{}, 0, errStr
	}

	str := make([][]uint8, len(destMatrix))
	for i := range str {
		str[i] = make([]uint8, maxWidth)
	}

	// Combine dest (left align) and time (right align)
	// TODO: what if strings are longer than max width?
	// first - left align
	for i, a := range destMatrix {
		copy(str[i][:], a)
	}

	// next - right align
	for i, a := range timeLeftMatrix {
		for j, b := range a {
			str[i][len(str[i])-len(timeLeftMatrix[0])+j] = b
		}
	}

	timeIndex := 0
	if len(str) > 0 && len(timeLeftMatrix) > 0 {
		timeIndex = len(str[0]) - len(timeLeftMatrix[0])
	}
	return str, timeIndex, nil
}

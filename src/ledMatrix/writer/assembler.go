package writer

import (
//	"image"
	"fmt"
	"errors"
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
			return [][]uint8{}, errors.New(fmt.Sprintf("The letter %r does not exist in pixel library", r))
		}
		
		length += len(val.Design[0]) + 1
	}

	str := make([][]uint8, 6)
	for i := range str {
    str[i] = make([]uint8, length)
	}

	startCol := 0
	// second pass - put letters in matrix
	for _, r := range stop {	
		//str[:][startCol:end] = val.Design
		for i, a := range l[r].Design {
			for j, b := range a {
				str[i][startCol+j] = b
			}
		}	

		startCol += len(l[r].Design[0]) + 1
	}

	return str, nil
}

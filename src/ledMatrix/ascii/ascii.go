package ascii

type Letter struct {
	Design [][]bool
}

func CreateVisualString(stop string) ([][]bool, error) {

}

func getLetters(map[rune]Letter) {
	letters := make(map[rune]Letter)

	// basic ascii letters and numbers
	letters['a'] = Letter{
		Design: [][]bool {
			{0, 0, 0},
			{0, 0, 0},
			{0, 1, 1},
			{1, 0, 1},
			{0, 1, 1},
			{0, 0, 0},
		},
	}

	letters['b'] = Letter{
    Design: [][]bool {
      {1, 0, 0},
      {1, 0, 0},
      {1, 1, 0},
      {1, 0, 1},
      {1, 1, 0},
      {0, 0, 0},
    },
  }

	letters['c'] = Letter{
    Design: [][]bool {
      {0, 0, 0},
      {0, 0, 0},
      {0, 1, 1},
      {1, 0, 0},
      {0, 1, 1},
      {0, 0, 0},
    },
  }

	letters['d'] = Letter{
    Design: [][]bool {
      {0, 0, 1},
      {0, 0, 1},
      {0, 1, 1},
      {1, 0, 1},
      {0, 1, 1},
      {0, 0, 0},
    },
  }

	return letters
}

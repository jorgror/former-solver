package game

import (
	"bufio"
	"os"
)

func (g *Grid) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	y := height - 1
	for scanner.Scan() {
		if y < 0 {
			break
		}
		line := scanner.Text()
		for x := 0; x < width && x < len(line); x++ {
			var color Color
			switch line[x] {
			case 'O':
				color = Orange
			case 'P':
				color = Pink
			case 'B':
				color = Blue
			case 'G':
				color = Green
			default:
				color = Empty
			}
			g.cells[y][x] = Cell{color: color}
		}
		y--
	}
	return scanner.Err()
}

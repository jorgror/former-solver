package game

import (
	"testing"
	// ...existing code...
)

func TestClone(t *testing.T) {
	originalGrid := NewGrid()
	clonedGrid := originalGrid.Clone()

	// Verify that the cloned grid has the same cells as the original
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if originalGrid.cells[y][x].color != clonedGrid.cells[y][x].color {
				t.Errorf("Cell mismatch at (%d, %d): original %v, clone %v", x, y, originalGrid.cells[y][x].color, clonedGrid.cells[y][x].color)
			}
		}
	}

	// Modify the cloned grid
	clonedGrid.ClearRegion(0, 0)

	// Verify that the original grid is not modified
	if originalGrid.cells[8][0].color == Empty {
		t.Errorf("Original grid was modified")
	}

	// Verify that the cloned grid is modified
	if clonedGrid.cells[8][0].color != Empty {
		t.Errorf("Cloned grid was not modified")
	}

}

func TestCloneAfterClear(t *testing.T) {
	originalGrid := NewGrid()
	originalGrid.ClearRegion(0, 0)
	clonedGrid := originalGrid.Clone()

	// Verify that the cloned grid has the same cells as the original
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if originalGrid.cells[y][x].color != clonedGrid.cells[y][x].color {
				t.Errorf("Cell mismatch at (%d, %d): original %v, clone %v", x, y, originalGrid.cells[y][x].color, clonedGrid.cells[y][x].color)
			}
		}
	}

	// Compare the actions
	if len(originalGrid.actions) != len(clonedGrid.actions) {
		t.Errorf("Action count mismatch: original %d, clone %d", len(originalGrid.actions), len(clonedGrid.actions))
	}
}

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const boardSize = 9
const subgridSize = 3
const defaultEmptyCells = 40

type SudokuBoard struct {
	Board   [boardSize][boardSize]int
	RowNums [boardSize]map[int]bool
	ColNums [boardSize]map[int]bool
}

func (s *SudokuBoard) generateHTMLTable() string {
	var htmlTable strings.Builder

	htmlTable.WriteString("<style>")
	htmlTable.WriteString("table { border-collapse: collapse; }\n")
	htmlTable.WriteString("td { border: 1px solid black; width: 20px; height: 20px; text-align: center; }\n")
	htmlTable.WriteString("td.empty { background-color: lightgray; }\n")
	htmlTable.WriteString("td.subgrid-separator { border: none; width: 10px; }\n")
	htmlTable.WriteString("</style>\n")

	htmlTable.WriteString("<table>\n")
	for i := 0; i < boardSize; i++ {
		htmlTable.WriteString("<tr>")
		for j := 0; j < boardSize; j++ {
			num := s.Board[i][j]
			cellClass := "empty"
			if num != 0 {
				cellClass = ""
			}
			htmlTable.WriteString("<td class=\"" + cellClass + "\">")
			if num != 0 {
				htmlTable.WriteString(fmt.Sprintf("%d", num))
			}
			htmlTable.WriteString("</td>")
			if (j+1)%subgridSize == 0 && j < boardSize-1 {
				htmlTable.WriteString("<td class=\"subgrid-separator\"></td>")
			}
		}
		htmlTable.WriteString("</tr>")
		if (i+1)%subgridSize == 0 && i < boardSize-1 {
			htmlTable.WriteString("<tr><td colspan=\"" + fmt.Sprintf("%d", boardSize+subgridSize-1) + "\" class=\"subgrid-separator\"></td></tr>")
		}
	}
	htmlTable.WriteString("</table>\n")

	return htmlTable.String()
}

func (s *SudokuBoard) initialize() {
	for i := 0; i < boardSize; i++ {
		s.RowNums[i] = make(map[int]bool)
		s.ColNums[i] = make(map[int]bool)
		for j := 1; j <= boardSize; j++ {
			s.RowNums[i][j] = true
			s.ColNums[i][j] = true
		}
	}
}

func (s *SudokuBoard) generateSolvedSudoku() bool {
	return s.solve(0, 0)
}

func (s *SudokuBoard) solve(row, col int) bool {
	if row == boardSize {
		row = 0
		col++
		if col == boardSize {
			return true
		}
	}

	availableNums := s.intersection(s.RowNums[row], s.ColNums[col])
	s.shuffle(availableNums)

	for _, num := range availableNums {
		if s.isValid(num, row, col) {
			s.Board[row][col] = num
			delete(s.RowNums[row], num)
			delete(s.ColNums[col], num)

			if s.solve(row+1, col) {
				return true
			}

			s.Board[row][col] = 0
			s.RowNums[row][num] = true
			s.ColNums[col][num] = true
		}
	}

	return false
}

func (s *SudokuBoard) isValid(num, row, col int) bool {
	subgridRow := row / subgridSize
	subgridCol := col / subgridSize
	subgridStartRow := subgridRow * subgridSize
	subgridStartCol := subgridCol * subgridSize

	for i := 0; i < subgridSize; i++ {
		for j := 0; j < subgridSize; j++ {
			if s.Board[subgridStartRow+i][subgridStartCol+j] == num {
				return false
			}
		}
	}

	return s.RowNums[row][num] && s.ColNums[col][num]
}

func (s *SudokuBoard) intersection(slice1, slice2 map[int]bool) []int {
	result := []int{}
	for num := 1; num <= boardSize; num++ {
		if slice1[num] && slice2[num] {
			result = append(result, num)
		}
	}
	return result
}

func (s *SudokuBoard) shuffle(nums []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
}

func (s *SudokuBoard) generatePlayableSudoku(clearedCells int) {
	rand.Seed(time.Now().UnixNano())

	for count := 0; count < clearedCells; {
		i := rand.Intn(boardSize)
		j := rand.Intn(boardSize)
		if s.Board[i][j] != 0 {
			s.Board[i][j] = 0
			count++
		}
	}
}

func setCellsToClear() int {
	// Read number of cells to clear from command-line argument
	if len(os.Args) == 2 {
		cellsToClear, err := strconv.Atoi(os.Args[1])
		if err != nil || cellsToClear < 0 || cellsToClear > boardSize*boardSize {
			return defaultEmptyCells
		}
		return cellsToClear
	}
	return defaultEmptyCells
}

func main() {
	// Create a slice to store generated Sudoku puzzles
	var sudokuPuzzles []SudokuBoard

	// Generate Sudoku puzzles and store them in the slice
	clearedCells := setCellsToClear()
	for i := 0; i < 3; i++ {
		var sudokuPuzzle SudokuBoard
		sudokuPuzzle.initialize()
		sudokuPuzzle.generateSolvedSudoku()
		sudokuPuzzles = append(sudokuPuzzles, sudokuPuzzle)
		sudokuPuzzle.generatePlayableSudoku(clearedCells)
		sudokuPuzzles = append(sudokuPuzzles, sudokuPuzzle)
	}

	htmlContent := generateHTML(sudokuPuzzles)

	// Write HTML content to the file
	err := writeToFile("sudoku.html", htmlContent)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("HTML content written to sudoku.html")
}

func generateHTML(sudokuPuzzles []SudokuBoard) string {
	htmlContent := "<html><body>"

	// Generate three rows of two puzzles each
	for row := 0; row < 3; row++ {
		htmlContent += "<div style=\"display: flex;\">"
		for col := 0; col < 2; col++ {
			htmlContent += "<div style=\"flex: 1;\">"
			htmlContent += "<hr/><br/>"
			htmlContent += sudokuPuzzles[row*2+col].generateHTMLTable()
			htmlContent += "</div>"
		}
		htmlContent += "</div><br/>"
	}

	htmlContent += "</body></html>"
	return htmlContent
}

func writeToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

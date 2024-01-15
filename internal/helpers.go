package helpers

import (
	"fmt"
	"strings"
)

func PrintTable(headers []string, data [][]string) {
	for _, header := range headers {
		fmt.Printf("\t%-10s", header)
	}
	fmt.Println("\n" + strings.Repeat("-", 80))
	for _, row := range data {
		for _, col := range row {
			fmt.Printf("\t%-10s", col)
		}
		fmt.Println()
	}
}
func FormatNumber(n int64) string {
	numStr := fmt.Sprintf("%d", n)
	var result strings.Builder
	length := len(numStr)
	for i := 0; i < length; i++ {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteString(" ")
		}
		result.WriteByte(numStr[i])
	}
	return result.String()
}

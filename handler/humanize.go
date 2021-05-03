package handler

import (
	"fmt"
	"math"
)

func HumanizeNumber(views int) string {
	millNames := []string{"", "k", "M", "Billion", "Trillion"}

	thousands := math.Min(float64(len(millNames)-1), math.Floor(math.Log10(math.Abs(float64(views)))/3.0))
	millidx := math.Max(0, thousands)

	return fmt.Sprintf("%.0f%s", float64(views)/math.Pow10(3*int(millidx)), millNames[int(millidx)])
}

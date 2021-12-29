package main

import "fmt"

type Time struct {
	days int
	hours int
	minutes int
}

func minutesToCost(min int) string {
	return fmt.Sprintf("%d руб.", int(float32(min) / 60 * COST))
}

func divMod(a int, b int) int {
	return (a - (a % b)) / b
}

func parseMinutes(min int) Time {
	var t Time
	t.days = divMod(min, 480)
	t.hours = divMod((min - t.days * 480), 60)
	t.minutes = min - (t.days * 480) - (t.hours * 60)
	return t
}

func minutesToString(min int) string {
	t := parseMinutes(min)
	var result string
	if t.days != 0 {
		result += fmt.Sprintf("%dд ", t.days)
	}
	if t.hours != 0 {
		result += fmt.Sprintf("%dч ", t.hours)
	}
	if t.minutes != 0 {
		result += fmt.Sprintf("%dм ", t.minutes)
	}

	return result
}
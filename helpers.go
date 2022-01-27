package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

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

func getMonthTranslate(date string) string {
	var str string
	switch strings.Split(date, "-")[1] {
	case "01":
		str = "янв."
	case "02":
		str = "фев."
	case "03":
		str = "мар."
	case "04":
		str = "апр."
	case "05":
		str = "май."
	case "06":
		str = "июнь."
	case "07":
		str = "июль."
	case "08":
		str = "авг."
	case "09":
		str = "сент."
	case "10":
		str = "окт."
	case "11":
		str = "нояб."
	case "12":
		str = "дек."
	}
	return str
}

func contains(s []string, searchterm string) bool {
    i := sort.SearchStrings(s, searchterm)
    return i < len(s) && s[i] == searchterm
}


func createDir(path string) {
	_, err := os.Stat(path)
    if err != nil {
        e := os.Mkdir(path, os.ModeDir)
		checkError(e, "Ошибка при создании директории")
    }
}

func createFile(fileName string) {
	f, err := os.Create(fileName)
	checkError(err, "Ошибка при создании файла")
	defer f.Close()
}
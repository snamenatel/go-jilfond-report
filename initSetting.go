package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
)

func getDateReport() time.Time {
	var currentMonthScrool int
	var result time.Time
	if len(os.Args) < 2 {
		prompt := promptui.Select{
			HideHelp: false,
			Label: "Select report year",
			Items: getSelectYearItems(),
		}
		_, year, err := prompt.Run()
		checkError(err, "Prompt failed")

		currentMonth, _ := strconv.Atoi(time.Now().Format("01"))
		currentMonth--
		if currentMonth > 5 {
			currentMonthScrool = currentMonth
		} else {
			currentMonthScrool = 0
		}
		prompt = promptui.Select{
			HideHelp: false,
			Label: "Select report month",
			Items: getSelectMonthItems(),
		}
		_, month, err := prompt.RunCursorAt(currentMonth, currentMonthScrool)
		checkError(err, "Prompt failed")

		t, _ := time.Parse("2006-January-02", fmt.Sprintf("%s-%s-01", year, month))
		result = t
	} else {
		t, _ := time.Parse("2006-01", os.Args[1])
		result = t
	}
	fmt.Printf("Поиск отчета за период %s \n", result)
	return result;
}

func getSelectYearItems() []string {
	var res []string
	for i := 0; i < 3; i++ {
		res = append(res, time.Now().AddDate(-i, 0, 0).Format("2006"))
	}
	return res; 
}

func getSelectMonthItems() []string {
	var res []string
	for i := 1; i <= 12; i++ {
		res = append(res, time.Date(2021, time.Month(i), 1, 0, 0, 0, 0, time.UTC).Format("January"))
	}
	return res; 
}

func isNeedFutureTasks() bool {
	prompt := promptui.Select{
		HideSelected: true,
		HideHelp: false,
		Label: "Вы хотите получить список заданий будущего спринта",
		Items: []string{"Да", "Нет"},
	}
	_, answer, err := prompt.Run()
	checkError(err, "Prompt failed")
	return answer == "Да"
}
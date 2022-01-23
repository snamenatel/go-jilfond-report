package main

import (
	"fmt"
	"os"
	"time"

	"github.com/manifoldco/promptui"
)

func getDateReport() string {
	if len(os.Args) < 2 {
		prompt := promptui.Select{
			Label: "Select report year",
			Items: getSelectYearItems(),
		}
		_, year, err := prompt.Run()
		checkError(err, "Prompt failed")

		prompt = promptui.Select{
			Label: "Select report month",
			Items: getSelectMonthItems(),
		}
		_, month, err := prompt.Run()
		checkError(err, "Prompt failed")

		t, _ := time.Parse("2006-January-02", fmt.Sprintf("%s-%s-01", year, month))
		return t.Format("2006-01")
	}
	return os.Args[1]
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
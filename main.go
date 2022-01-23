package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/joho/godotenv"
)

const LINE_LENGHT = 110
const COST = 700
var url string
var token string
var reportDate string

func init() {
	if err := godotenv.Load(); err != nil {
        fmt.Print("No .env file found")
    }
	token = "Bearer " + os.Getenv("TOKEN")
	url = os.Getenv("URL")

}

func main() {
	reportDate = getDateReport()

	reportIds := make(map[string]string)
	for _, report := range GetReportsList() {
		if report.Name == "Внедрено " + reportDate {
			reportIds["Implemented"] = report.Id
		}
		if report.Name == "Отработано " + reportDate {
			reportIds["Spent"] = report.Id
		}
	}
	if len(reportIds) == 0 {
		fmt.Println("Отчеты за выбранный период не найдены")
		return
	}

	fmt.Printf("Найдены отчеты: \n - по внедренным задачам - %s/reports/time/%s \n - по отработанному времени - %s/reports/time/%s\n\r", url, reportIds["Implemented"], url, reportIds["Spent"])	
	reports := make(map[string]Report)
	for k, v := range reportIds {
		reports[k] = GetReport(v)
	}
	var repotContent []string
	for _, report := range reports {
		for _, group := range report.Data.Groups {
			if group.Meta.LinkedUser.VisibleName == "Дударек Илья" {
				if strings.HasPrefix(report.Name, "Внедрено") {
					repotContent = append(repotContent, formatCompleatedReport(group))
				} else {
					repotContent = append(repotContent, formatSpentReport(group))
				}
			}
		}
	}

	writeToFile(repotContent)
}


func checkError(err error, text string) {
	if err != nil {
		panic(text)
	}
}

func formatSpentReport(group ReportGroupItem) string {
	var rowList []string
	for _, line := range group.Lines {
		title := line.IssueId + ": " + line.Description
		duration := minutesToString(line.TotalDuration.Value)
		if LINE_LENGHT - utf8.RuneCountInString(title) - utf8.RuneCountInString(duration) <= 0 {
			title = title[: LINE_LENGHT - utf8.RuneCountInString(duration) - 4] + "..."
		}

		dashLine := strings.Repeat(".", LINE_LENGHT - utf8.RuneCountInString(title) - utf8.RuneCountInString(duration))
		rowList = append(rowList, fmt.Sprintf("%s%s%s", title, dashLine, duration))

	}

	return fmt.Sprintf("<b>Отработано в %s %s</b>\n<pre>\n%s \n</pre>",
		getMonthTranslate(reportDate),
		strings.Split(reportDate, "-")[0],
		strings.Join(rowList, "\n"))
}

func formatCompleatedReport(group ReportGroupItem) string {
	var rowList []string
	lengthLeft := 80
	lengthRight := LINE_LENGHT - lengthLeft
	sumMinutes := 0
	for _, line := range group.Lines {
		sumMinutes += line.Estimation.Value
		title := line.IssueId + ": " + line.Description
		estimate := minutesToString(line.Estimation.Value)
		cost := minutesToCost(line.Estimation.Value)
		
		if lengthLeft - utf8.RuneCountInString(title) <= 0 {
			title = title[: lengthLeft - utf8.RuneCountInString(title) - 4] + "..."
		}
		dashLineLeft := strings.Repeat(".", lengthLeft - utf8.RuneCountInString(title))
		dashLineRight := strings.Repeat(".", lengthRight - utf8.RuneCountInString(estimate) - utf8.RuneCountInString(cost))
		rowList = append(rowList, fmt.Sprintf("%s%s%s%s%s", title, dashLineLeft, estimate, dashLineRight, cost))
	}
	totalMin := minutesToString(group.Estimation.Value)
	totalCost := minutesToCost(group.Estimation.Value)
	dashLineLeft := strings.Repeat(".", lengthLeft - utf8.RuneCountInString("Итого:"))
	dashLineRight := strings.Repeat(".", lengthRight - utf8.RuneCountInString(totalMin) - utf8.RuneCountInString(totalCost))
	rowList = append(rowList, fmt.Sprintf("Итого:%s%s%s%s",dashLineLeft, totalMin, dashLineRight, totalCost))

	return fmt.Sprintf("<b>Внедрено в %s %s</b>\n<pre>\n%s \n</pre>",
		getMonthTranslate(reportDate),
		strings.Split(reportDate, "-")[0],
		strings.Join(rowList, "\n"))
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

func writeToFile(contentList []string) {
	createDir("./reports")
	fileName := fmt.Sprintf("./reports/%s.txt", reportDate)
	createFile(fileName)
	
	err := os.WriteFile(fileName, []byte(strings.Join(contentList, "\n\r\n\r")), 0644)
	filePath, _ := filepath.Abs(fileName)
	fmt.Printf("Отчеты сформированы %s\n", filePath)

	checkError(err, "Ошибка при записи файла")
}


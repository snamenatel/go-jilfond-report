package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
	reportDate = getDateReport()
	fmt.Println("Поиск отчета за период", reportDate)
}

func main() {
	var curentReportId string
	for _, report := range GetReportsList() {
		if report.Name == "Внедрено " + reportDate {
			curentReportId = report.Id
		}
	}
	
	fmt.Println("Найден отчет " + curentReportId)
	GetReport(curentReportId)
}


func checkError(err error, text string) {
	if err != nil {
		panic(text)
	}
}

func formatReport(group ReportGroupItem) string {
	fmt.Printf("Отчет для %s \n", group.Meta.LinkedUser.VisibleName)
	var rowList []string
	for _, line := range group.Lines {
		title := line.IssueId + ": " + line.Description
		duration := strconv.Itoa(line.TotalDuration.Value)
		if LINE_LENGHT - utf8.RuneCountInString(title) - utf8.RuneCountInString(duration) <= 0 {
			title = title[: LINE_LENGHT - utf8.RuneCountInString(duration) - 4] + "..."
		}

		dashLine := strings.Repeat("_", LINE_LENGHT - utf8.RuneCountInString(title) - utf8.RuneCountInString(duration))
		rowList = append(rowList, fmt.Sprintf("%s%s%s", title, dashLine, duration))

	}
	return strings.Join(rowList, "\n")
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
	
	err := os.WriteFile(fileName, []byte(strings.Join(contentList, "\n")), 0644)
	checkError(err, "Ошибка при записи файла")
}

func getDateReport() string {
	if len(os.Args) < 2 {
		return time.Now().Format("2006-01")
	}
	return os.Args[1]
}

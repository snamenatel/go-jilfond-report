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

type ReportListItem struct {
	Own bool	`json:"own"`
	Name string	`json:"name"`
	Id string	`json:"id"`
	Type string	`json:"$type"`
}

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

func createDir(path string) {
	_, err := os.Stat(path)
    if err != nil {
        e := os.Mkdir(path, os.ModeDir)
		if e != nil {
			fmt.Println("Ошибка при создании директории: ", err)
			os.Exit(0)
		}
    }
}

func createFile(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Ошибка при создании файла: ", err)
		os.Exit(0)
	}
	defer f.Close()
}

func writeToFile(contentList []string) {
	createDir("./reports")
	fileName := fmt.Sprintf("./reports/%s.txt", reportDate)
	createFile(fileName)
	
	err := os.WriteFile(fileName, []byte(strings.Join(contentList, "\n")), 0644)
	if err != nil {
		fmt.Println("Ошибка при записи файла: ", err)
		os.Exit(0)
	}
}

func getDateReport() string {
	if len(os.Args) < 2 {
		return time.Now().Format("2006-01")
	}
	return os.Args[1]
}

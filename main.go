package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func writeToFile(contentList []string) {
	createDir("./reports")
	fileName := fmt.Sprintf("./reports/%s.txt", reportDate)
	createFile(fileName)
	
	err := os.WriteFile(fileName, []byte(strings.Join(contentList, "\n\r\n\r")), 0644)
	filePath, _ := filepath.Abs(fileName)
	fmt.Printf("Отчеты сформированы %s\n", filePath)

	checkError(err, "Ошибка при записи файла")
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

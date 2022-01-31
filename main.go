package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/joho/godotenv"
)

const LINE_LENGHT = 110
var costPerHour int
var url string
var token string
var reportDate time.Time

func init() {
	if err := godotenv.Load(); err != nil {
        fmt.Print("No .env file found")
    }
	token = "Bearer " + os.Getenv("TOKEN")
	url = os.Getenv("URL")
	t, _ := strconv.Atoi(os.Getenv("COST"))
	costPerHour = t
}

func main() {
	fmt.Println(costPerHour, )
	reportDate = getDateReport()

	reportIds := make(map[string]string)
	for _, report := range GetReportsList() {
		if report.Name == "Внедрено " + reportDate.Format("2006-01") {
			reportIds["Implemented"] = report.Id
		}
		if report.Name == "Отработано " + reportDate.Format("2006-01") {
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

	if isNeedFutureTasks() {
		dateNextSprint := reportDate.AddDate(0, 1, 0).Format("06-01")
		fmt.Printf("Поиск задач будет осуществляться в спринте %s \n", dateNextSprint)
		taskList := GetTaskList(GetTaskIDList(GetCurrentSprintID(dateNextSprint)))
		
		planTasks := []string{}
		prompt := &survey.MultiSelect{
			Message: "Выберите задачи для добавления их в план на следующий месяц",
			Options: taskList,
			Help: " ",
		}
		survey.AskOne(prompt, &planTasks)
		
		priorityTasks := []string{}
		prompt = &survey.MultiSelect{
			Message: "Выберите приоритетные задачи",
			Options: planTasks,
		}
		survey.AskOne(prompt, &priorityTasks)
		sort.Strings(planTasks)
		sort.Strings(priorityTasks)
		appendToFile(planTasks, priorityTasks)
	}
}

func writeToFile(contentList []string) {
	createDir("./reports")
	fileName := fmt.Sprintf("./reports/%s.txt", reportDate.Format("2006-01"))
	createFile(fileName)
	
	err := os.WriteFile(fileName, []byte(strings.Join(contentList, "\n\r\n\r")), 0644)
	filePath, _ := filepath.Abs(fileName)
	fmt.Printf("Отчеты сформированы %s\n", filePath)

	checkError(err, "Ошибка при записи файла")
}

func appendToFile(planTasks, priorityTasks []string) {
	plan, priority := futureTaskFormat(filter(planTasks, priorityTasks), priorityTasks)
	fileName := fmt.Sprintf("./reports/%s.txt", reportDate.Format("2006-01"))

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	checkError(err, "Ошибка при записи файла")
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("\n\r\n\r%s\n\r\n\r%s\n\r", plan, priority))
	checkError(err, "Ошибка при записи файла")
}

func filter(origin, target []string) []string {
	filtered := []string{}
	for _, str := range origin {
		if !contains(target, str) {
			filtered = append(filtered, str)
		}
	}
	return filtered
}

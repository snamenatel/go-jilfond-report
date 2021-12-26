package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

type TotalDuration struct {
	Value int
}

type ReportLine struct {
	Description string
	IssueId string
	TotalDuration TotalDuration
}

type LinkedUser struct {
	VisibleName string
	RingId string
}
type ReportGroupMeta struct {
	LinkedUser LinkedUser
}
type ReportGroupItem struct {
	Lines []ReportLine
	Meta ReportGroupMeta
}
type ReportData struct {
	Groups []ReportGroupItem
}
type Report struct {
	Data ReportData
	Name string
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
	for _, report := range getReportsList() {
		if report.Name == "Внедрено " + reportDate {
			curentReportId = report.Id
		}
	}
	
	fmt.Println("Найден отчет " + curentReportId)
	getReport(curentReportId)
}

func getReportsList()[]ReportListItem {
	req, err := http.NewRequest("GET", url + "/api/reports?$top=-1&fields=$type,id,name,own", nil)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при получении списка отчетов" + "\n")
    }
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error while reading the response bytes:", err)
    }

	var reporstList []ReportListItem
	json.Unmarshal(body, &reporstList)

	return reporstList
}

func getReport(id string) {
	currentReportUrl := url + "/api/reports/" + id + "?$top=-1&fields=$type,aggregationPolicy($type,field($type,id,presentation)),attachedFields($type,field($type,aggregateable,id,presentation)),authors(fullName,id,login,ringId),bubbleWorkItems,data(duration(id,value),estimation(id,value),fieldNames,groups(duration(id,value),estimation(id,value),id,issueId,issueSummary,lines(comment,date,description,duration(id,value),estimation(id,value),fieldValues,groupName,id,issueId,issueSummary,name,totalDuration(id,value),typeDurations(duration(value),workType),userAvatarUrl,userId,userVisibleName),meta(linkedIssue(idReadable,summary),linkedUser(ringId,visibleName,postfix)),name,totalDuration(id,value),typeDurations(duration(value),workType)),totalDuration(id,value),typeDurations(duration(value),workType)),effectiveQuery,effectiveQueryUrl,grouping(id,field(id,presentation)),id,invalidationInterval,name,own,owner(id,login),pin(pinned),projects(id,name,ringId,shortName),query,range($type,from,range(id),to),sprint(id,name),status(calculationInProgress,error(id),errorMessage,isOutdated,lastCalculated,progress,wikifiedErrorMessage),updateableBy(id,name),visibleTo(id,name),wikifiedError,windowSize,workTypes(id,name)&line=issue&query=for:ringId"
	req, err := http.NewRequest("GET", currentReportUrl, nil)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при получении отчета" + "\n")
    }
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error while reading the response bytes:", err)
    }

	var report Report
	json.Unmarshal(body, &report)

	var repotContent []string
	for _, group := range report.Data.Groups {
		if group.Meta.LinkedUser.VisibleName == "Дударек Илья" {
			repotContent = append(repotContent, formatReport(group))
			break
		}
	}

	writeToFile(repotContent)
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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type ReportListItem struct {
	Own bool	`json:"own"`
	Name string	`json:"name"`
	Id string	`json:"id"`
	Type string	`json:"$type"`
}

type TotalDuration struct {
	Value int32
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

	for _, group := range report.Data.Groups {
		if group.Meta.LinkedUser.VisibleName == "Дударек Илья" {
			formatReport(group)
			break
		}
	}
}

func formatReport(group ReportGroupItem) {
	fmt.Printf("Отчет для %s \n", group.Meta.LinkedUser.VisibleName)
	for _, line := range group.Lines {
		fmt.Printf("Задача: %s %s ----------------- %d \n", line.IssueId, line.Description, line.TotalDuration.Value)
	}
}

func getDateReport() string {
	if len(os.Args) < 2 {
		return time.Now().Format("2006-01")
	}
	return os.Args[1]
}

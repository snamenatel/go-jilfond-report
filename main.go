package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ReportListItem struct {
	Own bool `json: "own"`
	Name string `json: "name"`
	Id string `json: "id"`
	Type string `json: "$type"`
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
	reportDate = os.Args[1]
}

func main() {
	fmt.Println(url)
	fmt.Println(token)
	fmt.Println(reportDate)

	var curentReportId string
	for _, report := range getReportsList() {
		if report.Name == "Внедрено " + reportDate {
			curentReportId = report.Id
		}
	}
	
	fmt.Println(curentReportId)
}

func getReportsList()[]ReportListItem {
	req, err := http.NewRequest("GET", url + "/api/reports?$top=-1&fields=$type,id,name,own", nil)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("Ошибка при получении списка отчетов")
    }
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Print("Error while reading the response bytes:", err)
    }

	var reporstList []ReportListItem
	json.Unmarshal(body, &reporstList)

	return reporstList
}

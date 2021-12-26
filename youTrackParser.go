package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ReportListItem struct {
	Own bool	`json:"own"`
	Name string	`json:"name"`
	Id string	`json:"id"`
	Type string	`json:"$type"`
}

type ReportLine struct {
	Description   string
	IssueId       string
	TotalDuration struct {
		Value int
	}
}

type ReportGroupMeta struct {
	LinkedUser struct {
		VisibleName string
		RingId      string
	}
}
type ReportGroupItem struct {
	Lines []ReportLine
	Meta  ReportGroupMeta
}

type Report struct {
	Data struct {
		Groups []ReportGroupItem
	}
	Name string
}

func GetReportsList() []ReportListItem {
	req, err := http.NewRequest("GET", url+"/api/reports?$top=-1&fields=$type,id,name,own", nil)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err, "Ошибка при получении списка отчетов")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err, "Error while reading the response bytes")
	var reporstList []ReportListItem
	json.Unmarshal(body, &reporstList)

	return reporstList
}

func GetReport(id string) {
	currentReportUrl := url + "/api/reports/" + id + "?$top=-1&fields=$type,aggregationPolicy($type,field($type,id,presentation)),attachedFields($type,field($type,aggregateable,id,presentation)),authors(fullName,id,login,ringId),bubbleWorkItems,data(duration(id,value),estimation(id,value),fieldNames,groups(duration(id,value),estimation(id,value),id,issueId,issueSummary,lines(comment,date,description,duration(id,value),estimation(id,value),fieldValues,groupName,id,issueId,issueSummary,name,totalDuration(id,value),typeDurations(duration(value),workType),userAvatarUrl,userId,userVisibleName),meta(linkedIssue(idReadable,summary),linkedUser(ringId,visibleName,postfix)),name,totalDuration(id,value),typeDurations(duration(value),workType)),totalDuration(id,value),typeDurations(duration(value),workType)),effectiveQuery,effectiveQueryUrl,grouping(id,field(id,presentation)),id,invalidationInterval,name,own,owner(id,login),pin(pinned),projects(id,name,ringId,shortName),query,range($type,from,range(id),to),sprint(id,name),status(calculationInProgress,error(id),errorMessage,isOutdated,lastCalculated,progress,wikifiedErrorMessage),updateableBy(id,name),visibleTo(id,name),wikifiedError,windowSize,workTypes(id,name)&line=issue&query=for:ringId"
	req, err := http.NewRequest("GET", currentReportUrl, nil)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err, "Ошибка при получении отчета")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err, "Error while reading the response bytes")

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
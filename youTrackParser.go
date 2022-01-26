package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
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
	Estimation struct {
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
	Estimation struct {
		Value int
	}
}

type Report struct {
	Data struct {
		Groups []ReportGroupItem
	}
	Name string
}

type Sprint struct {
	Archived bool
	Id string
	Name string
}

type BoardCell struct {
	Issues []Issue
}

type Issue struct {
	Id string		`json:"id"`
	Summary string	`json:"summary"`
}

type Task struct {
	Id string
	IdReadable string
}

type TaskInfo struct {
	Id string
	IdReadable string
	Summary string
}

func GetReportsList() []ReportListItem {
	req, err := http.NewRequest("GET", url+"/api/reports?$top=-1&fields=$type,id,name,own", nil)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err, "Ошибка при получении списка отчетов. Проверьте доступронсть " + url)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err, "Error while reading the response bytes")
	var reporstList []ReportListItem
	json.Unmarshal(body, &reporstList)

	return reporstList
}

func GetReport(id string) Report {
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
	return report
}

func GetCurrentSprintID(date string) string {
	sprintListUrl := url + "/api/agiles/114-4/sprints/?issuesQuery=&$top=-1&fields=archived,finish,goal,id,isDefault,name,report(id),start"
	req, err := http.NewRequest("GET", sprintListUrl, nil)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err, "Ошибка при получении отчета")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err, "Error while reading the response bytes")

	var sprintList []Sprint
	json.Unmarshal(body, &sprintList)

	var sprintId string
	for _, sprint := range sprintList {
		if strings.HasPrefix(sprint.Name, date) {
			sprintId = sprint.Id
			break
		}
	}
	return sprintId
}

func GetTaskIDList(spritnId string) []Issue {
	currentUrl := url + "/api/agiles/114-4/sprints/" + spritnId + "?issuesQuery=%23%D1%8F%20&$top=-1&$topLinks=3&$topSwimlanes=10&fields=agile(hideOrphansSwimlane,id,name,orphansAtTheTop,status(valid,warnings,errors)),archived,board(columns(agileColumn(collapsed,color($type,id),fieldValues(canUpdate,column(id),id,isResolved,name,ordinal,presentation),id,isResolved,isVisible,ordinal,ordinal,parent($type,id),wipLimit($type,min,max)),collapsed,id,timeTrackingData(effectiveEstimation,estimation,hasExplicitSpentTime,spentTime)),id,name,notOnBoardCount,orphanRow($type,cells(column(collapsed,id),id,issues($type,id,summary),issuesCount,row(id),tooManyIssues),collapsed,id,issue($type,created,created,description,fields($type,hasStateMachine,id,isUpdatable,name,projectCustomField($type,bundle(id),canBeEmpty,emptyFieldText,field(fieldType(isMultiValue,valueType),id,localizedName,name,ordinal),id,isEstimation,isPublic,isSpentTime,ordinal,size),value($type,archived,avatarUrl,buildLink,color(id),fullName,id,isResolved,localizedName,login,markdownText,minutes,name,presentation,ringId,text)),id,idReadable,isDraft,numberInProject,project($type,id,name,plugins(timeTrackingSettings(enabled,estimate(field(id,name),id),timeSpent(field(id,name),id)),vcsIntegrationSettings(processors(enabled,url,upsourceHubResourceKey,server(enabled,url)))),ringId,shortName),reporter($type,avatarUrl,email,fullName,id,isLocked,issueRelatedGroup(icon),login,name,online,profiles(general(trackOnlineStatus)),ringId),resolved,resolved,summary,updated,watchers(hasStar),wikifiedDescription),matchesQuery,name,timeTrackingData(effectiveEstimation,estimation,hasExplicitSpentTime,spentTime),value(presentation)),sortByQuery,sprint,timeTrackingData(effectiveEstimation,estimation,hasExplicitSpentTime,spentTime),trimmedSwimlanes($type,cells(column(collapsed,id),id,issues($type,id,summary),issuesCount,row(id),tooManyIssues),collapsed,id,issue($type,created,created,description,fields($type,hasStateMachine,id,isUpdatable,name,projectCustomField($type,bundle(id),canBeEmpty,emptyFieldText,field(fieldType(isMultiValue,valueType),id,localizedName,name,ordinal),id,isEstimation,isPublic,isSpentTime,ordinal,size),value($type,archived,avatarUrl,buildLink,color(id),fullName,id,isResolved,localizedName,login,markdownText,minutes,name,presentation,ringId,text)),id,idReadable,isDraft,numberInProject,project($type,id,name,plugins(timeTrackingSettings(enabled,estimate(field(id,name),id),timeSpent(field(id,name),id)),vcsIntegrationSettings(processors(enabled,url,upsourceHubResourceKey,server(enabled,url)))),ringId,shortName),reporter($type,avatarUrl,email,fullName,id,isLocked,issueRelatedGroup(icon),login,name,online,profiles(general(trackOnlineStatus)),ringId),resolved,resolved,summary,updated,watchers(hasStar),wikifiedDescription),matchesQuery,name,timeTrackingData(effectiveEstimation,estimation,hasExplicitSpentTime,spentTime),value(presentation))),eventSourceTicket,finish,goal,id,isDefault,name,report(id),start"
	req, err := http.NewRequest("GET", currentUrl, nil)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err, "Ошибка при получении отчета")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err, "Error while reading the response bytes")

	var tempObj struct {
		Board struct {
			OrphanRow struct {
				Cells [] BoardCell 
			} 
		}
	}
	json.Unmarshal(body, &tempObj)
	var list []Issue
	for _, cell := range tempObj.Board.OrphanRow.Cells {
		list = append(list, cell.Issues...)
	}

	return list
}

func GetTaskList(taskIds []Issue) []string {
	temp := make([]struct{Id string `json:"id"`}, len(taskIds))
	for idx, t := range taskIds {
		temp[idx] = struct{Id string `json:"id"`}{Id: t.Id}
	}
	jsonData, err := json.Marshal(temp)
	checkError(err, "Произошла ошибка")

	currentUrl := url + "/api/issuesGetter?$top=-1&$topLinks=3&fields=attachments(id),fields($type,hasStateMachine,id,isUpdatable,name,projectCustomField($type,bundle(id),canBeEmpty,emptyFieldText,field(fieldType(isMultiValue,valueType),id,localizedName,name,ordinal),id,isEstimation,isPublic,isSpentTime,ordinal,size),value($type,archived,avatarUrl,buildLink,color(id),fullName,id,isResolved,localizedName,login,markdownText,minutes,name,presentation,ringId,text)),id,idReadable,isDraft,numberInProject,project($type,archived,id,name,plugins(timeTrackingSettings(enabled,estimate(field(id,name),id),timeSpent(field(id,name),id))),ringId,shortName),reporter($type,id,login,ringId),resolved,subtasks(id,issuesSize,unresolvedIssuesSize)"
	req, err := http.NewRequest("POST", currentUrl, bytes.NewBuffer(jsonData))
	
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err, "Ошибка при получении отчета")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err, "Error while reading the response bytes")

	var tasks []Task
	json.Unmarshal(body, &tasks)

	sort.SliceStable(tasks, func(i, j int) bool {
		return strings.Compare(tasks[i].Id, tasks[j].Id) > 0
	})
	sort.SliceStable(taskIds, func(i, j int) bool {
		return strings.Compare(taskIds[i].Id, taskIds[j].Id) > 0
	})

	taskInfoList := make([]string, len(tasks))
	for idx, _ := range taskInfoList {
		taskInfoList[idx] = fmt.Sprintf("%s: %s", tasks[idx].IdReadable, taskIds[idx].Summary)
	}
	return taskInfoList;
}
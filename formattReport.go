package main

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

func formatSpentReport(group ReportGroupItem) string {
	var rowList []string
	for _, line := range group.Lines {
		title := line.IssueId + ": " + line.Description
		duration := minutesToString(line.TotalDuration.Value)
		if LINE_LENGHT-utf8.RuneCountInString(title)-utf8.RuneCountInString(duration) <= 0 {
			title = title[:LINE_LENGHT-utf8.RuneCountInString(duration)-4] + "..."
		}

		dashLine := strings.Repeat(".", LINE_LENGHT-utf8.RuneCountInString(title)-utf8.RuneCountInString(duration))
		rowList = append(rowList, fmt.Sprintf("%s%s%s", title, dashLine, duration))

	}

	return fmt.Sprintf("<b>Отработано в %s %s</b>\n<pre>\n%s \n</pre>",
		getMonthTranslate(reportDate.Format("2006-01")),
		strings.Split(reportDate.Format("2006-01"), "-")[0],
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

		if lengthLeft-utf8.RuneCountInString(title) <= 0 {
			title = title[:lengthLeft-utf8.RuneCountInString(title)-4] + "..."
		}
		dashLineLeft := strings.Repeat(".", lengthLeft-utf8.RuneCountInString(title))
		dashLineRight := strings.Repeat(".", lengthRight-utf8.RuneCountInString(estimate)-utf8.RuneCountInString(cost))
		rowList = append(rowList, fmt.Sprintf("%s%s%s%s%s", title, dashLineLeft, estimate, dashLineRight, cost))
	}
	totalMin := minutesToString(group.Estimation.Value)
	totalCost := minutesToCost(group.Estimation.Value)
	dashLineLeft := strings.Repeat(".", lengthLeft-utf8.RuneCountInString("Итого:"))
	dashLineRight := strings.Repeat(".", lengthRight-utf8.RuneCountInString(totalMin)-utf8.RuneCountInString(totalCost))
	rowList = append(rowList, fmt.Sprintf("Итого:%s%s%s%s", dashLineLeft, totalMin, dashLineRight, totalCost))

	return fmt.Sprintf("<b>Внедрено в %s %s</b>\n<pre>\n%s \n</pre>",
		getMonthTranslate(reportDate.Format("2006-01")),
		strings.Split(reportDate.Format("2006-01"), "-")[0],
		strings.Join(rowList, "\n"))
}

func futureTaskFormat(planTasks, priorityTasks []string) (string, string) {
	plan := fmt.Sprintf("<b>Плановые задачи в %s %s</b>\n<pre>\n%s \n</pre>", 
		getMonthTranslate(reportDate.AddDate(0, 1, 0).Format("2006-01")),
		strings.Split(reportDate.AddDate(0, 1, 0).Format("2006-01"), "-")[0],
		strings.Join(sort.StringSlice(planTasks), "\n"))
	priority := fmt.Sprintf("<b>Приоритетные задачи в %s %s</b>\n<pre>\n%s \n</pre>", 
		getMonthTranslate(reportDate.AddDate(0, 1, 0).Format("2006-01")),
		strings.Split(reportDate.AddDate(0, 1, 0).Format("2006-01"), "-")[0],
		strings.Join(priorityTasks, "\n"))
	return plan, priority
}
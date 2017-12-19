package alert

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hubidu/e2e-backend/report-lib/model"
)

func hasConsecutiveFailures(reportGroup model.ReportGroup, numberOfFailures int) bool {
	if len(reportGroup.Items) >= numberOfFailures {
		failedReports := []model.ReportSlim{}

		for _, report := range reportGroup.Items {
			if report.Result != "error" {
				break
			}
			failedReports = append(failedReports, report)
		}

		return len(failedReports) >= numberOfFailures
	}
	return false
}

func getFailedReports(reportGroups []model.ReportGroup) []model.Report {
	failedReports := []model.Report{}
	for _, reportGroup := range reportGroups {
		// TODO Make this configurable
		if hasConsecutiveFailures(reportGroup, 2) {
			failedReports = append(failedReports, reportGroup.LastReport)
		}
	}
	return failedReports
}

func getSuccessfulReports(reportGroups []model.ReportGroup) []model.Report {
	successfulReports := []model.Report{}
	for _, reportGroup := range reportGroups {
		if reportGroup.LastReport.Result == "success" {
			successfulReports = append(successfulReports, reportGroup.LastReport)
		}
	}
	return successfulReports
}

func AlertTask() {
	fmt.Println("Checking alerts...")
	resp, err := GetReportCategories()
	if err != nil {
		log.Fatal(err)
	}

	groupedReports := make(map[uint32]*model.ReportGroup)
	json.Unmarshal(resp.Body(), &groupedReports)

	reportGroups := []model.ReportGroup{}
	for _, v := range groupedReports {
		reportGroups = append(reportGroups, *v)
	}

	failedReports := getFailedReports(reportGroups)
	successfulReports := getSuccessfulReports(reportGroups)

	newAlerts, _ := UpdateAlertedReports(failedReports, successfulReports)

	if len(newAlerts) > 0 {
		fmt.Println("Got new alerts", len(newAlerts))
		// email.SendAlert(newAlerts, fixedAlerts)
	}
}

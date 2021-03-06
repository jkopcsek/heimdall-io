package email

import (
	"crypto/tls"
	"fmt"

	"github.com/hubidu/e2e-backend/alert-service/config"
	"github.com/hubidu/e2e-backend/alert-service/service"
	"github.com/hubidu/e2e-backend/report-lib/model"
	gomail "gopkg.in/gomail.v2"
)

func formatMessage(reports []model.Report) string {
	greetingLine := fmt.Sprintf("Hi there!\n\nWe have %d new, repeated test failures. You should definitely have a look:\n\n", len(reports))

	content := ""
	for _, r := range reports {
		content += fmt.Sprintf("  - at %s [%s] %s\n", r.Started, r.DeviceSettings.Name, r.Title)
		content += fmt.Sprintf("    message \"%s\"\n", r.Screenshots[0].Message)
	}

	return greetingLine + content
}

func formatSubject(reports []model.Report) string {
	return fmt.Sprintf("We have %d NEW test failure(s)", len(reports))
}

// SendAlert sends an alert email using a list of new failing tests
func SendAlert(recipients []string, newAlerts []model.Report, fixedAlerts []model.Report, newAlertScreenshots []service.DownloadedScreenshot) {
	smtpConfig := config.NewSMTPConfig()
	alertConfig := config.NewAlertConfig()

	m := gomail.NewMessage()
	m.SetHeader("From", alertConfig.From)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", formatSubject(newAlerts))

	// var failedTests []string
	// for _, report := range newAlerts {
	// 	failedTests = append(failedTests, report.Title)
	// }
	m.SetBody("text/plain", formatMessage(newAlerts))

	fmt.Println("Attaching screenshots", newAlertScreenshots)
	for _, alertScreenshot := range newAlertScreenshots {
		m.Attach(alertScreenshot.Path)
	}

	// TODO Cleanup the screenshot files

	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, "", "")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send alert", err)
	}
}

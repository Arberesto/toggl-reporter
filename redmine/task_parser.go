package redmine

import (
	"fmt"
	"goreporter/report"
	"goreporter/utils"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ReportGenerator struct {
	BaseUrl string
}

// FindTaskId Finds task ID in the string "Task #42: description" or "42: description".
func FindTaskId(description string) (int, bool) {
	description = strings.TrimSpace(description)
	if strings.HasPrefix(description, "Task #") {
		description = description[6:]
		parts := strings.Split(description, ": ")
		if taskIdx, err := strconv.Atoi(parts[0]); err == nil {
			return taskIdx, true
		}
	} else {
		parts := strings.Split(description, ": ")
		if taskIdx, err := strconv.Atoi(parts[0]); err == nil {
			return taskIdx, true
		}
	}
	return 0, false
}

type TasksBlock struct {
	Tasks map[string]string
}

type RedmineForma struct {
	Date    string
	Hours   string
	Comment string
}

func DurationFormat(duration time.Duration) string {
	return fmt.Sprintf("%02d:%02d", utils.Hours(duration), utils.Minutes(duration))
}

func DateFormat(date time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", date.Year(), int(date.Month()), date.Day())
}

func BuildRedmineUrl(baseUrl string, id int, comment string, duration time.Duration, at time.Time) string {
	query := url.Values{}
	query.Set("time_entry[hours]", DurationFormat(duration))
	query.Set("time_entry[spent_on]", DateFormat(at))
	query.Set("time_entry[comments]", comment)

	rawUrl := fmt.Sprintf("%s/issues/%d/time_entries/new", baseUrl, id)
	u, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func (form *ReportGenerator) BuildRedmineReportForms(report report.Report) map[int]map[string]string {
	rreport := make(map[int]map[string]string)
	for projectID, project := range report.Projects {
		tasks := make(map[string]string)
		for task, duration := range project.Paid.Tasks {
			if idx, ok := FindTaskId(task); ok {
				tasks[task] = BuildRedmineUrl(form.BaseUrl, idx, task, duration, report.At)
			}
		}
		rreport[projectID] = tasks
	}
	return rreport
}

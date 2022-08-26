package weeklyRecurrance

import (
	dateUtils "gitlab-issue-automation/date_utils"
	gitlabUtils "gitlab-issue-automation/gitlab_utils"
	types "gitlab-issue-automation/types"
	"log"
	"math"
	"time"

	"github.com/xanzy/go-gitlab"
)

func GetNext(nextTime time.Time, data *types.Metadata, verbose bool) time.Time {
	if data.WeeklyRecurrence > 1 {
		git := gitlabUtils.GetGitClient()
		project := gitlabUtils.GetGitProject()
		orderBy := "created_at"
		options := &gitlab.ListProjectIssuesOptions{
			Search:  &data.Title,
			OrderBy: &orderBy,
		}
		issues, _, err := git.Issues.ListProjectIssues(project.ID, options)
		if err != nil {
			log.Fatal(err)
		}
		lastCreationDate := *issues[0].CreatedAt
		lastCreationWeek := dateUtils.GetStartOfWeek(lastCreationDate)
		currentWeek := dateUtils.GetStartOfWeek(time.Now())
		nextIssueWeek := lastCreationWeek
		log.Println("-- Degugging start")
		log.Println("--- Original next time")
		log.Println(nextTime)
		log.Println("--- Last creation date")
		log.Println(lastCreationDate)
		log.Println("--- Last creation week")
		log.Println(lastCreationWeek)
		log.Println("--- Current week")
		log.Println(currentWeek)
		for {
			nextIssueWeek = nextIssueWeek.AddDate(0, 0, 7*data.WeeklyRecurrence)
			log.Println("--- Next week iteration")
			log.Println(nextIssueWeek)
			if nextIssueWeek.Equal(currentWeek) || nextIssueWeek.After(currentWeek) {
				print("-- Break condition reached (next week equal or after current week)")
				break
			}
		}
		daysToAdd := math.Round(nextIssueWeek.Sub(currentWeek).Hours() / 24)
		if verbose {
			dueInWeeks := math.Round(math.Abs(daysToAdd / 7))
			log.Println("-- Next", data.WeeklyRecurrence, "weekly occurrence for", data.Title, "in", dueInWeeks, "week(s)")
		}
		nextTime = nextTime.AddDate(0, 0, int(daysToAdd))
		log.Println("--- Next week")
		log.Println(nextIssueWeek)
		log.Println("---Days to add")
		log.Println(daysToAdd)
		log.Println("--- Next time")
		log.Println(nextTime)
		log.Println("-- Degugging end")
	}
	return nextTime
}

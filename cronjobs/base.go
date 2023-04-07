package cronjobs

import (
	"fmt"
	"strings"
	"time"

	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

var (
	cronJobs = map[string]CronJobObject{
		"auto-verify-ids": {CronJob: VerifyIDs, Interval: time.Second * 10},
	}
)

type CronJob func(extReq request.ExternalRequest, db postgresql.Databases)

type CronJobObject struct {
	CronJob  CronJob
	Interval time.Duration
}

func Scheduler(extReq request.ExternalRequest, db postgresql.Databases, cronJob CronJob, interval time.Duration) {
	for {
		cronJob(extReq, db)
		time.Sleep(interval)
	}
}

func SetupCronJobs(extReq request.ExternalRequest, db postgresql.Databases, selectedJobs []string) {

	for _, v := range selectedJobs {
		jobName := strings.ToLower(v)
		cronJob, ok := cronJobs[jobName]

		if ok {
			utility.LogAndPrint(extReq.Logger, fmt.Sprintf("starting cronjob: %s", jobName))
			go Scheduler(extReq, db, cronJob.CronJob, cronJob.Interval)
		} else {
			utility.LogAndPrint(extReq.Logger, fmt.Sprintf("Cronjob not found: %s", jobName))
		}

	}

	select {}
}

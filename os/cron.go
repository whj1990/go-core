package os

import (
	"github.com/robfig/cron/v3"
)

var c *Cron

type Cron struct {
	*cron.Cron
	cronMap map[int64]cron.EntryID
}

type Job interface {
	Spec() string
	Cmd() func()
}

type JobsQuote struct {
	Jobs []Job
}

func Init(jobs ...Job) {
	c = &Cron{cron.New(cron.WithSeconds()), make(map[int64]cron.EntryID)}
	for _, job := range jobs {
		c.AddFunc(job.Spec(), job.Cmd())
	}
	c.Start()
}

package svc_cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func f1() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}

func f2(name string) func() {
	return func() {
		fmt.Printf("Hello, %s! (%v)\n", name, time.Now().Format("2006-01-02 15:04:05"))
	}
}

type Job struct {
	Name string
}

func (j Job) Run() {
	fmt.Printf("Hello, %s! (%v)\n", j.Name, time.Now().Format("2006-01-02 15:04:05"))
}

func CronInit() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))

	_, err := crontab.AddFunc("*/3 * * * * *", f1)
	if err != nil {
		fmt.Println(err)
	}

	_, err = crontab.AddFunc("*/3 * * * * *", f2("world"))
	if err != nil {
		fmt.Println(err)
	}

	_, err = crontab.AddJob("*/3 * * * * *", Job{"HelloCron"})
	if err != nil {
		fmt.Println(err)
	}

	crontab.Start() // Start a goroutine, we need to block main goroutine
}

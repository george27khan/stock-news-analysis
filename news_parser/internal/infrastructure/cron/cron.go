package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
)

type Scheduler struct {
	cron   *cron.Cron
	job    func(ctx context.Context) error
	option string
}

func NewScheduler(job func(ctx context.Context) error, option string) *Scheduler {
	return &Scheduler{cron.New(cron.WithSeconds()), job, option}
}

func (s *Scheduler) Start(ctx context.Context) {
	_, err := s.cron.AddFunc(s.option, func() {
		if err := s.job(ctx); err != nil {
			log.Printf("cron job failed: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("failed to add cron job: %v", err)
	}
	s.cron.Start()
	log.Println("cron started")
}

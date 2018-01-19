package publish2

import "time"

type Schedule struct {
	ScheduledStartAt *time.Time `gorm:"index"`
	ScheduledEndAt   *time.Time `gorm:"index"`
	ScheduledEventID *uint
}


package models

import "time"

type CalendarEventType string

const (
	EventExamType     CalendarEventType = "exam"
	EventTypeHoliday  CalendarEventType = "holiday"
	EventTypeEvent    CalendarEventType = "event"
	EventTypeDeadline CalendarEventType = "deadline"
	EventTypeOther    CalendarEventType = "other"
)

type CalendarBlock struct {
	ID          int               `db:"id" json:"id"`
	CollegeID   int               `db:"college_id" json:"college_id"`
	Title       string            `db:"title" json:"title"`
	Description string            `db:"description" json:"description"`
	EventType   CalendarEventType `db:"event_type" json:"event_type"`
	Date        time.Time         `db:"date" json:"date"`
}

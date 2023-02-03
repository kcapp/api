package models

import (
	"github.com/jordic/goics"
	"time"
)

// Feed is an iCal feed
type Feed struct {
	Content   string
	ExpiresAt time.Time
}

// Entry is a time entry
type Entry struct {
	DateStart   time.Time `json:"dateStart"`
	DateEnd     time.Time `json:"dateEnd"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
}

// Entries is a collection of entries
type Entries []*Entry

// EmitICal implements the interface for goics
func (e Entries) EmitICal() goics.Componenter {
	c := goics.NewComponent()
	c.SetType("VCALENDAR")
	c.AddProperty("CALSCAL", "GREGORIAN")
	c.AddProperty("VERSION", "2.0")

	for _, entry := range e {
		s := goics.NewComponent()
		s.SetType("VEVENT")
		s.AddProperty("ORGANIZER", "kcapp")

		k, v := goics.FormatDateTimeField("DTSTART", entry.DateStart)
		s.AddProperty(k, v)
		k, v = goics.FormatDateTimeField("DTEND", entry.DateEnd)
		s.AddProperty(k, v)

		s.AddProperty("SUMMARY", entry.Summary)
		s.AddProperty("DESCRIPTION", entry.Description)
		s.AddProperty("LOCATION", entry.Location)

		r := goics.NewComponent()
		r.SetType("VALARM")
		r.AddProperty("ACTION", "DISPLAY")
		r.AddProperty("DESCRIPTION", "REMINDER")
		r.AddProperty("TRIGGER", "-PT15M")
		s.AddComponent(r)

		c.AddComponent(s)
	}
	return c
}

package observability

import "time"

type Event struct {
	Kind       string
	Subject    string
	Tenant     string
	Attributes map[string]string
	At         time.Time
}

func NewEvent(kind, subject, tenant string, at time.Time) Event {
	return Event{Kind: kind, Subject: subject, Tenant: tenant, Attributes: map[string]string{}, At: at}
}
func (e *Event) Set(key, value string) { e.Attributes[key] = value }
func (e Event) Safe() Event            { return e }

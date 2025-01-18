package model

type Status int

const (
	UnkwownStatus Status = iota
	DefaultStatus
	SendingData
)

func (s Status) String() string {
	switch s {
	case DefaultStatus:
		return "default"
	case SendingData:
		return "sending data"
	default:
		return "unknown"
	}
}

package domain

type Role int

const (
	UnknownRole Role = iota
)

func (r Role) String() string {
	switch r {
	default:
		return "unknown"
	}
}

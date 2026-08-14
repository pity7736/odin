package odinerrors

type Tag uint8

const (
	Unknown Tag = iota
	Domain
	Render
	NotFound
	Unauthorized
)

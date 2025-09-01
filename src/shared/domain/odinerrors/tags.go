package odinerrors

type Tag uint8

const (
	UNKNOWN Tag = iota
	DOMAIN
	RENDER
	NOT_FOUND
)

package gofilemirror

type OperationType int

const (
	OT_NONE OperationType = iota
	// OT_READ
	OT_READ_AT
	// OT_SEEK
	// OT_TRUNCATE
	// OT_WRITE
	OT_WRITE_AT
	// OT_WRITE_STRING
	// OT_SYNC
)

package gofilemirror

type AsyncOperationType int

const (
	AOT_NONE AsyncOperationType = iota
	AOT_READ
	AOT_READ_AT
	AOT_SEEK
	AOT_TRUNCATE
	AOT_WRITE
	AOT_WRITE_AT
	AOT_WRITE_STRING
)

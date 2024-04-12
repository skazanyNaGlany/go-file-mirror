package gofilemirror

type AsyncOperationType int

const (
	READ = iota + 1
	READ_AT
	SEEK
	TRUNCATE
	WRITE
	WRITE_AT
	WRITE_STRING
)

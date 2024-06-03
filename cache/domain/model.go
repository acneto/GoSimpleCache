package domain

type CmdType string

const (
	CMDSet CmdType = "SET"
	CMDGet CmdType = "GET"
)

type Command struct {
	Type  CmdType
	Key   string
	Value string
}

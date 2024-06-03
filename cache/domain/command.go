package domain

import (
	"errors"
	"strings"
)

func ParseCommand(bytes []byte) (*Command, error) {
	var (
		rawStr = string(bytes)
		parts  = strings.Split(rawStr, " ")
	)

	var cmd Command
	cmdType := CmdType(parts[0])

	if cmdType == CMDSet {
		if len(parts) != 3 {
			return nil, errors.New("wrong arguments for SET command")
		}
		cmd = Command{
			Type:  CMDSet,
			Key:   strings.TrimSpace(parts[1]),
			Value: parts[2],
		}
		return &cmd, nil
	}

	if cmdType == CMDGet {
		if len(parts) != 2 {
			return nil, errors.New("wrong arguments for GET command")
		}
		cmd = Command{
			Type: CMDGet,
			Key:  strings.TrimSpace(parts[1]),
		}
		return &cmd, nil
	}

	switch cmd.Type {
	case CMDSet:
		return &cmd, nil
	case CMDGet:
		return &cmd, nil
	}
	return nil, errors.New("invalid command")
}

package cli

import (
	"flag"
	"fmt"
)

type TypeConnection string

const (
	Local TypeConnection = "local"
	SMB   TypeConnection = "smb"
	HTTP  TypeConnection = "http"
)

func NewTypeConnection(s string) (TypeConnection, error) {
	switch s {
	case string(Local):
		return Local, nil
	case string(SMB):
		return SMB, nil
	case string(HTTP):
		return HTTP, nil
	default:
		return "", fmt.Errorf("invalid TypeConnection: %s", s)
	}
}

func (t TypeConnection) Validate() bool {
	switch t {
	case Local, SMB, HTTP:
		return true
	}

	return false
}

type CMDFlags struct {
	IsDebug       bool
	EnableLogFile bool
	Path          string
	Conn          TypeConnection
}

func GetCMDFlags() (*CMDFlags, error) {
	var typeConn = flag.String("typeConn", "local", "")
	var path = flag.String("path", "", "config address")
	var isDebug = flag.Bool("debug", false, "")
	var enableLogFile = flag.Bool("enableLogFile", false, "")

	//var err error
	flag.Parse()

	// var path string

	if typeConn == nil || *typeConn == "" {
		return nil, fmt.Errorf("invalid typeConn to file")
	}

	if path == nil || *path == "" {
		return nil, fmt.Errorf("invalid path to file")
	}

	conn, err := NewTypeConnection(*typeConn)
	if err != nil {
		return nil, err
	}

	return &CMDFlags{
		Path:          *path,
		IsDebug:       *isDebug,
		Conn:          conn,
		EnableLogFile: *enableLogFile,
	}, nil
}

// func validateURL

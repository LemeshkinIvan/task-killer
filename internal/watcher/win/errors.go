package win

import "errors"

var ErrBlacklistLen error = errors.New("Win32Watcher: blacklist is empty")
var ErrLogNullable error = errors.New("Win32Watcher: log is nil")
var ErrProcessAccess error = errors.New("Win32Watcher: Access is denied.")
var ErrHandle error = errors.New("Win32Watcher: invalid handle value")

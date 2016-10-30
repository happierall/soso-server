package soso

import (
	"net/http"
	"strconv"
)

type Level int64

const (
	// System is unusable(Ex: This level should not be used by applications)
	LevelEmergency Level = iota
	// Should be corrected immediately(Ex: Loss of the primary ISP connection)
	LevelAlert
	// Critical conditions(Ex: A failure in the system's primary application)
	LevelCritical
	// Error conditions(Ex: An application has exceeded its file storage limit and attempts to write are failing)
	LevelError
	// May indicate that an error will occur if action is not taken (Ex: A non-root file system has only 2GB remaining)
	LevelWarning
	// Events that are unusual, but not error conditions.
	LevelNotice
	// Normal operation events that require no action (Ex: An application has started, paused or ended successfully.
	LevelInfo
	// Information useful to developers for debugging an application
	LevelDebug
)

var (
	LastLogID int = 0
)

var levels = [...]string{
	"emerg",
	"alert",
	"crit",
	"err",
	"warn",
	"notice",
	"info",
	"service/debug",
}

func (l Level) String() string { return levels[l] }

type Log struct {
	CodeKey string `json:"code_key"`
	CodeStr string `json:"code_str"`

	LevelInt int    `json:"level_int"`
	LevelStr string `json:"level_str"`

	LogID   string `json:"log_id"`
	UserMsg string `json:"user_msg"`
}

func NewLog(code_key int, lvl_str Level, user_msg string) Log {
	LastLogID++

	return Log{
		CodeKey: strconv.Itoa(code_key),
		CodeStr: http.StatusText(code_key),

		LevelInt: int(lvl_str),
		LevelStr: lvl_str.String(),

		LogID:   strconv.Itoa(LastLogID),
		UserMsg: user_msg,
	}
}

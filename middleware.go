package soso

import "time"

type middleware struct {
	before func(*Msg, time.Time)
	after  func(*Msg, time.Duration)
}

func (mw *middleware) Before(handler func(m *Msg, startTime time.Time)) {
	mw.before = handler
}

func (mw *middleware) After(handler func(m *Msg, elapsed time.Duration)) {
	mw.after = handler
}

func (mw *middleware) beforeExecute(m *Msg, startTime time.Time) {
	if mw.before != nil {
		mw.before(m, startTime)
	}
}

func (mw *middleware) afterExecute(m *Msg, elapsed time.Duration) {
	if mw.after != nil {
		mw.after(m, elapsed)
	}
}

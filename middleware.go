package soso

import "time"

type middleware struct {
	beforeList []func(*Msg, time.Time)
	afterList  []func(*Msg, time.Duration)
}

func (mw *middleware) Before(handler func(m *Msg, startTime time.Time)) {
	mw.beforeList = append(mw.beforeList, handler)
}

func (mw *middleware) After(handler func(m *Msg, elapsed time.Duration)) {
	mw.afterList = append(mw.afterList, handler)
}

func (mw *middleware) beforeExecute(m *Msg, startTime time.Time) {
	for _, handler := range mw.beforeList {
		handler(m, startTime)
	}
}

func (mw *middleware) afterExecute(m *Msg, elapsed time.Duration) {
	for _, handler := range mw.afterList {
		handler(m, elapsed)
	}
}

package soso

import (
	"errors"
	"net/http"
	"time"

	"github.com/fatih/color"
)

type (
	HandlerFunc func(*Msg)

	Route struct {
		Model  string
		Action string

		Handler HandlerFunc
	}

	Router struct {
		Delay int // ms. For testing only

		Routes     []Route
		Middleware middleware
	}

	Routes struct {
		List []Route
	}
)

func (r *Router) Handle(model string, action string, handler HandlerFunc) {
	route := Route{
		Model:   model,
		Action:  action,
		Handler: handler,
	}

	r.Routes = append(r.Routes, route)
}

func (r *Router) HandleList(routes []Route) {
	r.Routes = append(r.Routes, routes...)
}

func (r *Router) HandleRoutes(routes Routes) {
	r.Routes = append(r.Routes, routes.List...)
}

func (r *Router) processIncomingMsg(session Session, data []byte) {
	req, err := NewRequest(data)

	if err != nil {
		Loger.Errorf("%s Error: incorrect request - %s\n%s", logPrefix, data, err)
		return
	}

	msg := newMsgFromRequest(req, session)

	if msg == nil {
		return
	}

	// Find handler and execute
	found := false

	for _, route := range r.Routes {
		if req.Action == route.Action && req.Model == route.Model {

			startTime := time.Now()

			r.Middleware.beforeExecute(msg, startTime)

			if r.Delay != 0 {
				time.Sleep(time.Duration(r.Delay) * time.Millisecond)
			}

			if session.IsClosed() {
				Loger.Infof("%s %s | %s -> %s | %s\n",
					logPrefix,
					time.Now().Format("2006/01/02 - 15:04:05"),
					color.RedString(req.Model),
					color.GreenString(req.Action),
					"Session is closed",
				)
				return
			}

			route.Handler(msg)

			elapsedTime := time.Since(startTime)

			Loger.Infof("%s %s | %s -> %s | %s\n",
				logPrefix,
				startTime.Format("2006/01/02 - 15:04:05"),
				color.YellowString(req.Model),
				color.GreenString(req.Action),
				elapsedTime,
			)
			found = true

			r.Middleware.afterExecute(msg, elapsedTime)
		}
	}

	if found != true {
		Loger.Infof("%s %s | %s -> %s | %s\n",
			logPrefix,
			time.Now().Format("2006/01/02 - 15:04:05"),
			color.RedString(req.Model),
			color.GreenString(req.Action),
			"Route not found",
		)
		msg.Error(http.StatusNotFound, LevelError, errors.New("No model handler found"))
	}

}

func (r *Routes) Handle(model, action string, handler HandlerFunc) {
	route := Route{
		Model:   model,
		Action:  action,
		Handler: handler,
	}

	r.List = append(r.List, route)
}

package soso

import (
	"errors"
	"fmt"
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

		Routes []Route
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
		fmt.Printf("%s Error: incorrect request - %s\n", logPrefix, data)
		fmt.Println(err)
		return
	}

	msg := newMsgFromRequest(req, session)

	if msg == nil {
		return
	}

	// Find handler and execute
	found := false

	for _, route := range r.Routes {
		fmt.Println(req.Action, route.Action, req.Model, route.Model)
		if req.Action == route.Action && req.Model == route.Model {

			startTime := time.Now()

			if r.Delay != 0 {
				time.Sleep(time.Duration(r.Delay) * time.Millisecond)
			}

			if session.IsClosed() {
				fmt.Printf("%s %s | %s -> %s | %s\n",
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

			fmt.Printf("%s %s | %s -> %s | %s\n",
				logPrefix,
				startTime.Format("2006/01/02 - 15:04:05"),
				color.YellowString(req.Model),
				color.GreenString(req.Action),
				elapsedTime,
			)
			found = true
		}
	}

	if found != true {
		fmt.Printf("%s %s | %s -> %s | %s\n",
			logPrefix,
			time.Now().Format("2006/01/02 - 15:04:05"),
			color.RedString(req.Model),
			color.GreenString(req.Action),
			"Route not found",
		)
		msg.Error(http.StatusNotFound, LevelError, errors.New("No model handler found"))
	}

}

func (r *Routes) Add(model, action string, handler HandlerFunc) {
	route := Route{
		Model:   model,
		Action:  action,
		Handler: handler,
	}

	r.List = append(r.List, route)
}

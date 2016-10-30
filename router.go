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
		DataType  string
		ActionStr string

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

func (r *Router) Handle(data_type string, action_str string, handler HandlerFunc) {
	route := Route{
		ActionStr: action_str,
		DataType:  data_type,
		Handler:   handler,
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
		if req.ActionStr == route.ActionStr && req.DataType == route.DataType {

			startTime := time.Now()

			if r.Delay != 0 {
				time.Sleep(time.Duration(r.Delay) * time.Millisecond)
			}

			if session.IsClosed() {
				fmt.Printf("%s %s | %s -> %s | %s\n",
					logPrefix,
					time.Now().Format("2006/01/02 - 15:04:05"),
					color.RedString(req.DataType),
					color.GreenString(req.ActionStr),
					"Session is closed",
				)
				return
			}

			route.Handler(msg)

			elapsedTime := time.Since(startTime)

			fmt.Printf("%s %s | %s -> %s | %s\n",
				logPrefix,
				startTime.Format("2006/01/02 - 15:04:05"),
				color.YellowString(req.DataType),
				color.GreenString(req.ActionStr),
				elapsedTime,
			)
			found = true
		}
	}

	if found != true {
		fmt.Printf("%s %s | %s -> %s | %s\n",
			logPrefix,
			time.Now().Format("2006/01/02 - 15:04:05"),
			color.RedString(req.DataType),
			color.GreenString(req.ActionStr),
			"Route not found",
		)
		msg.Error(http.StatusNotFound, LevelError, errors.New("No model handler found"))
	}

}

func (r *Routes) Add(data_type string, action_str string, handler HandlerFunc) {
	route := Route{
		DataType:  data_type,
		ActionStr: action_str,
		Handler:   handler,
	}

	r.List = append(r.List, route)
}

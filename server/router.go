package main

import "errors"

type router struct {
	routes map[string]func(serverContext, string) int
}

func createRouter() router {
	var r router
	r.routes = make(map[string]func(serverContext, string) int)
	return r
}

func (r router) addRoute(route string, handler func(serverContext, string) int) {
	r.routes[route] = handler
}

func (r router) getHandler(route string) (handler func(serverContext, string) int, err error) {
	if h, ok := r.routes[route]; ok {
		return h, nil
	} else {
		return nil, errors.New("There is no handler for: " + route)
	}
}

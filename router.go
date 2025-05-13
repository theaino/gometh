package meth

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type RouteHandler = func(*fiber.Ctx) templ.Component

type Route struct {
	Path string
	Method string
	Handler RouteHandler
}

type Router struct {
	Routes []Route
}

func (r *Router) Get(path string, handler RouteHandler) {
	r.Match(path, handler, fiber.MethodGet)
}

func (r *Router) Put(path string, handler RouteHandler) {
	r.Match(path, handler, fiber.MethodPut)
}

func (r *Router) Head(path string, handler RouteHandler) {
	r.Match(path, handler, fiber.MethodHead)
}

func (r *Router) Post(path string, handler RouteHandler) {
	r.Match(path, handler, fiber.MethodPost)
}

func (r *Router) Patch(path string, handler RouteHandler) {
	r.Match(path, handler, fiber.MethodPatch)
}

func (r *Router) Trace(path string, handler RouteHandler) {
	r.Match(path, handler, fiber.MethodTrace)
}

func (r *Router) Delete(path string, handler RouteHandler) {
	r.Match(path, handler, fiber.MethodDelete)
}

func (r *Router) Connect(path string, handler RouteHandler) {
	r.Match(path, handler, fiber.MethodConnect)
}

func (r *Router) Options(path string, handler RouteHandler) {
	r.Match(path, handler, fiber.MethodOptions)
}

func (r *Router) Match(path string, handler RouteHandler, methods... string) {
	for _, method := range methods {
		r.Routes = append(r.Routes, Route{
			Path: path,
			Method: method,
			Handler: handler,
		})
	}
}

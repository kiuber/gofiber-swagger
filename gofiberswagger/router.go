package gofiberswagger

import (
	"github.com/gofiber/fiber/v3"
)

type SwaggerRouter struct {
	internalGroup string
	Router        fiber.Router
}

func NewRouter(app *fiber.App) SwaggerRouter {
	return NewRouterFromRouter(app.Group("/"))
}
func NewRouterFromRouter(r fiber.Router) SwaggerRouter {
	return SwaggerRouter{internalGroup: "", Router: r}
}

func (router SwaggerRouter) Get(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterRouteInternal("GET", path, router.internalGroup, docs)
	return router.Router.Get(path, handler, middleware...)
}
func (router SwaggerRouter) Head(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterRouteInternal("HEAD", path, router.internalGroup, docs)
	return router.Router.Head(path, handler, middleware...)
}
func (router SwaggerRouter) Post(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterRouteInternal("POST", path, router.internalGroup, docs)
	return router.Router.Post(path, handler, middleware...)
}
func (router SwaggerRouter) Put(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterRouteInternal("PUT", path, router.internalGroup, docs)
	return router.Router.Put(path, handler, middleware...)
}
func (router SwaggerRouter) Delete(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterRouteInternal("DELETE", path, router.internalGroup, docs)
	return router.Router.Delete(path, handler, middleware...)
}
func (router SwaggerRouter) Connect(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterRouteInternal("CONNECT", path, router.internalGroup, docs)
	return router.Router.Connect(path, handler, middleware...)
}
func (router SwaggerRouter) Options(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterRouteInternal("OPTIONS", path, router.internalGroup, docs)
	return router.Router.Options(path, handler, middleware...)
}
func (router SwaggerRouter) Trace(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterRouteInternal("TRACE", path, router.internalGroup, docs)
	return router.Router.Trace(path, handler, middleware...)
}
func (router SwaggerRouter) Patch(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterRouteInternal("PATCH", path, router.internalGroup, docs)
	return router.Router.Patch(path, handler, middleware...)
}
func (router *SwaggerRouter) Group(prefix string, handlers ...fiber.Handler) SwaggerRouter {
	return SwaggerRouter{internalGroup: router.internalGroup + prefix, Router: router.Router.Group(prefix, handlers...)}
}

func routerRegisterRouteInternal(method string, path string, internalGroup string, info *RouteInfo) {
	if info == nil {
		info = &RouteInfo{}
	}
	if internalGroup != "" {
		info.Tags = append(info.Tags, internalGroup)
	}
	RegisterRoute(method, internalGroup+path, info)
}

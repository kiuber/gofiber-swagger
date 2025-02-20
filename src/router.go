package swagger

import "github.com/gofiber/fiber/v3"

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
	routerRegisterPathInternal("GET", path, router.internalGroup, docs)
	return router.Router.Get(path, handler, middleware...)
}
func (router SwaggerRouter) Head(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterPathInternal("HEAD", path, router.internalGroup, docs)
	return router.Router.Head(path, handler, middleware...)
}
func (router SwaggerRouter) Post(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterPathInternal("POST", path, router.internalGroup, docs)
	return router.Router.Post(path, handler, middleware...)
}
func (router SwaggerRouter) Put(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterPathInternal("PUT", path, router.internalGroup, docs)
	return router.Router.Put(path, handler, middleware...)
}
func (router SwaggerRouter) Delete(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterPathInternal("DELETE", path, router.internalGroup, docs)
	return router.Router.Delete(path, handler, middleware...)
}
func (router SwaggerRouter) Connect(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterPathInternal("CONNECT", path, router.internalGroup, docs)
	return router.Router.Connect(path, handler, middleware...)
}
func (router SwaggerRouter) Options(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterPathInternal("OPTIONS", path, router.internalGroup, docs)
	return router.Router.Options(path, handler, middleware...)
}
func (router SwaggerRouter) Trace(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterPathInternal("TRACE", path, router.internalGroup, docs)
	return router.Router.Trace(path, handler, middleware...)
}
func (router SwaggerRouter) Patch(path string, docs *RouteInfo, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router {
	routerRegisterPathInternal("PATCH", path, router.internalGroup, docs)
	return router.Router.Patch(path, handler, middleware...)
}
func (router *SwaggerRouter) Group(prefix string, handlers ...fiber.Handler) SwaggerRouter {
	return SwaggerRouter{internalGroup: prefix, Router: router.Router.Group(prefix, handlers...)}
}

func routerRegisterPathInternal(method string, path string, internalGroup string, info *RouteInfo) {
	if info != nil && internalGroup != "" {
		info.Tags = append(info.Tags, internalGroup)
	}
	RegisterPath(method, path, info)
}

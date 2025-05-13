package meth

import (
	"net/http"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

type App struct {
	fiber *fiber.App
	router Router
	Conf Conf
}

func (a *App) route() {
	if a.fiber == nil {
		a.fiber = fiber.New()
	}
	for _, route := range a.router.Routes {
		a.fiber.Add(route.Method, route.Path, func(c *fiber.Ctx) error {
			c.Set("Content-Type", "text/html")
			return route.Handler(c).Render(c.Context(), c.Response().BodyWriter())
		})
	}
	if a.Conf.Build.DistDir != "" {
		a.fiber.Use(a.Conf.Build.DistDir, filesystem.New(filesystem.Config{
			Root: http.Dir(a.Conf.Build.DistDir),
			Browse: false,
		}))
	}
}

func (a *App) esbuild() {
	plugins := []api.Plugin{}
	if a.Conf.Build.Sass {
		plugins = append(plugins, sassPlugin())
	}
	api.Build(api.BuildOptions{
		EntryPoints: a.Conf.Build.Entrypoints,
		Outdir: a.Conf.Build.DistDir,
		Bundle: true,
		Write: true,
		LogLevel: api.LogLevelInfo,
		Plugins: plugins,
	})
}

func (a *App) Route(f func(r *Router)) {
	f(&a.router)
}

func (a *App) Run() error {
	a.route()
	if a.Conf.Build.Esbuild {
		a.esbuild()
	}
	return a.fiber.Listen(a.Conf.Server.Host + ":" + a.Conf.Server.Port)
}

package main

import (
	"flag"
	"log"
	"surflex-backend/api/controller"
	"surflex-backend/common/db"
	"time"

	"github.com/gofiber/fiber/v2/middleware/pprof"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	addr = flag.String("addr", ":8080", "TCP address to listen to")
)

func main() {
	flag.Parse()
	db.Init()

	app := fiber.New(fiber.Config{
		ReadTimeout: 60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		// AllowOriginsFunc: func(origin string) bool {
		// 	matched, _ := regexp.MatchString(`^https://([a-zA-Z0-9-]+)\.web\.app$`, origin)
		// 	return origin == "https://surflex.fun" || matched
		// },
		AllowHeaders: "*",
	}))

	api := app.Group("/api")

	api.Use(pprof.New(pprof.Config{Prefix: "/a18761"}))

	api.Use("/", func(c *fiber.Ctx) error {
		start := time.Now()
		defer func() {
			end := time.Since(start)

			log.Printf("[AccessLog] %s - %d %d %v - %s\n", getRemoteIP(c), c.Response().StatusCode(), len(c.Response().Body()), end, string(c.Request().RequestURI()))
		}()

		c.Response().Header.SetContentType("application/json")

		return c.Next()
	})

	api.Get("/health", func(ctx *fiber.Ctx) error {
		return nil
	})

	controller.SetPriceRouter(api.Group("/price"))
	controller.SetChartRouter(api.Group("/chart"))
	controller.SetRoundRouter(api.Group("/round"))
	controller.SetPositionRouter(api.Group("/position"))
	controller.SetAccountRouter(api.Group("/account"))
	controller.SetLeaderBoardRouter(api.Group("/leaderboard"))
	controller.SetAdminController(api.Group("/admin"))

	log.Println("Ready to Serve")

	if err := app.Listen(*addr); err != nil {
		log.Printf("Listen Error %+v\n", err)
	}
}

func getRemoteIP(ctx *fiber.Ctx) string {
	xForwardedFor := ctx.Request().Header.Peek("x-forwarded-for")
	if xForwardedFor != nil {
		return string(xForwardedFor)
	}

	return ctx.Context().RemoteIP().String()
}

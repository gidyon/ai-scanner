package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gidyon/ai-scanner/internal/imgcompress"
	"github.com/gidyon/gomicro/utils/errs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

var (
	configFile = flag.String("config-file", ".env", "Configuration file")
	keyFile    = flag.String("key-file", "", "tls private key")
	certFile   = flag.String("cert-file", "", "tls public key")
	httpPort   = flag.String("http-port", "30084", "Http port")
	httpsPort  = flag.String("https-port", "30085", "Https port")
)

func main() {
	flag.Parse()

	// Config in .env
	viper.SetConfigFile(*configFile)

	viper.AutomaticEnv()

	// Read config from .env
	err := viper.ReadInConfig()
	errs.Panic(err)

	app := fiber.New(fiber.Config{
		BodyLimit: 30 * 1024 * 1024,
	})

	// Fix cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,OPTIONS",
		ExposeHeaders: strings.Join([]string{
			"Accept",
			"Access-Control-Allow-Origin",
			"Authorization",
			"Cache-Control",
			"Content-Type",
			"DNT",
			"If-Modified-Since",
			"Keep-Alive",
			"Origin",
			"User-Agent",
			"X-Requested-With",
		}, ","),
		MaxAge:           1728,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	// API routes
	api := app.Group("/api")

	// API->V1 routes
	v1 := api.Group("/v1", func(c *fiber.Ctx) error { // middleware for /api/v1
		c.Set("Version", "v1")
		return c.Next()
	})

	// Compressing images
	v1.Post("/images:compress", imgcompress.CompressImage)

	v1.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("Am alive 👋!")
	})

	if *keyFile == "" || *certFile == "" {
		log.Fatal(app.Listen(":" + *httpPort))
	} else {
		log.Fatal(app.ListenTLS(":"+*httpsPort, *certFile, *keyFile))
	}
}

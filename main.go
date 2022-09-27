package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	_ "github.com/joho/godotenv/autoload"
)

var ctx context.Context
var rdb *redis.Client

func main() {
	initRedis()

	app := fiber.New()
	app.Get("/picture/:pictureName", handlePicture)
	app.Listen(":6060")
}

func initRedis() {
	redisPort, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}

	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisPort,
	})
}

func handlePicture(c *fiber.Ctx) error {
	url, err := rdb.Get(ctx, c.Params("pictureName")).Result()

	if queryArgs := c.Context().QueryArgs().String(); queryArgs != "" {
		url = url + "?" + queryArgs
	}

	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if err := proxy.Do(c, url); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Response().Header.Del(fiber.HeaderServer)
	c.Response().Header.Add("x-original-url", url)
	return nil
}

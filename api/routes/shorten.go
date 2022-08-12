package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/harshit2000/GoURL_Shortener/database"
)
package helpers

import (
	"time"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/harshit2000/GoURL_Shortener/database"

)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"custom_short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"x_rate_remaining"`
	XRateLimitReset time.Duration `json:"x_rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {

	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Implementation of rate Limiter

	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil{
		_ = r2.Set(database.Ctx, c.IP, os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else{
		val, _ = r2.Get(database.Ctx, c.IP.Result())
		valInt, _ := strconv.Atoi(val)
		// we will increment by 1 so thats why it is 0
		if valInt <= 0{
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Second,
			})
		}
	}

	// check if the input sent by user is actual URL

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})
	}

	//localhost check : for domain error

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fibre.StatusServiceUnavailable).JSON(fiber.Map{"error": "Service is down :) "})
	}

	// enforce HTTPS, SSL connection

	body.URL = helpers.EnforceHTTPS(body.URL)

	r2.Decr(database.Ctx, c.IP())
	

}

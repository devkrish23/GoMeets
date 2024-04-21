package handlers

import "github.com/gofiber/fiber/v2"

func Welcome(c *fiber.Ctx) error {
	return c.Render("welcome", nil, "layouts/main")// render our welcome page by sending response as http response to the fiber's built in template engine
}

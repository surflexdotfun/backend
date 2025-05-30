package controller

import (
	"surflex-backend/api/service"
	"surflex-backend/common/utils"

	"github.com/gofiber/fiber/v2"
)

func SetPriceRouter(r fiber.Router) {
	service.StartPriceUpdater()
	r.Get("/all", utils.Wrap(HandlePriceAll))
}

func HandlePriceAll(c *fiber.Ctx) (interface{}, error) {
	return service.GetPrice(), nil
}

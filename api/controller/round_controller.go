package controller

import (
	"surflex-backend/api/service"
	"surflex-backend/common/utils"

	"github.com/gofiber/fiber/v2"
)

func SetRoundRouter(r fiber.Router) {
	service.StartRoundUpdater()
	r.Get("/", utils.Wrap(HandleGetRound))
}

func HandleGetRound(c *fiber.Ctx) (interface{}, error) {
	return service.GetRound(), nil
}

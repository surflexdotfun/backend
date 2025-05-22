package controller

import (
	"surflex-backend/api/auth"
	"surflex-backend/api/service"
	"surflex-backend/common/utils"

	"github.com/gofiber/fiber/v2"
)

func SetAdminController(r fiber.Router) {
	r.Post("/leaderboard/update", utils.Wrap(HandleUpdateLeaderBoard))
}

func HandleUpdateLeaderBoard(c *fiber.Ctx) (interface{}, error) {
	if err := auth.CheckAdminAuthorization(c); err != nil {
		return nil, err
	}

	return service.UpdateLeaderBoard()
}

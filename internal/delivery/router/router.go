package router

import (
	"github.com/MCPutro/go-management-project/internal/delivery/handler"
	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router fiber.Router, handler handler.UserHandler) {
	users := router.Group("/users")

	users.Post("/", handler.CreateUser)
	users.Get("/:id", handler.GetUser)
	users.Put("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)
}

//// registerProjectRoutes registers all project-related routes
//func registerProjectRoutes(router fiber.Router, handler *Handler) {
//	projects := router.Group("/projects")
//
//	projects.Post("/", handler.CreateProject)
//	projects.Get("/:id", handler.GetProject)
//	projects.Put("/:id", handler.UpdateProject)
//	projects.Delete("/:id", handler.DeleteProject)
//}
//
//// registerListRoutes registers all list-related routes
//func registerListRoutes(router fiber.Router, handler *Handler) {
//	lists := router.Group("/lists")
//
//	lists.Post("/", handler.CreateList)
//	lists.Get("/:id", handler.GetList)
//	lists.Get("/project/:project_id", handler.GetListsByProject)
//	lists.Put("/:id", handler.UpdateList)
//	lists.Delete("/:id", handler.DeleteList)
//}
//
//// registerCardRoutes registers all card-related routes
//func registerCardRoutes(router fiber.Router, handler *Handler) {
//	cards := router.Group("/cards")
//
//	cards.Post("/", handler.CreateCard)
//	cards.Get("/:id", handler.GetCard)
//	cards.Get("/list/:list_id", handler.GetCardsByList)
//	cards.Put("/:id", handler.UpdateCard)
//	cards.Delete("/:id", handler.DeleteCard)
//}

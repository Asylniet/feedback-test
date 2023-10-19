package routes

import (
	"github.com/enzhas/feedback_back/controllers"
	"github.com/enzhas/feedback_back/middleware"
	"github.com/gin-gonic/gin"
)

type TodoRouteController struct {
	todoController controllers.TodoController
}

func NewTodoRouteController(todoController controllers.TodoController) TodoRouteController {
	return TodoRouteController{todoController}
}

func (tc *TodoRouteController) TodoRoute(rg *gin.RouterGroup) {
	router := rg.Group("todo")

	router.GET("/:orgId", middleware.DeserializeUser("any"), middleware.CheckSubscriptionExpireAt(), tc.todoController.GetTodos)
	router.POST("/", middleware.DeserializeUser("any"), middleware.CheckSubscriptionExpireAt(), tc.todoController.CreateTodo)
	router.POST("/vote", middleware.DeserializeUser("any"), middleware.CheckSubscriptionExpireAt(), tc.todoController.Vote)
	router.PUT("/:id", middleware.DeserializeUser("manager"), tc.todoController.EditTodo)
	router.DELETE("/:id", middleware.DeserializeUser("manager"), tc.todoController.DeleteTodo)
	router.GET("/changes/:id", middleware.DeserializeUser("manager"), tc.todoController.GetChangesOfToDo)
}

package controllers

import (
	"net/http"
	"sort"

	"github.com/enzhas/feedback_back/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TodoController struct {
	DB *gorm.DB
}

func NewTodoController(DB *gorm.DB) TodoController {
	return TodoController{DB}
}

func (tc *TodoController) getTodos(ctx *gin.Context, organizationId string) ([]map[string]interface{}, error) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	type resultTodo struct {
		Todo      models.Todo `json:"todo"`
		VoteCount int         `json:"vote_count"`
		UserVoted bool        `json:"user_voted"`
	}
	var todos []models.Todo
	if err := tc.DB.Preload("Votes").Preload("Sender").Where("organization_id = ?", organizationId).Find(&todos); err.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
		return nil, err.Error
	}
	dashboardTodos := map[string][]resultTodo{
		"todo":     nil,
		"planning": nil,
		"doing":    nil,
		"done":     nil,
	}
	stateTypes := []string{"todo", "planning", "doing", "done"}
	identifier := map[string]string{
		"todo":     "📄 Задача",
		"planning": "📝 Планируется",
		"doing":    "⚡ В процессе",
		"done":     "✅ Сделано",
	}

	// Group todos by state
	for i, todo := range todos {
		var todoResponse resultTodo
		todoResponse.Todo = todos[i]
		todoResponse.VoteCount = len(todos[i].Votes)

		if todoResponse.Todo.SenderID == currentUser.ID {
			todoResponse.UserVoted = true
		} else {
			for _, vote := range todos[i].Votes {
				if vote.SenderID == currentUser.ID {
					todoResponse.UserVoted = true
					break
				}
			}
		}

		dashboardTodos[todo.State] = append(dashboardTodos[todo.State], todoResponse)
	}
	for _, todo := range todos {
		sort.Slice(dashboardTodos[todo.State], func(i, j int) bool {
			return dashboardTodos[todo.State][i].VoteCount > dashboardTodos[todo.State][j].VoteCount
		})
	}
	// Create columns array
	var columns []map[string]interface{}
	for _, state := range stateTypes {
		column := map[string]interface{}{
			"id":    state,
			"name":  identifier[state],
			"tasks": dashboardTodos[state],
		}
		columns = append(columns, column)
	}

	return columns, nil
}

func (tc *TodoController) GetTodos(ctx *gin.Context) {
	orgId := ctx.Param("orgId")
	todos, err := tc.getTodos(ctx, orgId)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"todos": todos}})
}

func (tc *TodoController) Vote(ctx *gin.Context) {
	var todo models.Todo
	var voteInput models.VoteInput
	if err := ctx.ShouldBindJSON(&voteInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)

	vote := models.Vote{
		TodoID:   voteInput.TodoID,
		SenderID: currentUser.ID,
	}

	if result := tc.DB.Where("id = ?", voteInput.TodoID).First(&todo); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error})
		return
	}

	if todo.SenderID == vote.SenderID {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "creator can not vote his todo"})
		return
	}

	if result := tc.DB.Create(&vote); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error})
		return
	}

	orgID := (*todo.OrganizationID).String()
	todos, err := tc.getTodos(ctx, orgID)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"todos": todos}})
}

func (tc *TodoController) CreateTodo(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var todoInput models.TodoInput
	if err := ctx.ShouldBindJSON(&todoInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	todo := models.Todo{
		Name:           todoInput.Name,
		Description:    todoInput.Description,
		OrganizationID: todoInput.OrganizationID,
		SenderID:       currentUser.ID,
		State:          todoInput.State,
	}

	if result := tc.DB.Create(&todo); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error})
		return
	}
	orgID := (*todoInput.OrganizationID).String()
	todos, err := tc.getTodos(ctx, orgID)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"todos": todos}})
}

func (tc *TodoController) EditTodo(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	id := ctx.Param("id")
	var todo models.Todo
	var todoInput models.TodoEditInput

	if err := ctx.ShouldBindJSON(&todoInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if result := tc.DB.Where("id = ?", id).First(&todo); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "No todo found"})
		return
	}

	var todoChanges models.TodoChanges

	todoChanges.ManagerID = currentUser.ID
	if todo.Name != todoInput.Name {
		todoChanges.ChangedField = "name"
		todoChanges.OldValue = todo.Name
		todoChanges.NewValue = todoInput.Name
		todo.Name = todoInput.Name
	}

	if todo.Description != todoInput.Description {
		todoChanges.ToDoID = todo.ID
		todoChanges.ChangedField = "description"
		todoChanges.OldValue = todo.Description
		todoChanges.NewValue = todoInput.Description
		todo.Description = todoInput.Description
	}

	if todo.State != todoInput.State {
		todoChanges.ToDoID = todo.ID
		todoChanges.ChangedField = "state"
		todoChanges.OldValue = todo.State
		todoChanges.NewValue = todoInput.State
		todo.State = todoInput.State
	}

	if result := tc.DB.Create(&todoChanges); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error})
		return
	}

	if result := tc.DB.Save(&todo); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something went wrong"})
		return
	}

	orgID := (*todo.OrganizationID).String()
	todos, err := tc.getTodos(ctx, orgID)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"todos": todos}})
}

func (tc *TodoController) DeleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	var todo models.Todo
	if result := tc.DB.Where("id = ?", id).First(&todo); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Could not delete todo"})
		return
	}

	if result := tc.DB.Where("id = ?", id).Delete(&todo); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Could not delete todo"})
		return
	}

	orgID := (*todo.OrganizationID).String()
	todos, err := tc.getTodos(ctx, orgID)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"todos": todos}})
}

func (tc *TodoController) GetChangesOfToDo(ctx *gin.Context) {
	id := ctx.Param("id")
	var todoChanges []models.TodoChanges
	if result := tc.DB.Preload("Sender").Where("todo_id = ?", id).Find(&todoChanges); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Could not get todo history"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"todo_history": todoChanges}})
}

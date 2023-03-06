package dto

type CreateTodolistRequest struct {
	Title       string `json:"title" binding:"required,min=2"`
	Description string `json:"description" binding:"required,min=4"`
}

type UpdateTodolistRequest struct {
	Title       string `json:"title" binding:"required,min=2"`
	Description string `json:"description" binding:"required,min=4"`
	Status      bool   `json:"status"`
}

func (request *UpdateTodolistRequest) RequestUpdateTodolist() map[string]interface{} {
	update := make(map[string]interface{}, 0)
	if request.Title != "" {
		update["title"] = request.Title
	}

	if request.Description != "" {
		update["description"] = request.Description
	}

	update["status"] = request.Status

	return update
}

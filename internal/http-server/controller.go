package http_server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type TaskController interface {
	GetNextDate(log *slog.Logger)
	AddTask(log *slog.Logger, taskCreator TaskCreator)
	//GetTasks
	//GetTaskByID
	//TaskUpdate
	//TaskDone
	//TaskDelete
}

type Controller struct {
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) sendError(w http.ResponseWriter, r *http.Request, err error) {
	// TODO: error switch
	w.WriteHeader(http.StatusBadRequest)
	render.JSON(w, r, map[string]string{"error": err.Error()})
}

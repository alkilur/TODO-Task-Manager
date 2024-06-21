package http_server

import (
	"github.com/go-chi/render"
	"net/http"
)

type TaskController interface {
	GetNextDate(now, date, repeat string) (string, error)
	//GetTasks
	//GetTaskByID
	//TaskAdd
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

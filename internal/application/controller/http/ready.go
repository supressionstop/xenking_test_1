package http

import (
	"net/http"

	"github.com/supressionstop/xenking_test_1/internal/usecase"
)

type ReadyController struct {
	checker usecase.IsLineSyncedChecker
}

func NewReadyController(checker usecase.IsLineSyncedChecker) *ReadyController {
	return &ReadyController{
		checker: checker,
	}
}

func (c *ReadyController) Ready(w http.ResponseWriter, r *http.Request) {
	synced := c.checker.Execute()
	if synced {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

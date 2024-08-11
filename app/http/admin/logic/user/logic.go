package user

import (
	"github.com/gin-generator/ginctl/app/http/admin/logic"
	http "github.com/gin-generator/ginctl/package/respond"
	"github.com/gin-gonic/gin"
)

type Logic struct{}

func NewLogic() *Logic {
	return &Logic{}
}

// Index Get page list
func (l *Logic) Index(c *gin.Context) {
	params, err := logic.ParseAndCheckParams[Index](c)
	if err != nil {
		http.Alert400(c, http.StatusBadRequest, err.Error())
		return
	}
	// TODO: Your logic.

	// TODO: Replace your return struct.
	http.SuccessWithData(c, params.Data())
}

// Show Get info
func (l *Logic) Show(c *gin.Context) {
	params, err := logic.ParseAndCheckParams[Show](c)
	if err != nil {
		http.Alert400(c, http.StatusBadRequest, err.Error())
		return
	}
	// TODO: Your logic.

	// TODO: Replace your return struct.
	http.SuccessWithData(c, params.Data())
}

// Create Save one source
func (l *Logic) Create(c *gin.Context) {
	params, err := logic.ParseAndCheckParams[Create](c)
	if err != nil {
		http.Alert400(c, http.StatusBadRequest, err.Error())
		return
	}
	// TODO: Your logic.

	// TODO: Replace your return struct.
	http.SuccessWithData(c, params.Data())
}

// Update Modifying a resource
func (l *Logic) Update(c *gin.Context) {
	params, err := logic.ParseAndCheckParams[Update](c)
	if err != nil {
		http.Alert400(c, http.StatusBadRequest, err.Error())
		return
	}
	// TODO: Your logic.

	// TODO: Replace your return struct.
	http.SuccessWithData(c, params.Data())
}

// Destroy Delete a resource
func (l *Logic) Destroy(c *gin.Context) {
	params, err := logic.ParseAndCheckParams[Destroy](c)
	if err != nil {
		http.Alert400(c, http.StatusBadRequest, err.Error())
		return
	}
	// TODO: Your logic.

	// TODO: Replace your return struct.
	http.SuccessWithData(c, params.Data())
}

// Test test
func (l *Logic) Test(c *gin.Context) {
	params, err := logic.ParseAndCheckParams[Test](c)
	if err != nil {
		http.Alert400(c, http.StatusBadRequest, err.Error())
		return
	}
	// TODO: Your logic.

	// TODO: Replace your return struct.
	http.SuccessWithData(c, params.Data())
}

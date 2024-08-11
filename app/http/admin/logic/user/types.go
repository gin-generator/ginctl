package user

import (
	"github.com/gin-generator/ginctl/package/http"
	"github.com/gin-generator/ginctl/package/validator"
	"github.com/gin-gonic/gin"
)

// Demo demo.
//type Demo struct {
//	Page     int  `form:"page" validate:"numeric,min=1"`
//	Size     int  `form:"size" validate:"numeric,min=1,max=100"`
//	Keywords *string `form:"keywords" validate:"omitempty"`
//}
//
//func (r *Demo) ParseAndCheckParams(c *gin.Context) (err error) {
//	err = http.Parse(c, r)
//	if err != nil {
//		return
//	}
//	if err = validator.ValidateStructWithOutCtx(r); err != nil {
//		return
//	}
//	return
//}

type Index struct {
    // TODO: add your params.
}

func (r *Index) ParseAndCheckParams(c *gin.Context) (err error) {
	err = http.Parse(c, r)
	if err != nil {
		return
	}
	if err = validator.ValidateStructWithOutCtx(r); err != nil {
		return
	}
	// TODO: add your logic check.
	return
}

type Show struct {
    // TODO: add your params.
}

func (r *Show) ParseAndCheckParams(c *gin.Context) (err error) {
	err = http.Parse(c, r)
	if err != nil {
		return
	}
	if err = validator.ValidateStructWithOutCtx(r); err != nil {
		return
	}
	// TODO: add your logic check.
	return
}

type Create struct {
    // TODO: add your params.
}

func (r *Create) ParseAndCheckParams(c *gin.Context) (err error) {
	err = http.Parse(c, r)
	if err != nil {
		return
	}
	if err = validator.ValidateStructWithOutCtx(r); err != nil {
		return
	}
	// TODO: add your logic check.
	return
}

type Update struct {
    // TODO: add your params.
}

func (r *Update) ParseAndCheckParams(c *gin.Context) (err error) {
	err = http.Parse(c, r)
	if err != nil {
		return
	}
	if err = validator.ValidateStructWithOutCtx(r); err != nil {
		return
	}
	// TODO: add your logic check.
	return
}

type Destroy struct {
    // TODO: add your params.
}

func (r *Destroy) ParseAndCheckParams(c *gin.Context) (err error) {
	err = http.Parse(c, r)
	if err != nil {
		return
	}
	if err = validator.ValidateStructWithOutCtx(r); err != nil {
		return
	}
	// TODO: add your logic check.
	return
}

type Test struct {
    // TODO: add your params.
}

func (r *Test) ParseAndCheckParams(c *gin.Context) (err error) {
	err = http.Parse(c, r)
	if err != nil {
		return
	}
	if err = validator.ValidateStructWithOutCtx(r); err != nil {
		return
	}
	// TODO: add your logic check.
	return
}

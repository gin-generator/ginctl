package http

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

const (
	PATH = "path"
)

// RequestType is a generic struct that holds the parsed request data
type RequestType[T any] struct {
	data T
}

// NewRequestType is a constructor function for RequestType
func NewRequestType[T any](data T) RequestType[T] {
	return RequestType[T]{data: data}
}

// Data returns the parsed request data
func (r RequestType[T]) Data() T {
	return r.data
}

// Parse is a function that parses request parameters into the provided struct
func Parse[T any](c *gin.Context, obj *T) (err error) {

	// Handle the JSON binding first
	if err = c.ShouldBind(obj); err != nil {
		return
	}

	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag

		switch field.Kind() {
		case reflect.String:
			parseStringField(c, &field, tag)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			parseIntField(c, &field, tag)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			parseUintField(c, &field, tag)
		case reflect.Bool:
			parseBoolField(c, &field, tag)
		case reflect.Float32, reflect.Float64:
			parseFloatField(c, &field, tag)
		}
	}

	return
}

// Helper function to parse string fields
func parseStringField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) {
	if pathTag, ok := tag.Lookup(PATH); ok {
		field.SetString(c.Param(pathTag))
	}
}

// Helper function to parse integer fields
func parseIntField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) {
	if pathTag, ok := tag.Lookup(PATH); ok {
		if va, err := strconv.ParseInt(c.Param(pathTag), 10, 64); err == nil {
			field.SetInt(va)
		}
	}
}

// Helper function to parse unsigned integer fields
func parseUintField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) {
	if pathTag, ok := tag.Lookup(PATH); ok {
		if va, err := strconv.ParseUint(c.Param(pathTag), 10, 64); err == nil {
			field.SetUint(va)
		}
	}
}

// Helper function to parse boolean fields
func parseBoolField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) {
	if pathTag, ok := tag.Lookup(PATH); ok {
		if va, err := strconv.ParseBool(c.Param(pathTag)); err == nil {
			field.SetBool(va)
		}
	}
}

// Helper function to parse float fields
func parseFloatField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) {
	if pathTag, ok := tag.Lookup(PATH); ok {
		if va, err := strconv.ParseFloat(c.Param(pathTag), 64); err == nil {
			field.SetFloat(va)
		}
	}
}

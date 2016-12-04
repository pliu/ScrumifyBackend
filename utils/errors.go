package utils

import (
    "errors"
    "github.com/gin-gonic/gin"
)

var EmailExists error = errors.New("Email already exists")
var EmailDoesntExist error = errors.New("Email doesn't exist")
var UserDoesntExist error = errors.New("User doesn't exist")
var MappingDoesntExist error = errors.New("Mapping doesn't exist")
var StoryDoesntExist error = errors.New("Story doesn't exist")

var InternalErrorReturn gin.H = gin.H{"error": "Internal error"}
var UnauthorizedReturn gin.H = gin.H{"error": "Not authorized"}
var BadRequestReturn gin.H = gin.H{"error": "Required fields are empty"}

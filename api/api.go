package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"gin-framework/models"
	"gin-framework/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	router := gin.Default()
	router.POST("/user/create", CreateUser)
	router.GET("/user/all", GetAllUsers)
	router.DELETE("/user/delete", DeleteUser)
	router.GET("/user/get", GetUser)
	router.PUT("/user/update", UpdateUser)
	log.Println("Server is running...")
	if err := router.Run("localhost:8088"); err != nil {
		fmt.Println("Error while running server!")
	}
}

func CreateUser(c *gin.Context) {
	bodyByte, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("error while getting body", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var user *models.User
	if err = json.Unmarshal(bodyByte, &user); err != nil {
		log.Println("error while unmarshalling body", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	id := uuid.NewString()
	user.Id = id

	respUser, err := storage.CreateUser(user)
	if err != nil {
		log.Println("error while creating user", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusCreated, respUser)
}

func UpdateUser(c *gin.Context) {
	bodyByte, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("error while getting body", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var user *models.User
	if err = json.Unmarshal(bodyByte, &user); err != nil {
		log.Println("error while unmarshalling body", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	userId := c.Request.URL.Query().Get("id")

	respUser, err := storage.UpdateUser(userId, user)
	if err != nil {
		log.Println("error while updating user", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, respUser)
}

func DeleteUser(c *gin.Context) {
	userId := c.Request.URL.Query().Get("id")

	if err := storage.DeleteUser(userId); err != nil {
		log.Println("error while deleting user", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, "Deleted User")
}

func GetUser(c *gin.Context) {
	userId := c.Request.URL.Query().Get("id")

	respUser, err := storage.GetUser(userId)
	if err != nil {
		log.Println("Error while getting user", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, respUser)
}

func GetAllUsers(ctx *gin.Context) {
	page := ctx.Request.URL.Query().Get("page")
	intPage, err := strconv.Atoi(page)
	if err != nil {
		log.Println("Error while converting page")
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	limit := ctx.Request.URL.Query().Get("limit")

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("Error while converting limit")
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	users, err := storage.GetAllUsers(intPage, intLimit)
	if err != nil {
		log.Println("Error while getting all users", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, users)
}

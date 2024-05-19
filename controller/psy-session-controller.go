package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/putukrisna6/golang-api/dto"
	"github.com/putukrisna6/golang-api/entity"
	"github.com/putukrisna6/golang-api/helper"
	"github.com/putukrisna6/golang-api/service"
)

type PsySessionController interface {
	All(context *gin.Context)
	Get(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type psySessionController struct {
	psySessionService service.PsySessionService
	jwtService        service.JWTService
}

func NewPsySessionController(psySessionService service.PsySessionService, jwtService service.JWTService) PsySessionController {
	return &psySessionController{
		psySessionService: psySessionService,
		jwtService:        jwtService,
	}
}

func (c *psySessionController) All(context *gin.Context) {
	var psySessions []entity.PsySession = c.psySessionService.All()
	res := helper.BuildValidResponse("OK", psySessions)
	context.JSON(http.StatusOK, res)
}

func (c *psySessionController) Get(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("parameter ID must not be empty", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var session entity.PsySession = c.psySessionService.Get(id)
	if (session == entity.PsySession{}) {
		res := helper.BuildErrorResponse("failed to retrieve Book", "no data with given bookID", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	res := helper.BuildValidResponse("OK", session)
	context.JSON(http.StatusOK, res)
}

func (c *psySessionController) Insert(context *gin.Context) {
	var PsySessionCreateDTO dto.PsySessionCreateDTO
	errDTO := context.ShouldBind(&PsySessionCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	userID := c.getUserIDByToken(authHeader)
	convertedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		panic(err.Error())
	}

	PsySessionCreateDTO.ClientID = convertedUserID
	result := c.psySessionService.Insert(PsySessionCreateDTO)
	response := helper.BuildValidResponse("OK", result)
	context.JSON(http.StatusCreated, response)
}

func (c *psySessionController) Update(context *gin.Context) {
	var psySessionUpdateDTO dto.PsySessionUpdateDTO
	errDTO := context.ShouldBind(&psySessionUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.psySessionService.IsAllowedToEdit(userID, psySessionUpdateDTO.UserID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID != nil {
			panic(errID.Error())
		}

		psySessionUpdateDTO.UserID = id
		result := c.psySessionService.Update(psySessionUpdateDTO)
		response := helper.BuildValidResponse("OK", result)
		context.JSON(http.StatusOK, response)
		return
	}

	response := helper.BuildErrorResponse("you do not have permission to update this Book", "you are not the owner", helper.EmptyObj{})
	context.AbortWithStatusJSON(http.StatusForbidden, response)
}

func (c *psySessionController) Delete(context *gin.Context) {
	var psySession entity.PsySession
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("parameter ID must not be empty", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	psySession.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.psySessionService.IsAllowedToEdit(userID, psySession.ID) {
		c.psySessionService.Delete(psySession)
		message := fmt.Sprintf("Book with ID %v successfuly deleted", psySession.ID)
		res := helper.BuildValidResponse(message, helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
		return
	}

	response := helper.BuildErrorResponse("you do not have permission to delete this Book", "you are not the owner", helper.EmptyObj{})
	context.AbortWithStatusJSON(http.StatusForbidden, response)
}

func (c *psySessionController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}

	claims := aToken.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["user_id"])
}

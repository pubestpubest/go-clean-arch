package delivery

import (
	"net/http"
	"order-management/domain"
	"order-management/entity"
	"order-management/utils"
	"strconv"

	"order-management/middleware"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	userUsecase  domain.UserUsecase
	orderUsecase domain.OrderUsecase
}

func NewHandler(e *echo.Group, u domain.UserUsecase, o domain.OrderUsecase) *Handler {
	h := Handler{
		userUsecase:  u,
		orderUsecase: o,
	}

	publicGroup := e.Group("")
	publicGroup.POST("/register", h.CreateUser)
	publicGroup.POST("/login", h.Login)
	publicGroup.GET("/:id", h.GetUserByID)

	authGroup := e.Group("")
	authGroup.Use(middleware.UserAuth())
	authGroup.PUT("/:id", h.UpdateUser)
	// authGroup.GET("/orders", h.GetOrdersByUserID)
	authGroup.POST("/orders", h.CreateOrder)
	return &h
}

func (h *Handler) CreateOrder(c echo.Context) error {
	log.Trace("Entering function CreateOrder()")
	defer log.Trace("Exiting function CreateOrder()")

	req := entity.OrderRequest{}

	if err := c.Bind(&req); err != nil {
		err = errors.Wrap(err, "[Handler.CreateOrder]: invalid order data")

		log.WithError(err).Warn("Invalid order data")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: utils.StandardError(err)})
	}

	userID := c.Get("user").(*entity.UserJWT).ID

	if err := h.orderUsecase.CreateOrder(req, userID); err != nil {
		err = errors.Wrap(err, "[Handler.CreateOrder]: internal server error")

		log.WithFields(log.Fields{
			"order": req,
		}).WithError(err).Error("Internal server error during order creation")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: utils.StandardError(err)})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Order created successfully",
		Status:  http.StatusOK,
	})
}

func (h *Handler) Login(c echo.Context) error {
	//Bind
	req := entity.User{}
	if err := c.Bind(&req); err != nil {
		err = errors.Wrap(err, "[Handler.Login]: failed to bind request")

		log.WithError(err).Warn("Failed to bind request")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: utils.StandardError(err)})
	}
	//Check if email and password are provided
	if req.Email == "" || req.Password == "" {
		err := errors.New("[Handler.Login]: email and password are required")

		log.Warn("Email and password are required")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: utils.StandardError(err)})
	}
	//Login
	user, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		if err.Error() == "[UserUsecase.Login]: user not found" {
			err = errors.Wrap(err, "[Handler.Login]: user not found")

			log.WithFields(log.Fields{
				"email": req.Email,
			}).WithError(err).Warn("User not found during login")

			return c.JSON(http.StatusNotFound, entity.ResponseError{Error: utils.StandardError(err)})
		}
		if err.Error() == "[UserUsecase.Login]: invalid password" {
			err = errors.Wrap(err, "[Handler.Login]: invalid password")

			log.WithFields(log.Fields{
				"email": req.Email,
			}).WithError(err).Warn("Invalid password during login")

			return c.JSON(http.StatusUnauthorized, entity.ResponseError{Error: utils.StandardError(err)})
		}
		err = errors.Wrap(err, "[Handler.Login]: internal server error")

		log.WithFields(log.Fields{
			"email": req.Email,
		}).WithError(err).Error("Internal server error during login")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: utils.StandardError(err)})
	}

	c.Response().Header().Set("Authorization", "Bearer "+user)

	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Login successful",
		Status:  http.StatusOK,
	})
}

func (h *Handler) CreateUser(c echo.Context) error {
	req := entity.User{}

	if err := c.Bind(&req); err != nil {
		err = errors.Wrap(err, "[Handler.CreateUser]: invalid user data")

		log.WithError(err).Warn("Invalid user data")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: utils.StandardError(err)})
	}
	if req.Email == "" || req.Password == "" {
		err := errors.New("[Handler.CreateUser]: email and password are required")

		log.Warn("Email and password are required")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: utils.StandardError(err)})
	}

	if err := h.userUsecase.CreateUser(req); err != nil {
		if err.Error() == "[UserUsecase.CreateUser]: user already exists" {
			err = errors.Wrap(err, "[Handler.CreateUser]: user already exists")

			log.WithFields(log.Fields{
				"email": req.Email,
			}).WithError(err).Warn("User already exists during creation")

			return c.JSON(http.StatusConflict, entity.ResponseError{Error: utils.StandardError(err)})
		}
		err = errors.Wrap(err, "[Handler.CreateUser]: internal server error")

		log.WithFields(log.Fields{
			"email": req.Email,
		}).WithError(err).Error("Internal server error during user creation")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: utils.StandardError(err)})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "User created successfully",
		Status:  http.StatusOK,
	})
}

func (h *Handler) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetUserByID]: invalid user id")

		log.WithError(err).Warn("Invalid user id")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: utils.StandardError(err)})
	}

	user, err := h.userUsecase.GetUserByID(uint32(id))
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetUserByID]: internal server error")

		log.WithFields(log.Fields{
			"user_id": id,
		}).WithError(err).Error("Internal server error during user retrieval")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: utils.StandardError(err)})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = errors.Wrap(err, "[Handler.UpdateUser]: invalid user id")

		log.WithError(err).Warn("Invalid user id")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: utils.StandardError(err)})
	}

	user, err := h.userUsecase.GetUserByID(uint32(id))
	if err != nil {
		if err.Error() == "[UserUsecase.GetUserByID]: user not found" {
			err = errors.Wrap(err, "[Handler.UpdateUser]: user not found")

			log.WithFields(log.Fields{
				"user_id": id,
			}).WithError(err).Warn("User not found during user retrieval")

			return c.JSON(http.StatusNotFound, entity.ResponseError{Error: utils.StandardError(err)})
		}
		err = errors.Wrap(err, "[Handler.UpdateUser]: internal server error")

		log.WithFields(log.Fields{
			"user_id": id,
		}).WithError(err).Error("Internal server error during user retrieval")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: utils.StandardError(err)})
	}

	if err := c.Bind(&user); err != nil {
		err = errors.Wrap(err, "[Handler.UpdateUser]: invalid user data")

		log.WithError(err).Warn("Invalid user data")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: utils.StandardError(err)})
	}

	if err := h.userUsecase.UpdateUser(user); err != nil {
		err = errors.Wrap(err, "[Handler.UpdateUser]: internal server error")

		log.WithFields(log.Fields{
			"user_id": id,
		}).WithError(err).Error("Internal server error during user update")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: utils.StandardError(err)})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "User updated successfully",
		Status:  http.StatusOK,
	})
}

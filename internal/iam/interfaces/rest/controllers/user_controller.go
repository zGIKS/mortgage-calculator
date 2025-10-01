package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"finanzas-backend/internal/iam/domain/model/commands"
	"finanzas-backend/internal/iam/domain/model/queries"
	"finanzas-backend/internal/iam/domain/model/entities"
	"finanzas-backend/internal/iam/domain/services"
	"finanzas-backend/internal/iam/interfaces/rest/resources"
)

type UserController struct {
	userCommandService services.UserCommandService
	userQueryService   services.UserQueryService
	authService        services.AuthenticationService
}

func NewUserController(
	userCommandService services.UserCommandService,
	userQueryService services.UserQueryService,
	authService services.AuthenticationService,
) *UserController {
	return &UserController{
		userCommandService: userCommandService,
		userQueryService:   userQueryService,
		authService:        authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, password, and full name
// @Tags IAM
// @Accept json
// @Produce json
// @Param request body resources.RegisterUserResource true "User registration request"
// @Success 201 {object} resources.UserResource
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/iam/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req resources.RegisterUserResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd, err := commands.NewRegisterUserCommand(req.Email, req.Password, req.FullName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := c.userCommandService.HandleRegister(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	query, _ := queries.NewFindUserByEmailQuery(req.Email)
	user, err := c.userQueryService.HandleFindByEmail(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created user"})
		return
	}

	_ = userID
	response := c.transformUserToResource(user)
	ctx.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return token
// @Tags IAM
// @Accept json
// @Produce json
// @Param request body resources.LoginResource true "Login credentials"
// @Success 200 {object} resources.LoginResponseResource
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/iam/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req resources.LoginResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd, err := commands.NewLoginCommand(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authService.HandleLogin(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get user details
	query, _ := queries.NewFindUserByEmailQuery(req.Email)
	user, err := c.userQueryService.HandleFindByEmail(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	response := resources.LoginResponseResource{
		Token: token,
		User:  c.transformUserToResource(user),
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) transformUserToResource(user *entities.User) resources.UserResource {
	return resources.UserResource{
		ID:        user.ID().String(),
		Email:     user.Email().Value(),
		FullName:  user.FullName(),
		CreatedAt: user.CreatedAt(),
	}
}

package controllers

import (
	"net/http"

	"finanzas-backend/internal/profile/domain/model/commands"
	"finanzas-backend/internal/profile/domain/model/entities"
	"finanzas-backend/internal/profile/domain/model/queries"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
	"finanzas-backend/internal/profile/domain/services"
	"finanzas-backend/internal/profile/interfaces/rest/resources"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	commandService services.ProfileCommandService
	queryService   services.ProfileQueryService
}

func NewProfileController(
	commandService services.ProfileCommandService,
	queryService services.ProfileQueryService,
) *ProfileController {
	return &ProfileController{
		commandService: commandService,
		queryService:   queryService,
	}
}

// GetProfile godoc
// @Summary Get profile
// @Description Get authenticated user's profile
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} resources.ProfileResource
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/profile [get]
func (c *ProfileController) GetProfile(ctx *gin.Context) {
	// Get user_id from context (set by JWT middleware)
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := valueobjects.NewUserIDFromString(userIDValue.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query profile by user ID
	query, _ := queries.NewFindProfileByUserIDQuery(userID.String())
	profile, err := c.queryService.HandleFindByUserID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	response := c.transformProfileToResource(profile)
	ctx.JSON(http.StatusOK, response)
}

// UpdateProfile godoc
// @Summary Update profile
// @Description Update authenticated user's profile
// @Tags Profile
// @Accept json
// @Produce json
// @Param request body resources.UpdateProfileResource true "Update profile request"
// @Success 200 {object} resources.ProfileResource
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/profile [put]
func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	// Get user_id from context (set by JWT middleware)
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := valueobjects.NewUserIDFromString(userIDValue.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get profile to obtain profile ID
	queryProfile, _ := queries.NewFindProfileByUserIDQuery(userID.String())
	profile, err := c.queryService.HandleFindByUserID(ctx.Request.Context(), queryProfile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	var req resources.UpdateProfileResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create command
	cmd, err := commands.NewUpdateProfileCommand(
		profile.ID(),
		req.PhoneNumber,
		req.MonthlyIncome,
		req.Currency,
		req.MaritalStatus,
		req.IsFirstHome,
		req.HasOwnLand,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute command
	if err := c.commandService.HandleUpdate(ctx.Request.Context(), cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve updated profile
	query, _ := queries.NewFindProfileByIDQuery(profile.ID().String())
	updatedProfile, err := c.queryService.HandleFindByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated profile"})
		return
	}

	response := c.transformProfileToResource(updatedProfile)
	ctx.JSON(http.StatusOK, response)
}

func (c *ProfileController) transformProfileToResource(profile *entities.Profile) resources.ProfileResource {
	return resources.ProfileResource{
		ID:             profile.ID().String(),
		UserID:         profile.UserID().String(),
		DNI:            profile.DNI().Value(),
		FirstName:      profile.FirstName(),
		FirstLastName:  profile.FirstLastName(),
		SecondLastName: profile.SecondLastName(),
		FullName:       profile.FullName(),
		PhoneNumber:    profile.PhoneNumber().Value(),
		MonthlyIncome:  profile.MonthlyIncome().Amount(),
		Currency:       string(profile.MonthlyIncome().Currency()),
		MaritalStatus:  profile.MaritalStatus().String(),
		IsFirstHome:    profile.IsFirstHome(),
		HasOwnLand:     profile.HasOwnLand(),
		CreatedAt:      profile.CreatedAt(),
	}
}

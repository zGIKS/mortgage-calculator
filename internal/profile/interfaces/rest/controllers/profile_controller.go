package controllers

import (
	"net/http"

	"finanzas-backend/internal/profile/domain/model/commands"
	"finanzas-backend/internal/profile/domain/model/entities"
	"finanzas-backend/internal/profile/domain/model/queries"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
	"finanzas-backend/internal/profile/domain/services"
	"finanzas-backend/internal/profile/infrastructure/external"
	"finanzas-backend/internal/profile/interfaces/rest/resources"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	commandService services.ProfileCommandService
	queryService   services.ProfileQueryService
	reniecService  *external.ReniecService
}

func NewProfileController(
	commandService services.ProfileCommandService,
	queryService services.ProfileQueryService,
	reniecService *external.ReniecService,
) *ProfileController {
	return &ProfileController{
		commandService: commandService,
		queryService:   queryService,
		reniecService:  reniecService,
	}
}

// GetReniecData godoc
// @Summary Get person data from RENIEC by DNI
// @Description Consult RENIEC API to get person information by DNI
// @Tags Profile
// @Accept json
// @Produce json
// @Param dni query string true "DNI number (8 digits)"
// @Success 200 {object} resources.ReniecDataResource
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/profile/reniec [get]
func (c *ProfileController) GetReniecData(ctx *gin.Context) {
	dni := ctx.Query("dni")
	if dni == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "DNI is required"})
		return
	}

	// Validate DNI format
	if _, err := valueobjects.NewDNI(dni); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get data from RENIEC
	personData, err := c.reniecService.GetPersonDataByDNI(ctx.Request.Context(), dni)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := resources.ReniecDataResource{
		DNI:            personData.DocumentNumber,
		FirstName:      personData.FirstName,
		FirstLastName:  personData.FirstLastName,
		SecondLastName: personData.SecondLastName,
		FullName:       personData.FullName,
	}

	ctx.JSON(http.StatusOK, response)
}

// CreateProfile godoc
// @Summary Create a new profile
// @Description Create a new profile for authenticated user with DNI validation
// @Tags Profile
// @Accept json
// @Produce json
// @Param request body resources.CreateProfileResource true "Profile creation request"
// @Success 201 {object} resources.ProfileResource
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/profile [post]
func (c *ProfileController) CreateProfile(ctx *gin.Context) {
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

	var req resources.CreateProfileResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create command
	cmd, err := commands.NewCreateProfileCommand(
		userID,
		req.DNI,
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
	profileID, err := c.commandService.HandleCreate(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Retrieve created profile
	query, _ := queries.NewFindProfileByIDQuery(profileID.String())
	profile, err := c.queryService.HandleFindByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created profile"})
		return
	}

	response := c.transformProfileToResource(profile)
	ctx.JSON(http.StatusCreated, response)
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

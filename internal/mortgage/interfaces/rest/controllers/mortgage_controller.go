package controllers

import (
	"net/http"
	"strconv"

	"finanzas-backend/internal/mortgage/domain/model/commands"
	"finanzas-backend/internal/mortgage/domain/model/queries"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
	"finanzas-backend/internal/mortgage/domain/services"
	"finanzas-backend/internal/mortgage/interfaces/rest/resources"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MortgageController struct {
	commandService services.MortgageCommandService
	queryService   services.MortgageQueryService
}

func NewMortgageController(
	commandService services.MortgageCommandService,
	queryService services.MortgageQueryService,
) *MortgageController {
	return &MortgageController{
		commandService: commandService,
		queryService:   queryService,
	}
}

// CalculateMortgage godoc
// @Summary Calculate mortgage with French method
// @Description Calculates a mortgage loan using the French amortization method (constant installments). The bank_id is required and determines the rate type, payment frequency, and days in year.
// @Tags Mortgage
// @Accept json
// @Produce json
// @Param request body resources.CalculateMortgageRequest true "Mortgage calculation request"
// @Success 200 {object} resources.MortgageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/mortgage/calculate [post]
func (c *MortgageController) CalculateMortgage(ctx *gin.Context) {
	var req resources.CalculateMortgageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener user_id del contexto (guardado por el middleware JWT)
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID := userIDValue.(string)

	// Parse BankID (required)
	parsed, err := uuid.Parse(req.BankID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank ID format"})
		return
	}
	bankID := &parsed

	cmd, err := commands.NewCalculateMortgageCommand(
		userID,
		req.PrecioVenta,
		req.CuotaInicial,
		req.MontoPrestamo,
		req.BonoTechoPropio,
		req.TEA,
		bankID,
		req.PlazoMeses,
		req.MesesGracia,
		req.TipoGracia,
		req.Moneda,
		req.TasaDescuento,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mortgage, err := c.commandService.HandleCalculateMortgage(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := resources.TransformToMortgageResponse(mortgage)
	ctx.JSON(http.StatusOK, response)
}

// GetMortgageByID godoc
// @Summary Get mortgage by ID
// @Description Get a specific mortgage calculation by ID
// @Tags Mortgage
// @Accept json
// @Produce json
// @Param id path uint64 true "Mortgage ID"
// @Success 200 {object} resources.MortgageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/mortgage/{id} [get]
func (c *MortgageController) GetMortgageByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid mortgage ID"})
		return
	}

	query, err := queries.NewGetMortgageByIDQuery(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mortgage, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := resources.TransformToMortgageResponse(mortgage)
	ctx.JSON(http.StatusOK, response)
}

// GetMortgageHistory godoc
// @Summary Get mortgage calculation history
// @Description Get mortgage calculation history for authenticated user
// @Tags Mortgage
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} resources.MortgageSummaryResource
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/mortgage/history [get]
func (c *MortgageController) GetMortgageHistory(ctx *gin.Context) {
	// Obtener user_id del contexto (guardado por el middleware JWT)
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID := userIDValue.(string)

	query, err := queries.NewGetMortgageHistoryQuery(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse pagination
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			if offsetStr := ctx.Query("offset"); offsetStr != "" {
				if offset, err := strconv.Atoi(offsetStr); err == nil {
					query = query.WithPagination(limit, offset)
				}
			} else {
				query = query.WithPagination(limit, 0)
			}
		}
	}

	mortgages, err := c.queryService.HandleGetHistory(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]resources.MortgageSummaryResource, 0, len(mortgages))
	for _, mortgage := range mortgages {
		response = append(response, resources.TransformToMortgageSummary(mortgage))
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateMortgage godoc
// @Summary Update mortgage
// @Description Update an existing mortgage calculation. This will recalculate all values.
// @Tags Mortgage
// @Accept json
// @Produce json
// @Param id path uint64 true "Mortgage ID"
// @Param request body resources.UpdateMortgageRequest true "Mortgage update request"
// @Success 200 {object} resources.MortgageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/mortgage/{id} [put]
func (c *MortgageController) UpdateMortgage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid mortgage ID"})
		return
	}

	mortgageID, err := valueobjects.NewMortgageID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req resources.UpdateMortgageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse BankID if provided
	var bankID *uuid.UUID
	if req.BankID != nil {
		parsed, err := uuid.Parse(*req.BankID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank ID format"})
			return
		}
		bankID = &parsed
	}

	cmd, err := commands.NewUpdateMortgageCommand(
		mortgageID,
		req.PrecioVenta,
		req.CuotaInicial,
		req.MontoPrestamo,
		req.BonoTechoPropio,
		req.TEA,
		req.TipoTasa,
		bankID,
		req.FrecuenciaPago,
		req.DiasAnio,
		req.PlazoMeses,
		req.MesesGracia,
		req.TipoGracia,
		req.Moneda,
		req.TasaDescuento,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mortgage, err := c.commandService.HandleUpdateMortgage(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := resources.TransformToMortgageResponse(mortgage)
	ctx.JSON(http.StatusOK, response)
}

// DeleteMortgage godoc
// @Summary Delete mortgage
// @Description Delete a mortgage calculation by ID
// @Tags Mortgage
// @Accept json
// @Produce json
// @Param id path uint64 true "Mortgage ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/mortgage/{id} [delete]
func (c *MortgageController) DeleteMortgage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid mortgage ID"})
		return
	}

	mortgageID, err := valueobjects.NewMortgageID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener user_id del contexto (guardado por el middleware JWT)
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := valueobjects.NewUserID(userIDValue.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd, err := commands.NewDeleteMortgageCommand(mortgageID, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.commandService.HandleDeleteMortgage(ctx.Request.Context(), cmd); err != nil {
		if err.Error() == "unauthorized access to mortgage" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "mortgage not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

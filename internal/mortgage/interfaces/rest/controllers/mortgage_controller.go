package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"finanzas-backend/internal/mortgage/domain/model/commands"
	"finanzas-backend/internal/mortgage/domain/model/queries"
	"finanzas-backend/internal/mortgage/domain/services"
	"finanzas-backend/internal/mortgage/interfaces/rest/resources"
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
// @Description Calculates a mortgage loan using the French amortization method (constant installments)
// @Tags Mortgage
// @Accept json
// @Produce json
// @Param request body resources.CalculateMortgageRequest true "Mortgage calculation request"
// @Success 200 {object} resources.MortgageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/mortgage/calculate [post]
func (c *MortgageController) CalculateMortgage(ctx *gin.Context) {
	var req resources.CalculateMortgageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd, err := commands.NewCalculateMortgageCommand(
		req.UserID,
		req.PropertyPrice,
		req.DownPayment,
		req.LoanAmount,
		req.BonoTechoPropio,
		req.InterestRate,
		req.RateType,
		req.TermMonths,
		req.GracePeriodMonths,
		req.GracePeriodType,
		req.Currency,
		req.NPVDiscountRate,
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
// @Description Get mortgage calculation history for a user
// @Tags Mortgage
// @Accept json
// @Produce json
// @Param user_id query uint64 true "User ID"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} resources.MortgageSummaryResource
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/mortgage/history [get]
func (c *MortgageController) GetMortgageHistory(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

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

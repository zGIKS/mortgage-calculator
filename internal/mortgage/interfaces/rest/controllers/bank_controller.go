package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"finanzas-backend/internal/mortgage/domain/model/queries"
	"finanzas-backend/internal/mortgage/domain/services"
	"finanzas-backend/internal/mortgage/interfaces/rest/resources"
)

type BankController struct {
	queryService services.BankQueryService
}

func NewBankController(queryService services.BankQueryService) *BankController {
	return &BankController{
		queryService: queryService,
	}
}

// GetAllBanks godoc
// @Summary Get all banks
// @Description Get all available banks with their configuration
// @Tags Banks
// @Accept json
// @Produce json
// @Success 200 {array} resources.BankResource
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/banks [get]
func (c *BankController) GetAllBanks(ctx *gin.Context) {
	query := queries.NewGetAllBanksQuery()

	banks, err := c.queryService.HandleGetAll(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := resources.TransformToBankResources(banks)
	ctx.JSON(http.StatusOK, response)
}

// GetBankByID godoc
// @Summary Get bank by ID
// @Description Get a specific bank by ID
// @Tags Banks
// @Accept json
// @Produce json
// @Param id path string true "Bank ID (UUID)"
// @Success 200 {object} resources.BankResource
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/banks/{id} [get]
func (c *BankController) GetBankByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	bankID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank ID format"})
		return
	}

	query, err := queries.NewGetBankByIDQuery(bankID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bank, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := resources.TransformToBankResource(bank)
	ctx.JSON(http.StatusOK, response)
}

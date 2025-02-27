package handler

import (
	"go-playground/server/domain"
	"go-playground/server/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RewardsHandler struct {
	rewardsService domain.RewardsService
}

func NewRewardsHandler(rewardsService domain.RewardsService) *RewardsHandler {
	return &RewardsHandler{rewardsService: rewardsService}
}

// @Summary Create reward
// @Description Create a new reward
// @Tags rewards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param reward body domain.Reward true "Reward details"
// @Success 201 {object} domain.Reward
// @Failure 400 {object} map[string]string
// @Router /rewards [post]
func (h *RewardsHandler) Create(c *gin.Context) {
	var req domain.CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	reward, err := h.rewardsService.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, reward)
}

// @Summary Get reward by ID
// @Description Get reward details by ID
// @Tags rewards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Reward ID"
// @Success 200 {object} domain.Reward
// @Failure 404 {object} map[string]string
// @Router /rewards/{id} [get]
func (h *RewardsHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid reward ID",
		})
		return
	}

	reward, err := h.rewardsService.GetByID(c.Request.Context(), id)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, reward)
}

// @Summary Get all rewards
// @Description Get all available rewards
// @Tags rewards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param active query bool false "Filter active rewards only"
// @Success 200 {array} domain.Reward
// @Router /rewards [get]
func (h *RewardsHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// @Summary Update reward
// @Description Update reward details
// @Tags rewards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Reward ID"
// @Param reward body domain.Reward true "Updated reward details"
// @Success 200 {object} domain.Reward
// @Failure 400,404 {object} map[string]string
// @Router /rewards/{id} [put]
func (h *RewardsHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid reward ID",
		})
		return
	}

	var req domain.UpdateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	reward, err := h.rewardsService.Update(c.Request.Context(), id, &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, reward)
}

// @Summary Delete reward
// @Description Delete a reward
// @Tags rewards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Reward ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /rewards/{id} [delete]
func (h *RewardsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid reward ID",
		})
		return
	}

	if err := h.rewardsService.Delete(c.Request.Context(), id); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reward deleted successfully"})
}

// @Summary Get rewards by program ID
// @Description Get all rewards associated with a specific program
// @Tags rewards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param program_id path string true "Program ID"
// @Success 200 {array} domain.Reward
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /rewards/program/{program_id} [get]
func (h *RewardsHandler) GetByProgramID(c *gin.Context) {
	programIDStr := c.Param("program_id")
	programID, err := uuid.Parse(programIDStr)
	if err != nil {
		util.HandleError(c, domain.ValidationError{
			Field:   "program_id",
			Message: "invalid program ID",
		})
		return
	}

	rewards, err := h.rewardsService.GetByProgramID(c.Request.Context(), programID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, rewards)
}

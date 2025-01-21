package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RewardsHandler struct {
	rewardsService *service.RewardsService
}

func NewRewardsHandler(rewardsService *service.RewardsService) *RewardsHandler {
	return &RewardsHandler{rewardsService: rewardsService}
}

// @Summary Create reward
// @Description Create a new reward
// @Tags rewards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reward body domain.Reward true "Reward details"
// @Success 201 {object} domain.Reward
// @Failure 400 {object} map[string]string
// @Router /rewards [post]
func (h *RewardsHandler) Create(c *gin.Context) {
	var reward domain.Reward
	if err := c.ShouldBindJSON(&reward); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.rewardsService.Create(&reward); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
// @Param id path string true "Reward ID"
// @Success 200 {object} domain.Reward
// @Failure 404 {object} map[string]string
// @Router /rewards/{id} [get]
func (h *RewardsHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	reward, err := h.rewardsService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
// @Param active query bool false "Filter active rewards only"
// @Success 200 {array} domain.Reward
// @Router /rewards [get]
func (h *RewardsHandler) GetAll(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "false") == "true"
	rewards, err := h.rewardsService.GetAll(activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rewards)
}

// @Summary Update reward
// @Description Update reward details
// @Tags rewards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Reward ID"
// @Param reward body domain.Reward true "Updated reward details"
// @Success 200 {object} domain.Reward
// @Failure 400,404 {object} map[string]string
// @Router /rewards/{id} [put]
func (h *RewardsHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var reward domain.Reward
	if err := c.ShouldBindJSON(&reward); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reward.ID = id

	if err := h.rewardsService.Update(&reward); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
// @Param id path string true "Reward ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /rewards/{id} [delete]
func (h *RewardsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.rewardsService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reward deleted successfully"})
}

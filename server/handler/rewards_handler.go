package handler

import (
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"go-playground/server/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type RewardsHandler struct {
	rewardsService domain.RewardsService
	logger         zerolog.Logger
}

func NewRewardsHandler(rewardsService domain.RewardsService) *RewardsHandler {
	return &RewardsHandler{
		rewardsService: rewardsService,
		logger:         logging.GetLogger(),
	}
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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming create reward request")

	var req domain.CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind create reward request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	reward, err := h.rewardsService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Interface("request", req).
			Msg("Failed to create reward")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("reward_id", reward.ID.String()).
		Str("program_id", reward.ProgramID.String()).
		Msg("Reward created successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get reward request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing reward ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid reward ID",
		})
		return
	}

	reward, err := h.rewardsService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("reward_id", id).
			Msg("Failed to get reward")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("reward_id", reward.ID.String()).
		Str("program_id", reward.ProgramID.String()).
		Msg("Reward retrieved successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get all rewards request")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming update reward request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing reward ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid reward ID",
		})
		return
	}

	var req domain.UpdateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind update reward request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	reward, err := h.rewardsService.Update(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("reward_id", id).
			Interface("request", req).
			Msg("Failed to update reward")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("reward_id", reward.ID.String()).
		Str("program_id", reward.ProgramID.String()).
		Msg("Reward updated successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming delete reward request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing reward ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid reward ID",
		})
		return
	}

	if err := h.rewardsService.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error().
			Err(err).
			Str("reward_id", id).
			Msg("Failed to delete reward")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("reward_id", id).
		Msg("Reward deleted successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get rewards by program request")

	programIDStr := c.Param("program_id")
	programID, err := uuid.Parse(programIDStr)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("program_id", programIDStr).
			Msg("Invalid program ID format")
		util.HandleError(c, domain.ValidationError{
			Field:   "program_id",
			Message: "invalid program ID",
		})
		return
	}

	rewards, err := h.rewardsService.GetByProgramID(c.Request.Context(), programID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("program_id", programID.String()).
			Msg("Failed to get rewards")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("program_id", programID.String()).
		Int("rewards_count", len(rewards)).
		Msg("Rewards retrieved successfully")

	c.JSON(http.StatusOK, rewards)
}

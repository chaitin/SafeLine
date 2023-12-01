package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chaitin/SafeLine/internal/service"
)

// GitHubHandler provides HTTP handlers for GitHub operations.
type GitHubHandler struct {
	gitHubService *service.GitHubService
}

// NewGitHubHandler creates a new instance of GitHubHandler.
func NewGitHubHandler(gitHubService *service.GitHubService) *GitHubHandler {
	return &GitHubHandler{
		gitHubService: gitHubService,
	}
}

// GetIssues handles GET requests for fetching GitHub issues.
// @Summary get issues
// @Description get issues from GitHub
// @Tags	GitHub
// @Accept	json
// @Produce	json
// @Param	q		query	string 	false  "search by"
// @Success 200 {array} service.Issue
// @Router /repos/issues [get]
func (h *GitHubHandler) GetIssues(c *gin.Context) {
	filter := c.Query("q")

	issues, err := h.gitHubService.GetIssues(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, issues)
}

// GetDiscussions handles GET requests for fetching GitHub discussions.
// @Summary get discussions
// @Description get discussions from GitHub
// @Tags	GitHub
// @Accept	json
// @Produce	json
// @Param	q		query	string 	false  "search by"
// @Success 200 {array} service.Discussion
// @Router /repos/discussions [get]
func (h *GitHubHandler) GetDiscussions(c *gin.Context) {
	filter := c.Query("q")

	discussions, err := h.gitHubService.GetDiscussions(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discussions)
}

// GetRepo handles GET requests for fetching GitHub repo info.
// @Summary get repo info
// @Description get repo info from GitHub
// @Tags	GitHub
// @Accept	json
// @Produce	json
// @Success 200 {object} service.Repo
// @Router /repos/info [get]
func (h *GitHubHandler) GetRepo(c *gin.Context) {
	repo, err := h.gitHubService.GetRepo(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, repo)
}

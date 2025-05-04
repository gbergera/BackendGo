package api

import (
	"net/http"
	"strconv"

	follow "ualabackend/entities/follow"
	followRepo "ualabackend/repositories/follow"

	"github.com/gin-gonic/gin"
)

func followRoutes(router *gin.Engine, repo *followRepo.Repository) {
	follows := router.Group("/follows")
	{
		follows.GET("/", func(c *gin.Context) { getAllFollows(c, repo) })
		follows.POST("/", func(c *gin.Context) { createFollow(c, repo) })
		follows.GET("/:follower_id", func(c *gin.Context) { getFollowedByFollowerID(c, repo) })
		follows.GET("/:follower_id/:followed_id", func(c *gin.Context) { getFollowByID(c, repo) })
		follows.DELETE("/:follower_id/:followed_id", func(c *gin.Context) { deleteFollow(c, repo) })

	}
}

// getAllFollows godoc
// @Summary Obtener todos los follows
// @Description Devuelve una lista de todos los follows
// @Tags follows
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /follows/ [get]
func getAllFollows(c *gin.Context, repo *followRepo.Repository) {
	follows, err := repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener follows"})
		return
	}
	c.JSON(http.StatusOK, follows)
}

// createFollow godoc
// @Summary Crear un nuevo follow
// @Description Crea un follow entre usuarios
// @Tags follows
// @Accept json
// @Produce json
// @Param follow body follow.FollowInput true "Datos del follow"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /follows/ [post]
func createFollow(c *gin.Context, repo *followRepo.Repository) {
	var payload follow.FollowInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err := repo.Create(payload.FollowerID, payload.FollowedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el follow"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Follow creado"})
}

// getFollowByID godoc
// @Summary Obtener un follow por ID
// @Description Devuelve un follow según el ID provisto
// @Tags follows
// @Produce json
// @Param follower_id path int true "ID del seguidor"
// @Param followed_id path int true "ID del seguido"
// @Success 200 {object} map[string]interface{}
// @Router /follows/{follower_id}/{followed_id} [get]
func getFollowByID(c *gin.Context, repo *followRepo.Repository) {
	followerID, err1 := strconv.Atoi(c.Param("follower_id"))
	followedID, err2 := strconv.Atoi(c.Param("followed_id"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs inválidos"})
		return
	}

	f, err := repo.GetByIDs(followerID, followedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar follow"})
		return
	}
	if f == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Follow no encontrado"})
		return
	}

	c.JSON(http.StatusOK, f)
}

// deleteFollow godoc
// @Summary Eliminar un follow
// @Description Elimina un follow por ID
// @Tags follows
// @Produce json
// @Param follower_id path int true "ID del seguidor"
// @Param followed_id path int true "ID del seguido"
// @Success 200 {object} map[string]interface{}
// @Router /follows/{follower_id}/{followed_id} [delete]
func deleteFollow(c *gin.Context, repo *followRepo.Repository) {
	followerID, err1 := strconv.Atoi(c.Param("follower_id"))
	followedID, err2 := strconv.Atoi(c.Param("followed_id"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs inválidos"})
		return
	}

	err := repo.Delete(followerID, followedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar follow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Follow eliminado"})
}

// getFollowedByFollowerID godoc
// @Summary Listar seguidos por un follower
// @Description Devuelve todos los usuarios seguidos por un follower
// @Tags follows
// @Produce json
// @Param follower_id path int true "ID del seguidor"
// @Success 200 {array} follow.Follow
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /follows/{follower_id} [get]
func getFollowedByFollowerID(c *gin.Context, repo *followRepo.Repository) {
	followerID, err := strconv.Atoi(c.Param("follower_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	follows, err := repo.GetFollowedByFollowerID(followerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener seguidos"})
		return
	}

	c.JSON(http.StatusOK, follows)
}

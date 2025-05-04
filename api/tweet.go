package api

import (
	"net/http"
	"strconv"
	"time"

	tweet "ualabackend/entities/tweet"
	tweetRepo "ualabackend/repositories/tweet"

	"github.com/gin-gonic/gin"
)

type UpdateTweetRequest struct {
	Message string `json:"message" binding:"required"`
}

func tweetRoutes(router *gin.Engine, repo *tweetRepo.Repository) {

	tweets := router.Group("/tweets")
	{
		tweets.GET("/", func(c *gin.Context) { getAllTweets(c, repo) })
		tweets.POST("/", func(c *gin.Context) { createTweet(c, repo) })
		tweets.GET("/:id", func(c *gin.Context) { getTweetByID(c, repo) })
		tweets.PUT("/:id", func(c *gin.Context) { updateTweet(c, repo) })
		tweets.DELETE("/:id", func(c *gin.Context) { deleteTweet(c, repo) })
	}
}

// getAllTweets godoc
// @Summary Obtener todos los tweets
// @Description Devuelve una lista de todos los tweets
// @Tags tweets
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /tweets/ [get]
func getAllTweets(c *gin.Context, repo *tweetRepo.Repository) {
	tweets, err := repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los tweets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tweets": tweets})
}

// @Summary Crear un nuevo tweet
// @Description Crea un nuevo tweet en el sistema
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet body tweet.TweetInput true "Datos del tweet"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tweets/ [post]
func createTweet(c *gin.Context, repo *tweetRepo.Repository) {
	var input tweet.TweetInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	t := tweet.Tweet{
		Timestamp: time.Now(),
		Message:   input.Message,
		Author_id: input.AuthorID,
	}

	if err := repo.Create(t.Author_id, t.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el tweet"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tweet creado"})
}

// getTweetByID godoc
// @Summary Obtener tweet por ID
// @Description Devuelve un tweet según el ID proporcionado
// @Tags tweets
// @Produce json
// @Param id path int true "ID del tweet"
// @Success 200 {object} map[string]interface{}
// @Router /tweets/{id} [get]
func getTweetByID(c *gin.Context, repo *tweetRepo.Repository) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	t, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}
	if t == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tweet no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tweet": t})
}

// updateTweet godoc
// @Summary Actualizar un tweet
// @Description Actualiza el contenido de un tweet
// @Tags tweets
// @Accept json
// @Produce json
// @Param id path int true "ID del tweet"
// @Param message body UpdateTweetRequest true "Nuevo mensaje en el cuerpo"
// @Success 200 {object} map[string]interface{}
// @Router /tweets/{id} [put]
func updateTweet(c *gin.Context, repo *tweetRepo.Repository) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mensaje requerido"})
		return
	}

	if err := repo.Update(id, input.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el tweet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tweet actualizado"})
}

// deleteTweet godoc
// @Summary Eliminar un tweet
// @Description Elimina un tweet por ID
// @Tags tweets
// @Produce json
// @Param id path int true "ID del tweet"
// @Success 200 {object} map[string]interface{}
// @Router /tweets/{id} [delete]
func deleteTweet(c *gin.Context, repo *tweetRepo.Repository) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el tweet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tweet eliminado"})
}

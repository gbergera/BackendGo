package api

import (
	"net/http"
	"strconv"

	user "ualabackend/entities/user"
	userRepo "ualabackend/repositories/user"

	"github.com/gin-gonic/gin"
)

func userRoutes(router *gin.Engine, repo *userRepo.Repository) {
	users := router.Group("/users")
	{
		users.GET("/", func(c *gin.Context) { getAllUsers(c, repo) })
		users.POST("/", func(c *gin.Context) { createUser(c, repo) })
		users.GET("/:id", func(c *gin.Context) { getUserByID(c, repo) })
		users.PUT("/:id", func(c *gin.Context) { updateUser(c, repo) })
		users.DELETE("/:id", func(c *gin.Context) { deleteUser(c, repo) })
	}
}

// getAllUsers godoc
// @Summary Obtener todos los usuarios
// @Description Devuelve una lista de todos los usuarios
// @Tags usuarios
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users/ [get]
func getAllUsers(c *gin.Context, repo *userRepo.Repository) {
	users, err := repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// createUser godoc
// @Summary Crear un nuevo usuario
// @Description Crea un usuario en el sistema
// @Tags usuarios
// @Accept json
// @Produce json
// @Param user body user.UserInput true "Datos del usuario"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/ [post]
func createUser(c *gin.Context, repo *userRepo.Repository) {
	var payload user.UserInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err := repo.Create(payload.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado"})
}

// getUserByID godoc
// @Summary Obtener un usuario por ID
// @Description Devuelve un usuario según el ID proporcionado
// @Tags usuarios
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [get]
func getUserByID(c *gin.Context, repo *userRepo.Repository) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	u, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}
	if u == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, u)
}

// updateUser godoc
// @Summary Actualizar un usuario
// @Description Actualiza los datos de un usuario existente
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param name query string true "Nuevo nombre del usuario"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [put]
func updateUser(c *gin.Context, repo *userRepo.Repository) {
	// Get the user ID from the URL path
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Get the new name from the query parameter
	name := c.DefaultQuery("name", "") // Default to empty if not provided
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'name' es requerido"})
		return
	}

	// Call the repository's Update method
	err = repo.Update(id, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el usuario"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado"})
}

// deleteUser godoc
// @Summary Eliminar un usuario
// @Description Elimina un usuario por ID
// @Tags usuarios
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [delete]
func deleteUser(c *gin.Context, repo *userRepo.Repository) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el usuario"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado"})
}

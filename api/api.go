package api

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ualabackend/db"
	followRepo "ualabackend/repositories/follow"
	tweetRepo "ualabackend/repositories/tweet"
	userRepo "ualabackend/repositories/user"
)

func InitAPI() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("❌ Could not initialize database:", err)
	}

	router := gin.Default()
	router.GET("/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	followRepository := followRepo.NewRepository(database)
	tweetRepository := tweetRepo.NewRepository(database)
	userRepository := userRepo.NewRepository(database)

	userRepository.Create("Usuario1")
	userRepository.Create("Usuario2")
	userRepository.Create("Usuario3")
	followRepository.Create(1, 2)
	followRepository.Create(2, 1)
	followRepository.Create(3, 1)
	tweetRepository.Create(1, "esto es una prueba")
	tweetRepository.Create(2, "esto tambien")
	tweetRepository.Create(3, "ola")
	userRoutes(router, userRepository)
	tweetRoutes(router, tweetRepository)
	followRoutes(router, followRepository)

	router.Run(":9090")

}

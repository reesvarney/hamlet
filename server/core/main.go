package main

import (
	"crypto/rand"
	"hamlet/server/core/auth"
	"hamlet/server/core/database"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
)

func GenerateSecret(secret_size int) []byte {
	buf := make([]byte, secret_size)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal("Failed to generate server secret")
	}
	return buf
}

func main() {
	r := gin.Default()
	secret := GenerateSecret(32)
	db := database.Connect()
	store, err := postgres.NewStore(db, secret)
	if err != nil {
		log.Fatal("Could not initialise postgres session store")
	}
	r.Use(sessions.Sessions("hamlet_session", store))
	auth.Routes(r, db)
	r.Run(":5376")
}

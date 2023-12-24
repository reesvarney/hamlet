package auth

import (
	"bytes"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChallengeResponseBody struct {
	Decoded string `json:"decoded" binding:"required"`
}

const RECAPTCHA_ENABLED = true

func Routes(route *gin.Engine, db *sql.DB) {
	auth := route.Group("/auth")

	// First verify ownership of the public key
	auth.GET("/challenge", func(ctx *gin.Context) {
		// Get the public key from the query params
		url_values := ctx.Request.URL.Query()
		if !url_values.Has("public_key") {
			ctx.JSON(http.StatusBadRequest, "")
			return
		}

		// Parse the public key
		pub, pubkey_bytes, err := ParsePubKey(url_values.Get("public_key"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		// Get the challenge
		challenge, encoded, err := GenerateChallenge(&pub)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": ""})
			return
		}

		// Store the challenge in the database
		_, err = db.Query(`INSERT INTO core_auth_challenges (public_key, challenge) VALUES ($1, $2)`, pubkey_bytes, challenge)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": ""})
			return
		}

		// Send the encoded challenge to the client
		ctx.JSON(http.StatusOK, gin.H{"challenge": encoded})
	})

	// Validate challenge response
	auth.POST("/challenge", func(ctx *gin.Context) {
		// Get the public key from the query params
		url_values := ctx.Request.URL.Query()
		if !url_values.Has("public_key") {
			ctx.JSON(http.StatusBadRequest, "")
			return
		}
		public_key := url_values.Get("public_key")

		// Read challenge response from request body
		var response ChallengeResponseBody
		err := ctx.BindJSON(response)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body_data"})
			return
		}
		decoded_bytes := []byte(response.Decoded)

		// Parse the public key
		_, pubkey_bytes, err := ParsePubKey(public_key)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		// Get original challenge from the database
		var challenge []byte
		err = db.QueryRow(`SELECT challenge FROM core_auth_challenges WHERE public_key = $1`, pubkey_bytes).Scan(&challenge)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "challenge_not_found"})
			return
		}

		// Check if original challenge is equal to decoded challenge
		if !bytes.Equal(decoded_bytes, challenge) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "challenge_not_match"})
			return
		}

		// Check if client is registring in the query params
		// Eventually this could use further parameters to detect potentially suspicious behaviour such as logging in from a new region
		if url_values.Has("register") {
			// Verify client is not automated
			ctx.JSON(http.StatusOK, gin.H{
				"verify_url": "/register_verify",
			})
			return
		}
		// Else, authenticate the session
		AuthenticateSession(ctx, db)
	})

	// Serve the verification page
	auth.GET("/verify", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	// Finally, register the user
	auth.POST("/verify", func(ctx *gin.Context) {
		// If recaptcha disabled or recaptcha API response valid
		// Insert user values into the database
		AuthenticateSession(ctx, db)
	})
}

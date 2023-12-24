package auth

import (
	"database/sql"
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthenticateSession(ctx *gin.Context, db *sql.DB) error {
	session := sessions.Default(ctx)
	_, pubkey_bytes, err := ParsePubKey(ctx.Query("public_key"))
	if err != nil {
		return err
	}
	var user_id string
	err = db.QueryRow("SELECT (id) FROM hamlet_users WHERE public_key = $1 LIMIT 1", pubkey_bytes).Scan(&user_id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		// Add new user to db
		_, err = db.Query("INSERT INTO hamlet_users () VALUES ()")
		if err != nil {
			return err
		}
	}
	if err = session.Save(); err != nil {
		return err
	}
	return nil
}

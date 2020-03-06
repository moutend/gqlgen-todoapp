package auth

import (
	"context"
	"log"
	"net/http"

	dbmodel "github.com/moutend/gqlgen-todoapp/internal/db/model"
	database "github.com/moutend/gqlgen-todoapp/internal/db/mysql"
	"github.com/moutend/gqlgen-todoapp/internal/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	userCtxKey = &contextKey{"user"}
)

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("TOKEN")

			if err != nil || c == nil {
				log.Println(err)

				next.ServeHTTP(w, r)
				return
			}

			token := c.Value
			username, err := jwt.ParseToken(token)

			if err != nil {
				log.Println(err)
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			tx, err := database.Db.Begin()

			if err != nil {
				log.Println(err)
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			defer func() {
				if err != nil {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}()

			user, err := dbmodel.Users(dbmodel.UserWhere.Name.EQ(username)).One(r.Context(), tx)

			if err != nil {
				log.Println(err)
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *dbmodel.User {
	user, _ := ctx.Value(userCtxKey).(*dbmodel.User)
	return user
}

func IsValidPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

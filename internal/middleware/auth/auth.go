package auth

import (
	"context"
	"log"
	"net/http"

	dbmodel "github.com/moutend/gqlgen-todoapp/internal/db/model"
	"github.com/moutend/gqlgen-todoapp/internal/jwt"
	"github.com/moutend/gqlgen-todoapp/internal/middleware/db"
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
			defer func() {
				if r := recover(); r != nil {
					log.Println("auth:", r)
					http.Error(w, `{"errors":[{"message":"access denied","path":[]}],"data":null}`, http.StatusForbidden)
				}
			}()

			c, err := r.Cookie("TOKEN")

			if err != nil || c == nil {
				log.Println("auth:", err)

				next.ServeHTTP(w, r)
				return
			}

			name, err := jwt.ParseToken(c.Value)

			if err != nil {
				log.Println("auth:", err)
				http.Error(w, `{"errors":[{"message":"invalid token","path":[]}],"data":null}`, http.StatusForbidden)
				return
			}

			tx := db.ForContext(r.Context())

			if tx == nil {
				log.Println("auth: internal error")
				http.Error(w, `{"errors":[{"message":"internal error","path":[]}],"data":null}`, http.StatusForbidden)
				return
			}

			user, err := dbmodel.Users(dbmodel.UserWhere.Name.EQ(name)).One(r.Context(), tx)

			if err != nil {
				log.Println("auth:", err)
				http.Error(w, `{"errors":[{"message":"internal error","path":[]}],"data":null}`, http.StatusForbidden)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), userCtxKey, user))

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

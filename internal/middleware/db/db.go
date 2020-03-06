package db

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	database "github.com/moutend/gqlgen-todoapp/internal/db/mysql"
)

type Transaction struct {
	*sql.Tx
	error error
}

func (tx *Transaction) Error(err error) {
	tx.error = err
}

var (
	transactionCtxKey = &contextKey{"transaction"}
)

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := database.Db.Begin()

			if err != nil {
				log.Println("db:", err)
				http.Error(w, `{"errors":[{"message":"internal error","path":[]}],"data":null}`, http.StatusInternalServerError)
				return
			}

			log.Println("db: Begin transaction")
			transaction := &Transaction{Tx: tx}
			r = r.WithContext(context.WithValue(r.Context(), transactionCtxKey, transaction))

			defer func() {
				if transaction.error != nil {
					log.Println("db: Rollback transaction")
					tx.Rollback()
				} else {
					log.Println("db: Commit transaction")
					tx.Commit()
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *Transaction {
	tx, _ := ctx.Value(transactionCtxKey).(*Transaction)
	return tx
}

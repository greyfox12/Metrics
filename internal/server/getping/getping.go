package getping

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func GetPing(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		//		fmt.Printf("CounterPage \n")

		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			fmt.Printf("Error Ping DB: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(nil)
	}
}

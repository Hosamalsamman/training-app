package dbutil

import (
	"errors"
	"net/http"

	"github.com/jackc/pgconn"
)

func ParseDBError(err error) (int, string) {

	if err == nil {
		return http.StatusOK, ""
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		case "23505":
			return http.StatusBadRequest, "Duplicate data."

		case "23503":
			return http.StatusBadRequest, "Referenced record does not exist."

		case "23502":
			return http.StatusBadRequest, "Required field is missing."

		case "23514":
			return http.StatusBadRequest, "Data violates a database rule."

		default:
			return http.StatusInternalServerError, pgErr.Message
		}
	}

	return http.StatusInternalServerError, err.Error()
}
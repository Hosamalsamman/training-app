package dbutil

import (
	"errors"
	"net/http"

	"github.com/jackc/pgconn"
)

func TranslateDBError(err error) (int, string) {

	if err == nil {
		return http.StatusOK, ""
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		case "23505":
			return http.StatusBadRequest,
				"البيانات المدخلة موجودة بالفعل."

		case "23503":
			return http.StatusBadRequest,
				"لا يمكن إتمام العملية لوجود بيانات مرتبطة غير صحيحة."

		case "23502":
			return http.StatusBadRequest,
				"يرجى إدخال جميع البيانات المطلوبة."

		case "22001":
			return http.StatusBadRequest,
				"إحدى القيم المدخلة أطول من الحد المسموح."

		case "23514":
			return http.StatusBadRequest,
				"البيانات المدخلة لا تحقق شروط النظام."

		default:
			return http.StatusInternalServerError,
				"حدث خطأ في قاعدة البيانات."
		}
	}

	return http.StatusInternalServerError,
		"حدث خطأ غير متوقع."
}
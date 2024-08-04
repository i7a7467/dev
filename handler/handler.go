package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/i7a7467/dev/db"
	"github.com/i7a7467/dev/model"
)

func HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	res, err := db.GetCurrentTime()
	if err != nil {
		http.Error(w, "error occured.", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, res+" health check is ok.")
	}
}

func StatusCheckHandler(w http.ResponseWriter, req *http.Request) {
	dbConn, err := db.DBConn()
	if err != nil {
		http.Error(w, "error occured.", http.StatusInternalServerError)
	}

	ctx := context.Background()

	var res string
	err = dbConn.NewSelect().ColumnExpr("current_timestamp").Scan(ctx, &res) //result, err := dbConn.NewSelect().Model(&nodejs).ColumnExpr("lower(url)").Where("? = ?", bun.Ident("id"), "42").Scan(ctx)

	if err != nil {
		http.Error(w, "error occured. get random data", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "status is ok. :"+res)

}

func GetOnePersonHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

    _ , err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbConn, err := db.DBConn()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person := new(model.Person)
	err = dbConn.NewSelect().Model(person).Where("id = ?", req.PathValue("id")).Scan(context.Background())

	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(person)
	} else if err == sql.ErrNoRows {
		// Check no rows in result set.
		// https://github.com/uptrace/bun/issues/876
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return
}

package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/i7a7467/dev/cache"
	"github.com/i7a7467/dev/db"
	"github.com/i7a7467/dev/model"
)

type Handler struct {
    Cache *cache.Cache
}

func NewHandler(cache *cache.Cache) *Handler {
    return &Handler{Cache: cache}
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	res, err := db.GetCurrentTime()
	if err != nil {
		http.Error(w, "error occured.", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, res+" health check is ok.")
	}
}

func (h *Handler) StatusCheckHandler(w http.ResponseWriter, req *http.Request) {
	dbConn, err := db.DBConn()
	if err != nil {
		http.Error(w, "error occured.", http.StatusInternalServerError)
	}

	ctx := context.Background()

	var res string
	err = dbConn.NewSelect().ColumnExpr("current_timestamp").Scan(ctx, &res)

	if err != nil {
		http.Error(w, "error occured. get random data", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "status is ok. :"+res)

}

func (h *Handler) GetOneAccountHandler(w http.ResponseWriter, req *http.Request) {

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

func (h *Handler) GetAccountsHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	cacheKey := req.RequestURI
	accountsCache ,err := h.Cache.GetCache(cacheKey)

	if err == bigcache.ErrEntryNotFound  {

		dbConn, err := db.DBConn()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		accounts := []model.Person{}
		err = dbConn.NewSelect().Model(&accounts).OrderExpr("id DESC").Limit(10).Scan(context.Background())
	    if err != nil {
			http.Error(w, "Failed to get data", http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&accounts)
			cacheImportData, err := json.Marshal(&accounts)
			if err != nil {
				fmt.Printf("Failed to parse import cache data . : %v\n", err)
				return
			}
			if err := h.Cache.SetCache(cacheKey, cacheImportData); err != nil {
        		fmt.Printf("Failed to set cache. : %v\n", err)
				return
			}
			return
		} 
	} else if err != nil {
		http.Error(w, "Failed to get cache", http.StatusInternalServerError)
	} else if accountsCache != nil {
		w.WriteHeader(http.StatusOK)
		w.Write(accountsCache)
		return
	}
}

func (h *Handler) CacheTestHandler(w http.ResponseWriter, req *http.Request) {

	now := time.Now()
	millis := now.UnixNano() / int64(time.Millisecond)
    humanReadableTime := time.Unix(0, millis*int64(time.Millisecond))

    value := humanReadableTime.Format("2006-01-02 15:04:05")

	key := "key"

	ck ,err := h.Cache.GetCache(key)

	if err == bigcache.ErrEntryNotFound  {
		if err := h.Cache.SetCache(key, []byte(value)); err != nil {
        	http.Error(w, "Failed to set cache", http.StatusInternalServerError)
        	return
		} 
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, value)
		fmt.Printf("cache.ErrEntryNotFound: %v\n", value)
		return
    } else if err != nil {
		fmt.Printf("cache.Err: %v\n", value)
		http.Error(w, "Failed to get cache", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(ck))
	fmt.Printf("cache.Found: %v\n", string(ck))
	return
}
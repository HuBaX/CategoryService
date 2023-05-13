package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"vsmlab/categoryservice/datahandling"

	"encoding/json"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

func main() {
	dbUser := os.Getenv("MYSQL_USER")
	dbAddr := os.Getenv("MYSQL_ADDRESS")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Net:                  "tcp",
		Addr:                 dbAddr,
		DBName:               dbName,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	queries := datahandling.New(db)

	http.HandleFunc("/addCategory", handleAddCategory(ctx, queries))
	http.HandleFunc("/getCategory", handleGetCategoryById(ctx, queries))
	http.HandleFunc("/getCategories", handleGetCategories(ctx, queries))
	http.HandleFunc("/getCategoryByName", handleGetCategoryByName(ctx, queries))
	http.HandleFunc("/delCategoryById", handleDelCategory(ctx, queries))

	http.ListenAndServe("0.0.0.0:8081", nil)
}

func handleAddCategory(ctx context.Context, queries *datahandling.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			setError(w, ErrMethodNotAllowed)
			return
		}
		var nameJSON map[string]any
		err := readJSON(r, &nameJSON)
		defer r.Body.Close()
		if err != nil {
			setError(w, ErrReadJSON)
			return
		}

		name := nameJSON["name"].(string)
		_, err = queries.AddCategory(ctx, name)
		if err != nil {
			setError(w, ErrQueryDatabase)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleGetCategoryById(ctx context.Context, queries *datahandling.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			setError(w, ErrMethodNotAllowed)
			return
		}

		idString := r.URL.Query().Get("id")
		if len(idString) == 0 {
			setError(w, ErrIDNotSet)
		}

		idInt, err := strconv.Atoi(idString)
		if err != nil {
			setError(w, ErrStrToInt)
			return
		}

		id := int32(idInt)

		if id < 0 {
			setError(w, ErrIDNegative)
		}

		category, err := queries.GetCategory(ctx, id)
		if err != nil {
			setError(w, ErrQueryDatabase)
			return
		}

		err = writeJSON(w, http.StatusOK, category)
		if err != nil {
			setError(w, ErrWriteJSON)
		}
	}
}

func handleGetCategories(ctx context.Context, queries *datahandling.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			setError(w, ErrMethodNotAllowed)
			return
		}

		categories, err := queries.GetCategories(ctx)
		if err != nil {
			setError(w, ErrQueryDatabase)
			return
		}

		categoryMap := map[string][]datahandling.Category{
			"categories": categories,
		}

		err = writeJSON(w, http.StatusOK, categoryMap)
		if err != nil {
			setError(w, ErrWriteJSON)
		}
	}
}

func handleGetCategoryByName(ctx context.Context, queries *datahandling.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			setError(w, ErrMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")

		if len(name) == 0 {
			setError(w, ErrNameNotSet)
		}

		categories, err := queries.GetCategoryByName(ctx, name)
		if err != nil {
			setError(w, ErrQueryDatabase)
			return
		}

		categoryMap := map[string][]datahandling.Category{
			"categories": categories,
		}

		err = writeJSON(w, http.StatusOK, categoryMap)
		if err != nil {
			setError(w, ErrWriteJSON)
		}
	}
}

func handleDelCategory(ctx context.Context, queries *datahandling.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			setError(w, ErrMethodNotAllowed)
			return
		}

		var idJSON map[string]any
		err := readJSON(r, &idJSON)
		defer r.Body.Close()
		if err != nil {
			setError(w, ErrReadJSON)
			return
		}

		id := idJSON["id"].(int32)
		err = queries.DelCategory(ctx, id)
		if err != nil {
			setError(w, ErrQueryDatabase)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}

func readJSON(r *http.Request, value *map[string]any) error {
	return json.NewDecoder(r.Body).Decode(value)
}

func setError(w http.ResponseWriter, err apiError) {
	w.WriteHeader(err.Status)
	w.Write([]byte(err.Msg))
}

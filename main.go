package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"vsmlab/categoryservice/datahandling"

	"encoding/json"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

type ApiResponse struct {
	Hostname string `json:"hostname"`
}

type CategoryResponse struct {
	Category    datahandling.Category `json:"category"`
	Apiresponse ApiResponse           `json:"apiresponse"`
}

type CategoriesResponse struct {
	Categories  []datahandling.Category `json:"categories"`
	Apiresponse ApiResponse             `json:"apiresponse"`
}

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
	http.HandleFunc("/delCategoryById", handleDelCategory(ctx, db, queries))

	http.ListenAndServe("0.0.0.0:8081", nil)
}

func handleAddCategory(ctx context.Context, queries *datahandling.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			setError(w, ErrMethodNotAllowed, "")
			return
		}
		var nameJSON map[string]any
		err := readJSON(r, &nameJSON)
		defer r.Body.Close()
		if err != nil {
			setError(w, ErrReadJSON, err.Error())
			return
		}

		name := nameJSON["name"].(string)
		err = queries.AddCategory(ctx, name)
		if err != nil {
			setError(w, ErrQueryDatabase, err.Error())
			return
		}

		hostname, err := os.Hostname()
		if err != nil {
			setError(w, ErrHostname, err.Error())
			return
		}
		apiResponse := ApiResponse{Hostname: hostname}

		err = writeJSONResponse(w, http.StatusOK, apiResponse)
		if err != nil {
			setError(w, ErrWriteJSON, err.Error())
		}
	}
}

func handleGetCategoryById(ctx context.Context, queries *datahandling.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			setError(w, ErrMethodNotAllowed, "")
			return
		}

		idString := r.URL.Query().Get("id")
		if len(idString) == 0 {
			setError(w, ErrIDNotSet, "")
			return
		}

		idInt, err := strconv.Atoi(idString)
		if err != nil {
			setError(w, ErrStrToInt, "")
			return
		}

		id := int32(idInt)

		if id < 0 {
			setError(w, ErrIDNegative, "")
			return
		}

		category, err := queries.GetCategory(ctx, id)
		if err != nil {
			setError(w, ErrQueryDatabase, err.Error())
			return
		}

		hostname, err := os.Hostname()
		if err != nil {
			setError(w, ErrHostname, err.Error())
			return
		}

		catResponse := CategoryResponse{
			Category:    category,
			Apiresponse: ApiResponse{Hostname: hostname},
		}

		err = writeJSONResponse(w, http.StatusOK, catResponse)
		if err != nil {
			setError(w, ErrWriteJSON, err.Error())
		}
	}
}

func handleGetCategories(ctx context.Context, queries *datahandling.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			setError(w, ErrMethodNotAllowed, "")
			return
		}

		categories := make([]datahandling.Category, 0)
		queriedCategories, err := queries.GetCategories(ctx)
		if err != nil {
			setError(w, ErrQueryDatabase, err.Error())
			return
		}

		categories = append(categories, queriedCategories...)

		hostname, err := os.Hostname()
		if err != nil {
			setError(w, ErrHostname, err.Error())
			return
		}

		catResponse := CategoriesResponse{
			Categories:  categories,
			Apiresponse: ApiResponse{Hostname: hostname},
		}

		err = writeJSONResponse(w, http.StatusOK, catResponse)
		if err != nil {
			setError(w, ErrWriteJSON, err.Error())
		}
	}
}

func handleGetCategoryByName(ctx context.Context, queries *datahandling.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			setError(w, ErrMethodNotAllowed, "")
			return
		}

		name := r.URL.Query().Get("name")

		if len(name) == 0 {
			setError(w, ErrNameNotSet, "")
		}

		categories := make([]datahandling.Category, 0)

		queriedCategories, err := queries.GetCategoryByName(ctx, name)
		if err != nil {
			setError(w, ErrQueryDatabase, err.Error())
			return
		}

		categories = append(categories, queriedCategories...)

		hostname, err := os.Hostname()
		if err != nil {
			setError(w, ErrHostname, err.Error())
		}

		catResponse := CategoriesResponse{
			Categories:  categories,
			Apiresponse: ApiResponse{Hostname: hostname},
		}

		err = writeJSONResponse(w, http.StatusOK, catResponse)
		if err != nil {
			setError(w, ErrWriteJSON, err.Error())
		}
	}
}

func handleDelCategory(ctx context.Context, db *sql.DB, queries *datahandling.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			setError(w, ErrMethodNotAllowed, "")
			return
		}

		idStr := r.URL.Query().Get("id")

		apiErr := delProductsByCategoryId(idStr)

		if apiErr != nil {
			setError(w, *apiErr, "")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			setError(w, ErrStrToInt, err.Error())
			return
		}

		err = queries.DelCategory(ctx, int32(id))
		if err != nil {
			setError(w, ErrQueryDatabase, err.Error())
			return
		}

		hostname, err := os.Hostname()
		if err != nil {
			setError(w, ErrHostname, err.Error())
		}

		apiResponse := ApiResponse{Hostname: hostname}
		err = writeJSONResponse(w, http.StatusOK, apiResponse)
		if err != nil {
			setError(w, ErrWriteJSON, err.Error())
		}
	}
}

func delProductsByCategoryId(categoryId string) *apiError {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://product-service:8082/delProductsByCategoryId?id="+categoryId, nil)
	if err != nil {
		return &ErrCreateRequest
	}

	resp, err := client.Do(req)
	if err != nil {
		return &ErrRequestProductDeletion
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &ErrRequestProductDeletion
	}
	return nil
}

func writeJSONResponse(w http.ResponseWriter, status int, value any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}

func readJSON(r *http.Request, value *map[string]any) error {
	return json.NewDecoder(r.Body).Decode(value)
}

func setError(w http.ResponseWriter, err apiError, returnedErrMsg string) {
	fmt.Println(err.Msg)
	fmt.Println(returnedErrMsg)
	w.WriteHeader(err.Status)
	w.Write([]byte(err.Msg))
}

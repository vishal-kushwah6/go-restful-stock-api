package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-postgress/model"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message"`
}

var db *sql.DB

func InitDB() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("POSTGRES_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	fmt.Println(" Database connected successfully")
}

func GetAstock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	stock, err := getstocks(int64(id))
	if err != nil {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllstocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stocks, err := getAllstocks()
	if err != nil {
		http.Error(w, "Error fetching stocks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stocks)
}

func Createstock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var stock model.Stock
	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if stock.Name == "" || stock.Company == "" || stock.Price <= 0 {
		http.Error(w, "Missing or invalid stock data", http.StatusBadRequest)
		return
	}

	insertID := insertStock(stock)

	res := Response{
		ID:      insertID,
		Message: "Stock created successfully",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func Deletestocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	deletedRows := deleteStock(int64(id))
	msg := fmt.Sprintf("Deleted successfully, total rows affected: %d", deletedRows)

	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func Updatestocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	var stock model.Stock
	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedRows := updateStock(int64(id), stock)
	msg := fmt.Sprintf("Updated successfully, total rows affected: %d", updatedRows)

	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func insertStock(stock model.Stock) int64 {
	sqlStatement := `INSERT INTO stocks (name, price, company) VALUES ($1, $2, $3) RETURNING stockid`
	var id int64
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Printf("Insert error: %v", err)
		return 0
	}
	return id
}

func getstocks(id int64) (model.Stock, error) {
	sqlStatement := `SELECT stockid, name, price, company FROM stocks WHERE stockid = $1`
	var stock model.Stock
	err := db.QueryRow(sqlStatement, id).Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	return stock, err
}

func getAllstocks() ([]model.Stock, error) {
	sqlStatement := `SELECT stockid, name, price, company FROM stocks`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []model.Stock
	for rows.Next() {
		var stock model.Stock
		if err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func deleteStock(id int64) int64 {
	sqlStatement := `DELETE FROM stocks WHERE stockid = $1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Printf("Delete error: %v", err)
		return 0
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("RowsAffected error: %v", err)
		return 0
	}
	return rowsAffected
}

func updateStock(id int64, stock model.Stock) int64 {
	sqlStatement := `UPDATE stocks SET name = $1, price = $2, company = $3 WHERE stockid = $4`
	res, err := db.Exec(sqlStatement, stock.Name, stock.Price, stock.Company, id)
	if err != nil {
		log.Printf("Update error: %v", err)
		return 0
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("RowsAffected error: %v", err)
		return 0
	}
	return rowsAffected
}

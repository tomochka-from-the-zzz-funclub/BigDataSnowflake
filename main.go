package main

import (
	config "bigdata/cfg"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

var (
	CSV_DIR  = getEnv("CSV_DIR", ".")
	WORK_DIR = getEnv("WORK_DIR", ".")
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func initDB(db *sql.DB) {
	fmt.Println("init")
	ddlPath := filepath.Join(WORK_DIR, "ddl.sql")
	ddlSQL, err := os.ReadFile(ddlPath)
	if err != nil {
		log.Fatalf("Ошибка чтения файла DDL: %v", err)
	}
	fmt.Println("init exec")
	_, err = db.Exec(string(ddlSQL))
	if err != nil {
		log.Fatalf("Ошибка выполнения DDL: %v", err)
	}
	fmt.Println("DDL выполнен")
}

func loadData(db *sql.DB) {
	fmt.Println("load")
	pattern := filepath.Join(CSV_DIR, "исходные данные", "MOCK_DATA*.csv")
	fmt.Println("Путь поиска:", pattern)

	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Printf("Ошибка поиска файлов: %v", err)
	}
	fmt.Printf("Найдено файлов: %d\n", len(files))
	fmt.Println("range")
	for _, path := range files {
		fmt.Printf("Загрузка %s\n", path)
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Ошибка открытия файла %s: %v", path, err)
			return
		}

		reader := csv.NewReader(file)
		// Пропускаем заголовок
		if _, err := reader.Read(); err != nil {
			fmt.Printf("Ошибка чтения заголовка из %s: %v", path, err)
		}
		fmt.Println("insert")
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("Ошибка чтения строки из %s: %v", path, err)
			}

			_, err = db.Exec(`INSERT INTO mock_data VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, 
			$34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50)`, toInterfaceSlice(record)...)
			if err != nil {
				fmt.Printf("Ошибка вставки данных из %s: %v", path, err)
			}
		}
		file.Close()
	}
	fmt.Printf("CSV загружены в mock_data")
}

func toInterfaceSlice(strs []string) []interface{} {
	res := make([]interface{}, len(strs))
	for i, v := range strs {
		res[i] = v
	}
	return res
}

func dml(db *sql.DB) {
	fmt.Println("dml")
	dmlSQL, err := os.ReadFile("dml.sql")
	if err != nil {
		log.Fatalf("Ошибка чтения файла DML: %v", err)
	}
	fmt.Println("dml exec")
	_, err = db.Exec(string(dmlSQL))
	if err != nil {
		log.Fatalf("Ошибка выполнения DML: %v", err)
	}
	fmt.Println("DML выполнен")
}

func main() {

	cfg := config.LoadConfig()
	//fmt.Printf("%s", cfg.DBHost, cfg.SslMode)
	//fmt.Println("create")
	base := NewPostgres(cfg)
	// fmt.Println("create")
	// fmt.Println("initDB")
	initDB(base.Connection)
	//fmt.Println("loadData")
	loadData(base.Connection)
	//fmt.Println("dml")
	dml(base.Connection)
}

type Postgres struct {
	Connection *sql.DB
}

func NewPostgres(cfg config.Config) *Postgres {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort, cfg.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return nil
	}
	fmt.Println("1min")
	//time.Sleep(time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	err = db.PingContext(ctx)
	cancel()
	if err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
		return nil
	} else {

		log.Printf("ping success")
	}
	fmt.Println("2min")
	//time.Sleep(time.Minute)

	return &Postgres{
		Connection: db,
	}
}

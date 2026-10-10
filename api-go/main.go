package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type GraphRequest struct {
	Edges [][]string `json:"edges"`
}

func main() {
	app := fiber.New()

	pythonURL := os.Getenv("PYTHON_WORKER_URL")
	if pythonURL == "" {
		pythonURL = "http://localhost:8000"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		
		dbURL = "postgres://admin:password@localhost:5432/algonode?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных: ", err)
	}
	defer db.Close()

	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			log.Println("Успешное подключение к PostgreSQL!")
			break
		}
		log.Println("Ожидание базы данных...")
		time.Sleep(2 * time.Second)
	}

	app.Post("/solve/cycle", func(c *fiber.Ctx) error {
		
		var req GraphRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Неверный формат JSON"})
		}

		payloadBytes, _ := json.Marshal(req)
		payloadStr := string(payloadBytes)
		taskType := "cycle_detection"

		var cachedResult string
		err := db.QueryRow(
			"SELECT result FROM graph_cache WHERE task_type = $1 AND graph_payload::jsonb = $2::jsonb LIMIT 1",
			taskType, payloadStr,
		).Scan(&cachedResult)

		if err == nil {
		
			var resultData map[string]interface{}
			json.Unmarshal([]byte(cachedResult), &resultData)
			return c.JSON(fiber.Map{
				"status": "success",
				"source": "database_cache",
				"data":   resultData,
			})
		}

		resp, err := http.Post(pythonURL+"/api/v1/cycle", "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil || resp.StatusCode != 200 {
			return c.Status(500).JSON(fiber.Map{"error": "Python-микросервис недоступен"})
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var pythonResult map[string]interface{}
		json.Unmarshal(body, &pythonResult)

		resultBytes, _ := json.Marshal(pythonResult)
		_, err = db.Exec(
			"INSERT INTO graph_cache (task_type, graph_payload, result) VALUES ($1, $2, $3)",
			taskType, payloadStr, string(resultBytes),
		)
		if err != nil {
			log.Println("Ошибка сохранения в кэш:", err)
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"source": "python_worker", 
			"data":   pythonResult,
		})
	})

	log.Println("AlgoNode API Gateway запущен на порту 8080")
	log.Fatal(app.Listen(":8080"))
}
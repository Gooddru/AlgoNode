# 🕸️ AlgoNode

A high-performance, microservice-based algorithmic engine designed to solve complex graph theoretical problems via a REST API. It demonstrates a modern approach to system design by combining the I/O speed of Go with the computational power of Python.

## Architecture
* **API Gateway (Go / Fiber):** Handles incoming REST requests, checks the PostgreSQL cache, and routes tasks.
* **Algorithm Worker (Python / FastAPI / NetworkX):** Performs heavy graph calculations (e.g., cycle detection, LaTeX/TikZ code generation).
* **Cache Database (PostgreSQL):** Stores previous calculations to return instant results for duplicate requests.

## How to Run (Docker)
1. Clone the repository.
2. Run `docker-compose up --build`.
3. The API will be available at `http://localhost:8085`.

## 📡 Example Usage
**POST** `/solve/cycle`
```json
{
    "edges": [
        ["A", "B"],
        ["B", "C"],
        ["C", "A"]
    ]
}
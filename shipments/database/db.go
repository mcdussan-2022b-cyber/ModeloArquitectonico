package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	// Variables de entorno con valores por defecto
	host := getEnv("PGHOST", "localhost")
	user := getEnv("PGUSER", "postgres")
	pass := getEnv("PGPASSWORD", "Liam1309")
	name := getEnv("PGDATABASE", "fasttrack_shipments")
	port := getEnv("PGPORT", "5432")

	// Cadena de conexión válida para PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Error al abrir conexión con PostgreSQL:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("⚠️ No se pudo verificar conexión con PostgreSQL:", err)
	}

	DB = db
	fmt.Println("✅ Conectado a PostgreSQL exitosamente")
}

// getEnv obtiene una variable de entorno o usa un valor por defecto
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

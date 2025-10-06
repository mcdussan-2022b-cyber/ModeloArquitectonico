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
	host := getEnv("PGHOST", "localhost")
	user := getEnv("PGUSER", "postgres")
	pass := getEnv("PGPASSWORD", "Liam1309") // ← tu password
	name := getEnv("PGDATABASE", "fasttrack_notifications")
	port := getEnv("PGPORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Error al abrir conexión con PostgreSQL:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("⚠️ No se pudo verificar conexión con PostgreSQL:", err)
	}

	DB = db
	fmt.Println("✅ Conectado a PostgreSQL exitosamente")

	if err := migrate(); err != nil {
		log.Fatal("❌ Error en migración:", err)
	}
}

func migrate() error {
	ddl := `
	CREATE TABLE IF NOT EXISTS notifications (
	  id UUID PRIMARY KEY,
	  recipient TEXT NOT NULL,
	  channel   TEXT NOT NULL,
	  message   TEXT NOT NULL,
	  status    TEXT NOT NULL DEFAULT 'PENDING',
	  created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`
	_, err := DB.Exec(ddl)
	return err
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

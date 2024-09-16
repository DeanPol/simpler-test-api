package db

import (
	"database/sql"
	"fmt" // i'm guessing to parse the variables to the string
	"log" // printing to console purposes
	"os"  // should be used for loading our .env files

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var DB *sql.DB

/*
	Grab .env variables regarding our db connection credentials.
	Initialize db and check connection.
	Finally, run migrations.

	Regarding fallbacks: Notice there are no return statements in the case of an error.
	With log.Fatalf(), we are already handling the error, as this function prints the error message, including any details,
	and then terminates the program. It is the equivalent to log.Printf() + os.Exit(1).
*/
func InitializeDB() *sql.DB {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatalf("Error opening database connection: %v", err)
    }

    // Check the connection
    err = db.Ping()
    if err != nil {
        log.Fatalf("Error pinging the database: %v", err)
    }

    log.Println("Successfully connected to the database!")

    runMigrations(db)

    DB = db
    return db
}

// Applies the latest migrations using golang-migrate
func runMigrations(db *sql.DB) {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatalf("Could not create migration driver: %v", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "postgres", driver)
    if err != nil {
        log.Fatalf("Could not start migration: %v", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Could not apply migrations: %v", err)
    }

    log.Println("Migrations applied successfully!")
}

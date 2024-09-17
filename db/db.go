package db

import (
	"database/sql"
	"fmt" // i'm guessing to parse the variables to the string
	"log" // printing to console purposes
	"os"  // should be used for loading our .env files

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
    dbHost := os.Getenv("POSTGRES_HOST")
    if dbHost == "" {
		dbHost = "postgres" // Fallback to the Docker service name
	}

	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort == "" {
		dbPort = "5432" // Default PostgreSQL port
	}

    dbUser := os.Getenv("POSTGRES_USER")
    dbPassword := os.Getenv("POSTGRES_PASSWORD")
    dbName := os.Getenv("POSTGRES_DB")

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

/*
    Running migrations manually here.
    Would rather have a migrations directory and have a more modular function body:

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "postgres", driver)

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Could not apply migrations: %v", err)
    }

    log.Println("Migrations applied successfully!")
}
    But this is causing me a lot of trouble when running docker for some reason.
*/

func runMigrations(db *sql.DB) {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            description TEXT,
            price DECIMAL(10, 2) NOT NULL,
            stock INTEGER NOT NULL
        );
    `)
    if err != nil {
        log.Fatalf("Error running migrations: %v", err)
    }
}


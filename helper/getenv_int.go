package helper

import (
	"os"
	"strconv"
)

/*
	Helper function to fetch .env variable with fallback value

	About exporting functions to use in another file.
	We have to define the package. Here our directory is 'helper' (hence our package name).
	In order to load our functions define within this 'helper' package in other files, we have to
	import that package. (i.e. import "simpler-test-api/helper")

	In Go, only functions, types and variables that START WITH AN UPPERCASE LETTER are exported (accessible from other packages).

*/
func GetEnvInt(key string, defaultValue int) int {
    valueStr := os.Getenv(key)
    if valueStr == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}

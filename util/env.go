package util

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func checkVars() []string {
	vars := []string{"GO_REST_ENV", "ADDR"}
	missing := []string{}
	for _, v := range vars {
		_, set := os.LookupEnv(v)
		if !set {
			missing = append(missing, v)
		}
	}
	return missing
}

// LoadEnv : Se cargan variables de entorno
func LoadEnv() {
	env := os.Getenv("TEMPLATE_ENV")
	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
	if vars := checkVars(); len(vars) != 0 {
		log.Printf("ERROR: Variables de entorno necesarias no definidas: %v", vars)
		panic(fmt.Sprintf("ERROR: Variables de entorno necesarias no definidas: %v", vars))
	}
}

package readEnv

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)




func ReadVariable(key string) string {

  // load .env file
  return os.Getenv(key)
}

func Init() {

	isTesting := os.Getenv("testing")

	if isTesting == "true" {
		return
	}

	fmt.Printf("Loading .env file testing: %s\n",isTesting )

	err := godotenv.Load("api.env")

  if err != nil {
		fmt.Printf(" %v\n", err)
    log.Fatalf("Error loading .env file",)
  }

}
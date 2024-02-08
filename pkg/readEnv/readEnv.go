package readEnv

import (
	"os"
)




func ReadVariable(key string) string {

  // load .env file
  return os.Getenv(key)
}

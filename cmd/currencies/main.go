package main

import (
	"net/http"
	"os"

	currencies "github.com/spicyt/currencies/internal/app"
)




func main() {
	r := currencies.InitRouter()

	println("Server running on port" + os.Getenv("PORT"))
	http.ListenAndServe(":" + os.Getenv("PORT"), r)
}

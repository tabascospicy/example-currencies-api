package main

import (
	"net/http"
	"os"

	app "github.com/spicyt/currencies/internal/currencies"
)




func main() {
	r := app.InitRouter()

	println("Server running on port " + os.Getenv("PORT"))
	http.ListenAndServe(":" + os.Getenv("PORT"), r)
}

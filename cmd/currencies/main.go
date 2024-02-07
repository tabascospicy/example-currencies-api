package main

import (
	"net/http"

	currencies "github.com/spicyt/currencies/internal/app"
	readEnv "github.com/spicyt/currencies/pkg"
)




func main() {
	readEnv.Init()
	r := currencies.InitRouter()
	println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}

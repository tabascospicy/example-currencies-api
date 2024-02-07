package currencies

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)




func InitRouter() *chi.Mux{


	r := chi.NewRouter()

	r.Use(middleware.Logger)
	
  // listen the root path and call the GetCurrencyExchangeHandler
	r.Get("/" , GetCurrencyExchangeHandler)
	
	return r
}
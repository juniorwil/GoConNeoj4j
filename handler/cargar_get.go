package handler

import (
	"encoding/json"
	"net/http"

	"github.com/juniorwil/chi/tienda/cargar"
)

func CargarGet(feed cargar.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := feed.GetAll()
		json.NewEncoder(w).Encode(items)
	}
}

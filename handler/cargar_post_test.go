package handler

import (
	"net/http"
	"testing"

	"github.com/juniorwil/chi/tienda/cargar"
	"github.com/juniorwil/chi/tienda/mock_http"
)

func TestCargarPost(t *testing.T) {

	feed := cargar.New()

	headers := http.Header{}
	headers.Add("content-type", "application/json")

	w := &mock_http.ResponseWriter{}
	r := &http.Request{
		Header: headers,
	}

	r.Body = mock_http.RequestBody(map[string]string{
		"Title": "Prueba",
		"Post":  "Tres",
	})

	handler := CargarPost()
	handler(w, r)

	result := w.GetBodyString()

	if result != "Guardado!" {
		t.Errorf("Handler did not complete")
	}

	if len(feed.GetAll()) != 1 {
		t.Errorf("Item did not add")
	}

	if feed.GetAll()[0].Title != "Prueba" {
		t.Errorf("Item bad")
	}

}

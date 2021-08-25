package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/juniorwil/chi/handler"
)

// Para pruebas de tokens
var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	// Para depuracion y pruebas creo un
	// jwt token con `user_id:123` aqui:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func main() {

	router := chi.NewRouter()
	// Cors para el tema de permisos desde la peticion del front
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Url especifica
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// Rutas protegidas
	router.Group(func(r chi.Router) {
		// Verificar y validar JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("Area protegida. hi %v", claims)))
		})
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Prueba de inicio golang"))
	})
	// Post a neoj4 para crear bd
	router.Post("/cargar", handler.CargarPost())
	// Listar compradores
	router.Post("/compradores", handler.ListarCompradores())
	// Listar compradores e historial
	router.Post("/historial", handler.ListarHistoria())
	// Listar otros compradores por la misma Ip
	router.Post("/ocompradores", handler.ListarOtros())
	// Listar de productos sueridos comprador por otros clientes
	router.Post("/sugeridos", handler.ListarSugeridos())
	// Llamado a server golang
	http.ListenAndServe(":3333", router)
}

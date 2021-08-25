package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/juniorwil/chi/internal"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func ListarCompradores() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		driver := internal.ConexNeo()
		sess := driver.NewSession(neo4j.SessionConfig{})
		defer sess.Close()

		batch := 1
		result, err := sess.Run("MATCH (n:compradores) RETURN n LIMIT 25",
			map[string]interface{}{"batch": batch})
		if err != nil {
			fmt.Println(err)
		}
		var rec *neo4j.Record
		var cadena string
		cadena += `[`
		for result.NextRecord(&rec) {
			node := rec.Values[0].(neo4j.Node)
			// Extraccion de los datos en los nodos
			ageInt := node.Props["age"].(int64)
			age := strconv.FormatInt(ageInt, 10)
			id := string(node.Props["id"].(string))
			name := string(node.Props["name"].(string))
			cadena += `{"id": "` + id + `", "name": "` + name + `", "age": "` + age + `" },`

		}
		cadena += `]`
		// Busco la ultima coma y la retiro porque genera error en la construccion del json
		cadFinal := strings.Replace(cadena, ",]", "]", -1)
		cadConver := []byte(cadFinal)
		var arrayJson interface{}
		json.Unmarshal(cadConver, &arrayJson)
		json.NewEncoder(w).Encode(arrayJson)
	}
}

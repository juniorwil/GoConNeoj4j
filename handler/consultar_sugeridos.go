package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/juniorwil/chi/internal"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func ListarSugeridos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraccion de campos
		comprador := r.PostFormValue("comprador")
		// Conexion
		driver := internal.ConexNeo()
		sess := driver.NewSession(neo4j.SessionConfig{})
		defer sess.Close()
		// Extra los nodos de productos que no han sido comprador por el cliente
		result, err := sess.Run(`match (w:Articulos) 
            where not exists{
                  match (a:compradores{ id : $idCom } )-[:COMPRA]->(b:Transacciones)<-[:ARTICULOS]-(c:Articulos) 
                   where w.id = c.id   
            }
			return w , rand() as r limit 10`, map[string]interface{}{"idCom": comprador})
		if err != nil {
			fmt.Println(err)
		}
		var rec *neo4j.Record
		var cadena string
		cadena += `[`
		for result.NextRecord(&rec) {
			// Nodo 0 correpondiente a los articulos sugeridos
			node := rec.Values[0].(neo4j.Node)
			// Extraccion de los datos en los nodos
			id := string(node.Props["id"].(string))
			name := string(node.Props["name"].(string))
			price := node.Props["price"].(string)
			// Extraccion de los datos en los nodos
			cadena += `{"id": "` + id + `", "name": "` + name + `", "price": "` + price + `" },`
		}
		cadena += `]`
		// Busco la ultima coma y la retiro porque genera error en la construccion del json
		cadFinal := strings.Replace(cadena, ",]", "]", -1)
		cadConver := []byte(cadFinal)
		var arrayJson interface{}
		json.Unmarshal(cadConver, &arrayJson)
		json.NewEncoder(w).Encode(arrayJson)

		// conjunto de registros si es un solo se usa Single
		//record, err := result.Collect()
		//if err != nil {
		//	fmt.Println(err)
		//}
		// Envio json
		//json.NewEncoder(w).Encode(record)
	}
}

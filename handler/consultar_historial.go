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

func ListarHistoria() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraccion de campos
		comprador := r.PostFormValue("comprador")
		// Conexion
		driver := internal.ConexNeo()
		sess := driver.NewSession(neo4j.SessionConfig{})
		defer sess.Close()

		result, err := sess.Run(`match (a:compradores{ id : $idCom })-[:COMPRA]->(b:Transacciones)<-[:ARTICULOS]-(c:Articulos)
		                               return a, b, c`,
			map[string]interface{}{"idCom": comprador})
		if err != nil {
			fmt.Println(err)
		}
		var rec *neo4j.Record
		var cadena string
		cadena += `[`
		for result.NextRecord(&rec) {
			// Nodo 0 correpondiente a los datos del comprador
			node := rec.Values[0].(neo4j.Node)
			// Extraccion de los datos en los nodos
			ageInt := node.Props["age"].(int64)
			age := strconv.FormatInt(ageInt, 10)
			id := string(node.Props["id"].(string))
			name := string(node.Props["name"].(string))

			// Nodo 1 correpondiente a los datos de la compra
			node1 := rec.Values[1].(neo4j.Node)
			idCom := string(node1.Props["id"].(string))
			device := string(node1.Props["device"].(string))
			ip := string(node1.Props["ip"].(string))

			// Nodo 2 correpondiente a los datos del producto
			node2 := rec.Values[2].(neo4j.Node)
			nameProd := string(node2.Props["name"].(string))
			idProd := string(node2.Props["id"].(string))
			price := node2.Props["price"].(string)

			// Extraccion de los datos en los nodos
			cadena += `{"id": "` + id + `", "name": "` + name + `", "age": "` + age + `", "idCom": "` + idCom + `", "idProd": "` + idProd + `", 
			            "nameProd": "` + nameProd + `", "price": "` + price + `", "device": "` + device + `", "ip": "` + ip + `" },`

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

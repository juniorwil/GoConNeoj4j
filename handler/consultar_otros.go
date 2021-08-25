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

func ListarOtros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraccion de campos
		comprador := r.PostFormValue("comprador")
		// Conexion
		driver := internal.ConexNeo()
		sess := driver.NewSession(neo4j.SessionConfig{})
		defer sess.Close()

		result, err := sess.Run(`match (a:compradores{id:"113c7c8c"})-[:COMPRA]->(b:Transacciones)<-[:IP]-(c:Transacciones)<-[:COMPRA]-(d:compradores)  
		                               return a, b, c, d`,
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

			// Nodo 3 correpondiente a los datos del otro comprador que hizo compras con la misma Ip
			node3 := rec.Values[3].(neo4j.Node)
			ageIntO := node3.Props["age"].(int64)
			ageO := strconv.FormatInt(ageIntO, 10)
			idO := string(node3.Props["id"].(string))
			nameO := string(node3.Props["name"].(string))

			// Extraccion de los datos en los nodos
			cadena += `{"id": "` + id + `", "name": "` + name + `", "age": "` + age + `", "idCom": "` + idCom + `", "idOtro": "` + idO + `", 
			            "nameOtro": "` + nameO + `", "ageOtro": "` + ageO + `", "device": "` + device + `", "ip": "` + ip + `" },`

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

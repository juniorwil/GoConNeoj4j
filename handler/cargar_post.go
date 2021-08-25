package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/juniorwil/chi/internal"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Cargar datos a neoj4
func CargarPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := map[string]string{}
		json.NewDecoder(r.Body).Decode(&request)
		// Conexion drivers neoj4
		driver := internal.ConexNeo()
		// funcion insertItem insericcion de datos neo4j
		item, err := insertItem(driver)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v\n", item)
		cadConver := []byte(`[{
			"Result":  "Ok"
		}]`)
		var arrayJson interface{}
		json.Unmarshal(cadConver, &arrayJson)
		json.NewEncoder(w).Encode(arrayJson)
	}
}

// Creacion de nodos neoj4
func createItemFn(tx neo4j.Transaction) (interface{}, error) {
	// 1. Lectura por Json via http
	resp, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers")
	if err != nil {
		fmt.Println("No es posible conectarse a la url")
	}
	defer resp.Body.Close()
	// response body is []byte
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	arregloComoJson := []byte(body)

	var nombre []campos
	// Cargar arreglo a nombredesde el responsive
	err2 := json.Unmarshal(arregloComoJson, &nombre)
	if err2 != nil {
		fmt.Println("error:", err2)
		os.Exit(1)
	}
	// recorrido de datos para enviar a neo4f.js compradores o clientes
	for _, values := range nombre {
		records, err := tx.Run("CREATE (a:compradores { id: $id, name: $name, age: $age, dt:datetime() }) RETURN a.id, a.name", map[string]interface{}{
			"id":   values.Id,
			"name": values.Name,
			"age":  values.Age,
		})
		_ = records
		if err != nil {
			return nil, err
		}
	}
	// 2. Lectura por Json
	//Bajar dese la Url el archivo csv
	Download("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products", "productos.csv")
	//Apertura archivo csv
	csvFile, err := os.Open("productos.csv")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	r.Comma = '\''
	r.Comment = '#'
	csvLines, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	// recorrido de datos para enviar a neo4f.js articulos o productos
	for _, line := range csvLines {
		records, err := tx.Run("CREATE (a:Articulos { id: $id, name: $name, price: $price, dt:datetime() }) RETURN a.id, a.name, a.price", map[string]interface{}{
			"id":    line[0],
			"name":  line[1],
			"price": line[2],
		})
		_ = records
		if err != nil {
			return nil, err
		}
	}

	// 3.Formato: No standard
	//Download("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions", "transactions")
	bytesLeidos, err := ioutil.ReadFile("transactions")
	if err != nil {
		fmt.Printf("Error leyendo archivo: %v", err)
	}
	// COnversion a cadena de texto
	cadena := string(bytesLeidos)
	lon := len(cadena)
	// Contamos totales de cierre de registros para nodos con el ) en cada parte del la cadena
	// para establecer el registro de cada nodo de ransaccion
	// Ejemplo:#000061009e00 e7114a0a 222.123.167.115 android (d3ad22dc,7812a47,fc5de8c5,a1b88b02,79571680)
	//         #000061009e01 7a71b889 215.33.130.224 linux (7379fda8,72c3e407,8b89140b)  , donde ) seria componente de cierre nodo
	total := strings.Count(cadena, ")")
	var conteo int = 1
	// recorro el total de ) para establecer fin de ciclo de recorrido
	for conteo <= total {
		//fmt.Println(greeting[i])Ejemplo para recorrer cada letra del archivo

		// Encuentra posicion del ) para ubicar el que sera el final de la entidad o nodo
		posicion := strings.Index(strings.TrimSpace(cadena), ")")
		// Extraccion de la cadena para armar el nodo
		data := cadena[0 : posicion+1]
		// Busco primer espacio para extraer id de la transaccion
		posId := strings.Index(strings.TrimSpace(cadena), " ")
		// Extraer el id de la transaccion para dalr id al nodo
		idNodo := cadena[1 : posId+1]
		// reemplazos para convertir en funcion de nodos para chyper
		// se reemplaza listado de productos en la transaccion
		der := strings.Replace(data, "(", ",product_ids:'", -1)
		// se reemplaza espacio entre id_tras e id_comprador
		der1 := strings.Replace(der, " ", "',buyer_id:'", 1)
		// se reemplaza espacio buyer_id e Ip
		der2 := strings.Replace(der1, " ", "',ip:'", 1)
		// se reemplaza espacio Ip a So
		der3 := strings.Replace(der2, " ", "',device:'", 1)
		// se reemplaza espacio Device y productos
		der4 := strings.Replace(der3, " ", "'", 1)
		// se arma primera sentencia create de chyper
		der5 := strings.Replace(der4, "#", "CREATE(T"+idNodo+":Transacciones{id:'", -1)
		der6 := strings.Replace(der5, ")", "' , dt:datetime()  })", -1)
		// recorrido de datos para enviar a neo4f.js transaccions de compras
		records, err := tx.Run(der6, map[string]interface{}{
			"id": 0,
		})
		_ = records
		if err != nil {
			return nil, err
		}
		// Cadena que va cortando la cadena siguiente para simpflicar
		cadena = strings.TrimSpace(cadena[posicion+1 : lon])
		lon = len(cadena)
		conteo = conteo + 1 // Incrementar para cerrar ciclo
	}
	fmt.Print("Ciclo terminado")

	// 4. crear relaciones
	// Relacion compradores/clientes - transacciones
	records, err := tx.Run("MATCH (u:compradores) match (f:Transacciones ) where u.id = f.buyer_id create (u)-[:COMPRA{}]->(f); ", map[string]interface{}{
		"id": 0,
	})
	_ = records
	if err != nil {
		return nil, err
	}
	// Relacion transacciones - productos o articulos
	recordsT, err := tx.Run("match (b:Articulos) match (a:Transacciones ) where a.product_ids CONTAINS b.id create (b)-[:ARTICULOS{}]->(a); ", map[string]interface{}{
		"id": 0,
	})
	_ = recordsT
	if err != nil {
		return nil, err
	}
	// Relacion transaccaciones - Ip
	recordsIp, err := tx.Run("match(a:Transacciones) match(b:Transacciones) where b.ip=a.ip create (a)-[:IP]->(b) ", map[string]interface{}{
		"id": 0,
	})
	_ = recordsIp
	if err != nil {
		return nil, err
	}
	// 5. Respuesta final
	return &Item{
		Name: "Nodos creados",
	}, nil
}

// Session a neo4
func insertItem(driver neo4j.Driver) (*Item, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	result, err := session.WriteTransaction(createItemFn)
	if err != nil {
		return nil, err
	}
	return result.(*Item), nil
}

func Download(url string, filename string) error {
	out, _ := os.Create(filename)
	defer out.Close()

	resp, err := http.Get(url)
	//defer resp.Body.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		log.Fatal(err)
	}

	return err
}

type Item struct {
	Id   int64
	Name string
}

type campos struct {
	Id   string `json:"Id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

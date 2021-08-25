package internal

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func ConexNeo() neo4j.Driver {
	dbUri := "neo4j://localhost:7687"
	driver, err := neo4j.NewDriver(dbUri, neo4j.BasicAuth("neo4j", "123", ""))
	if err != nil {
		panic(err)
	}
	//	defer driver.Close()
	return driver
}

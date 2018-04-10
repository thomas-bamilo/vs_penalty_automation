package scitemid

import (
	"database/sql"
	"log"

	dbconf "github.com/thomas-bamilo/vs_penalty_automation/dbconf"
	// driver for MySQL
	_ "github.com/go-sql-driver/mysql"
)

// ScItemID is a struct representing the table of relation between oms_item_number and sc_item_number
type ScItemID struct {
	OmsItemNumber int `json:"oms_item_number"`
	ScItemNumber  int `json:"sc_item_number"`
}

// CreateScItemID queries seller center and write result to ScItemID struct
func CreateScItemID(omsItemNumberFilter string) []ScItemID {

	// fetch database configuration
	var dbConf dbconf.DbConf
	dbConf.ReadYamlDbConf()
	// create connection string
	connStr := dbConf.ScUser + ":" + dbConf.ScPw + "@tcp(" + dbConf.ScHost + ")/" + dbConf.ScDb

	// connect to database
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// test connection with ping
	err = db.Ping()
	if err != nil {
		log.Println("Connection failed")
		log.Fatal(err)
	} else {
		log.Println("Connection successful!")
	}

	// store the query in a string
	query := `SELECT 
	COALESCE(soi.src_id,0) oms_item_number
   ,COALESCE(soi.id_sales_order_item,0) sc_item_number
   FROM sales_order_item soi WHERE COALESCE(soi.src_id,0) IN (` + omsItemNumberFilter + `)`

	var scItemNumber, omsItemNumber int
	var scItemID []ScItemID

	rows, _ := db.Query(query)

	for rows.Next() {
		err := rows.Scan(&omsItemNumber, &scItemNumber)
		if err != nil {
			log.Fatal(err)
		}
		scItemID = append(scItemID,
			ScItemID{
				ScItemNumber:  scItemNumber,
				OmsItemNumber: omsItemNumber,
			})
	}

	return scItemID
}

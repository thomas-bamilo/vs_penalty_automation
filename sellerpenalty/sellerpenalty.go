package sellerpenalty

import (
	"database/sql"
	"log"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/thomas-bamilo/sql/dbconf"
)

// SellerPenalty is a struct representing the table of penalty amounts per sales order item
type SellerPenalty struct {
	SupplierName  string `json:"supplier_name"`
	OrderNr       int    `json:"order_nr"`
	BobItemNumber int    `json:"bob_item_number"`
	OmsItemNumber int    `json:"oms_item_number"`
	ReturnReason  string `json:"return_reason"`
	CancelReason  string `json:"cancel_reason"`
	YearMonth     int    `json:"year_month"`
	Amount        int    `json:"amount"`
}

// CreateSellerPenalty queries oms and write result to SellerPenalty struct
func CreateSellerPenalty() []SellerPenalty {

	// fetch database configuration
	var dbConf dbconf.DbConf
	dbConf.ReadYamlDbConf()
	// create connection string
	connStr := dbConf.OmsUser + ":" + dbConf.OmsPw + "@tcp(" + dbConf.OmsHost + ")/" + dbConf.OmsDb

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
	COALESCE(is1.name_en,'NA') supplier_name
	,COALESCE(iso.order_nr,0)
	,COALESCE(isoi.bob_id_sales_order_item,0) bob_item_number
	  ,COALESCE(isoi.id_sales_order_item,0) oms_item_number
	,COALESCE(or2.name,'NA') return_reason
	  ,COALESCE(or1.name,'NA') cancel_reason
	,CONCAT(CASE WHEN MONTH(CURRENT_DATE()) = 1 THEN YEAR(CURRENT_DATE())-1 ELSE YEAR(CURRENT_DATE()) END
			  ,CASE WHEN MONTH(CURRENT_DATE()) = 1 THEN 12 ELSE MONTH(CURRENT_DATE())-1 END) 'year_month'
	,250000 amount
	  
	  FROM ims_sales_order_item isoi
	  
	  LEFT JOIN ims_sales_order iso
	  ON isoi.fk_sales_order = iso.id_sales_order
	  
	  LEFT JOIN ims_supplier is1
	  ON is1.bob_id_supplier = isoi.bob_id_supplier
	  
	  LEFT JOIN oms_reason or1
	  ON or1.id_reason = isoi.fk_cancel_reason
	  
	  LEFT JOIN oms_return_ticket ort
	  ON ort.fk_sales_order_item = isoi.id_sales_order_item
	  
	  LEFT JOIN oms_reason or2
	  ON or2.id_reason = ort.fk_return_reason
	  
	  WHERE MONTH(isoi.created_at) = CASE WHEN MONTH(CURRENT_DATE()) = 1 THEN 12 ELSE MONTH(CURRENT_DATE())-1 END
	  AND YEAR(isoi.created_at) = CASE WHEN MONTH(CURRENT_DATE()) = 1 THEN YEAR(CURRENT_DATE())-1 ELSE YEAR(CURRENT_DATE()) END
	  AND isoi.fk_sales_order_item_status <> 10
	  AND (LOWER(or2.name) LIKE '%merchant%' AND LOWER(or2.name) NOT LIKE '%item%'
		   OR (isoi.fk_cancel_reason = 642)
		   OR (isoi.fk_cancel_reason = 652)
		  )
	  AND LOWER(is1.name_en) NOT LIKE '%bamilo%'
	  AND is1.name_en <> 'RSH Co.'
	  
	  GROUP BY isoi.id_sales_order_item`

	var supplierName, returnReason, cancelReason string
	var orderNr, bobItemNumber, omsItemNumber, yearMonth, amount int
	var sellerPenalty []SellerPenalty

	rows, _ := db.Query(query)

	for rows.Next() {
		err := rows.Scan(&supplierName, &orderNr, &bobItemNumber, &omsItemNumber, &returnReason, &cancelReason, &yearMonth, &amount)
		if err != nil {
			log.Fatal(err)
		}
		sellerPenalty = append(sellerPenalty,
			SellerPenalty{SupplierName: supplierName,
				OrderNr:       orderNr,
				BobItemNumber: bobItemNumber,
				OmsItemNumber: omsItemNumber,
				ReturnReason:  returnReason,
				CancelReason:  cancelReason,
				YearMonth:     yearMonth,
				Amount:        amount})
	}

	return sellerPenalty
}

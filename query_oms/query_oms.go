package queryoms

import (
	"database/sql"
	"log"

	// I don't know what I am doing
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/sqltocsv"
)

// QueryOms queries oms and write result to csv file
func QueryOms(csvPath string) {

	db, err := sql.Open("mysql",
		"thomas:TH#)dec@)!&@tcp(178.22.70.43:3306)/oms_live_ir")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Connection failed")
		log.Fatal(err)
	} else {
		log.Println("Connection successful!")
	}

	query := `SELECT 
	soi_table.supplier_name
	,soi_table.order_nr
	,COUNT(soi_table.id_sales_order_item) count_item_penalty
	,COUNT(soi_table.id_sales_order_item)*2200000 order_penalty_amount
	,COUNT(soi_table.cancel_reason) cancel_reason_count
	,COUNT(soi_table.return_reason) return_reason_count
	,CONCAT(CASE WHEN MONTH(CURRENT_DATE()) = 1 THEN YEAR(CURRENT_DATE())-1 ELSE YEAR(CURRENT_DATE()) END
			,CASE WHEN MONTH(CURRENT_DATE()) = 1 THEN 12 ELSE MONTH(CURRENT_DATE())-1 END) 'year_month'
	
	FROM  
	(SELECT 
	isoi.id_sales_order_item
	,iso.order_nr
	,is1.name_en supplier_name
	,or1.name cancel_reason
	,or2.name return_reason
	
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
	
	GROUP BY isoi.id_sales_order_item) soi_table
	
	GROUP BY 
	soi_table.supplier_name
	,soi_table.order_nr`

	rows, _ := db.Query(query)

	err = sqltocsv.WriteFile(csvPath, rows)
	if err != nil {
		log.Fatal(err)
	}

}

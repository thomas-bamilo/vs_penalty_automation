package main

import (
	"log"
	"os"

	goemail "github.com/thomas-bamilo/vs_penalty_automation/goemail"
	queryoms "github.com/thomas-bamilo/vs_penalty_automation/queryoms"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// used for logging
	f, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	queryoms.QueryOms("seller_penalty.csv")

	goemail.GoEmail()

}

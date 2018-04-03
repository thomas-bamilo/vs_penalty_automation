package main

import (
	"log"
	"os"

	query_oms "github.com/tbeaudouin05/vs_penalty_automation/query_oms"
	send_email "github.com/tbeaudouin05/vs_penalty_automation/send_email"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	f, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	query_oms.QueryOms("file.csv")

	send_email.SendEmail("file.csv", []string{"mohammad.goudarzi@bamilo.com"})

}

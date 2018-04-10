package main

import (
	"log"
	"os"
	"strconv"

	"github.com/thomas-bamilo/vs_penalty_automation/goemail"
	joinscomstocsv "github.com/thomas-bamilo/vs_penalty_automation/joinscomstocsv"
	scitemid "github.com/thomas-bamilo/vs_penalty_automation/scitemid"
	sellerpenalty "github.com/thomas-bamilo/vs_penalty_automation/sellerpenalty"
)

func main() {

	// used for logging
	f, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// get seller penalty data from oms
	sellerPenalty := sellerpenalty.CreateSellerPenalty()

	// get omsID to only fetch coreesponding ID from Seller Center
	var omsIDStr string
	for i := 0; i < len(sellerPenalty)-1; i++ {
		omsIDStr += strconv.Itoa(sellerPenalty[i].OmsItemNumber) + ","
	}
	omsIDStr += strconv.Itoa(sellerPenalty[len(sellerPenalty)-1].OmsItemNumber)

	// get mapping OmsItemNumber and ScItemNumber
	scItemID := scitemid.CreateScItemID(omsIDStr)

	// join seller_penalty and sc_item_id tables on oms_item_number and write result to csv file in the same folder as the application
	joinscomstocsv.JoinScOmsToCsv(sellerPenalty, scItemID)

	goemail.GoEmail()

}

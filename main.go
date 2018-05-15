package main

import (
	"strconv"

	"github.com/thomas-bamilo/email/goemail"
	joinscomstocsv "github.com/thomas-bamilo/vspenaltyautomation/joinscomstocsv"
	scitemid "github.com/thomas-bamilo/vspenaltyautomation/scitemid"
	sellerpenalty "github.com/thomas-bamilo/vspenaltyautomation/sellerpenalty"
)

func main() {

	// get seller penalty data from oms
	sellerPenalty := sellerpenalty.CreateSellerPenalty()

	// get omsID to only fetch corresponding ID from Seller Center
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

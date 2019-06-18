package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"gopkg.in/Iwark/spreadsheet.v2"
)

//Request request
type Request struct {
	Title string     `json:"title"`
	Rows  [][]string `json:"rows"`
}

//Response response
type Response struct {
	Msg string `json:"message"`
	ID  string `json:"sheet-id"`
}

func prepareSheet(w http.ResponseWriter, r *http.Request) {

	rq := Request{}
	resp := Response{}

	req, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(req, &rq)
	if err != nil {
		log.Println("error unmarshaling request body: ", err)
	}

	id := createSheet(rq)

	w.Header().Set("Content-Type", "application/json")
	resp.Msg = "success"
	resp.ID = id
	rp, _ := json.Marshal(resp)

	w.Write(rp)
}

func authenticate() *http.Client {

	data, err := ioutil.ReadFile("client_secret.json")
	checkError(err)

	config, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/spreadsheets")
	checkError(err)

	client := config.Client(context.TODO())

	return client

}

func createSheet(req Request) string {

	client := authenticate()

	service := spreadsheet.NewServiceWithClient(client)

	// google
	sheetsService, err := sheets.New(client)

	ss, err := service.CreateSpreadsheet(spreadsheet.Spreadsheet{
		Properties: spreadsheet.Properties{
			Title: req.Title,
		},
	})
	checkError(err)

	sheet, err := ss.SheetByIndex(0)
	checkError(err)

	for _, row := range sheet.Rows {
		for _, cell := range row {
			fmt.Println(cell.Value)
		}
	}

	for i, r := range req.Rows {
		for z, c := range r {
			sheet.Update(i, z, c)
		}
	}

	// BatchUpdate
	err = sheet.Synchronize()
	checkError(err)

	reziseSheet(ss, sheetsService)

	return ss.ID

}

func reziseSheet(ss spreadsheet.Spreadsheet, service *sheets.Service) {

	id := ss.ID

	dr := sheets.DimensionRange{Dimension: "COLUMNS", StartIndex: 0, EndIndex: 999, SheetId: 0}
	ardr := sheets.AutoResizeDimensionsRequest{Dimensions: &dr}

	request := sheets.Request{AutoResizeDimensions: &ardr}

	requests := make([]*sheets.Request, 1)
	requests[0] = &request

	rb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}

	resp, err := service.Spreadsheets.BatchUpdate(id, rb).Context(context.TODO()).Do()
	if err != nil {
		log.Fatal(err)
	}

	if resp.ServerResponse.HTTPStatusCode != 200 {
		log.Println("Error updating sheet id: ", ss.ID)
	}

}

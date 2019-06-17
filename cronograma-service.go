package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
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
}

var ctx context.Context

func prepareSheet(w http.ResponseWriter, r *http.Request) {

	rq := Request{}
	resp := Response{}

	req, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(req, &rq)
	if err != nil {
		log.Println("error unmarshaling request body: ", err)
	}

	createSheet(rq)

	w.Header().Set("Content-Type", "application/json")
	resp.Msg = "success"
	rp, _ := json.Marshal(resp)

	w.Write(rp)
}

func createSheet(req Request) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	ctx = context.Background()

	client := config.Client(ctx, tok)

	service := spreadsheet.NewServiceWithClient(client)
	checkError(err)

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

}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
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

	resp, err := service.Spreadsheets.BatchUpdate(id, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", resp)

}

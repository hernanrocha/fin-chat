package stooq

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
)

func StooqHandler(s string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://stooq.com/q/l/?s=%s.us&f=sd2t2ohlcv&h&e=csv", s))
	reader := csv.NewReader(bufio.NewReader(resp.Body))
	_, err = reader.Read()
	if err != nil {
		log.Printf("Error reading header: %s \n", err)
		return fmt.Sprintf("Error obtaining info for %s", s), err
	}
	row, err := reader.Read()
	if err != nil || len(row) <= 4 {
		log.Printf("Error reading row: %s \n", err)
		return fmt.Sprintf("Error obtaining info for %s", s), err
	}
	return fmt.Sprintf("%s quote is $%s per share", s, row[3]), nil
}

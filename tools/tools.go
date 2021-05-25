package tools

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sensor-data-simulator/global"
	"strconv"
)

/*------------------------------*/

func SendJSON(resp http.ResponseWriter, obj interface{}) {

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
}

/*------------------------------*/

func GetLimitOffset(req *http.Request) (int, int, int) {
	qryParams := req.URL.Query()

	page := 1
	if _, ok := qryParams["page"]; ok {

		var err error
		page, err = strconv.Atoi(qryParams["page"][0])
		if err != nil {
			log.Printf("Error in page number: %v", err)
			page = 1
		}
		if page <= 0 {
			page = 1
		}
	}

	limit := global.RowsPerPage
	offset := (page - 1) * limit

	return limit, offset, page
}

/*------------------------------*/

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	return nil
}

func ReadAll(rc io.ReadCloser) ([]byte, error) {
	defer rc.Close()

	if cb, ok := rc.(*ClosingBuffer); ok {
		return cb.Bytes(), nil
	}

	return ioutil.ReadAll(rc)
}

/*------------------------------*/

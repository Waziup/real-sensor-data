package api

import (
	"net/http"
	"real-sensor-data/global"
	"real-sensor-data/tools"

	routing "github.com/julienschmidt/httprouter"
)

/*-------------*/

/*
* This function implements GET /dataCollection/status
 */

func GetDataCollectionStatus(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	tools.SendJSON(resp, global.DataCollectorProgress)
}

/*-------------*/

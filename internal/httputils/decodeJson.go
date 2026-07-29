package httputils

import (

	"encoding/json"
	"net/http"
)


func DecodeJSONReq(r *http.Request, dts any){
	// dts -> pointer to where to write to after marshall
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(dts)
}
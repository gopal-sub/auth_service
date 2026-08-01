package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func DecodeJSONReq(r *http.Request, dts any) error {
	// dts -> pointer to where to write to after marshall

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(dts)
	

	if errors.Is(err, SyntaxError) {
		return JsonSyntaxError
	}
	if err != nil {
		if strings.HasPrefix(err.Error(), "json: unknown field"){
			return DisallowedFieldsError
		}
		return err
	}
	return nil

}

func DecodeAndWriteError(w http.ResponseWriter, r *http.Request, dts any) error{
	err := DecodeJSONReq(r, dts)
	if errors.Is(err, JsonSyntaxError){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Incorrect body",
		})
		return err
	}
	if errors.Is(err, DisallowedFieldsError){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid fields",
		})
		return err
	}
	fmt.Println(err)
	if err != nil{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "bad data from client",
		})
		return err
	}
	return nil

}


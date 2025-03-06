package helpers

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter,resp any , status int)error{
	w.Header().Set("content-type","application/json")
	w.WriteHeader(status)
	data ,err :=json.Marshal(resp)
	if err !=nil {
		return err; 
	}
 _,err = w.Write(data)
 return err; 
}


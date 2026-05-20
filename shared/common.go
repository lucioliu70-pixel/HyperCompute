package shared

import (
  "encoding/json"
  "net/http"
)

type ErrorBody struct { Error struct { Message string `json:"message"`; Type string `json:"type"`; Code string `json:"code"`} `json:"error"` }

func WriteError(w http.ResponseWriter, status int, msg, code string){
  w.Header().Set("Content-Type","application/json")
  w.WriteHeader(status)
  var e ErrorBody
  e.Error.Message=msg; e.Error.Type="invalid_request_error"; e.Error.Code=code
  _ = json.NewEncoder(w).Encode(e)
}

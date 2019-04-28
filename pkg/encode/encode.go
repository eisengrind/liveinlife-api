package encode

import "net/http"

//Encoder for api http endpoints
type Encoder interface {
	//Encode an interface{} to a data format.
	//Can only be called once in a request because it sets data format specific http headers
	Encode(w http.ResponseWriter, v interface{}) error
}

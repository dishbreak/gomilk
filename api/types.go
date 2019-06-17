package api

/*
Envelope represents the API response from RTM
*/
type Envelope struct {
	Rsp  interface{}
	Stat string
	Err  struct {
		
	}
}


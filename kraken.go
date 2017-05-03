package kraken

import "net/http"

// Kraken main client for the API.
type Kraken struct {
	Client *http.Client
}

// Init initialize the client instance.
func (k *Kraken) Init() {
	k.Client = &http.Client{}
}

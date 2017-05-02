package kraken_api

import "net/http"

type Kraken struct {
	Client *http.Client
}

func (k *Kraken) Init() {
	k.Client = &http.Client{}
}

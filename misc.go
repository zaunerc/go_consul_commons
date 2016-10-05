// Package consul provides utility functions for the consul key/value store.
package go_consul_commons

import (
	"fmt"
	"github.com/zaunerc/consul/api"
	"sync"
)

// internalGet returns nil if key does not exist or if
// there is an error during lookup. In both cases a
// message is logged.
func InternalGet(kv *api.KV, key string) []byte {

	kvp, _, err := kv.Get(key, nil)

	if err != nil {
		fmt.Printf("Error while reading key >%s<. Returning nil as key value.\n", err)
		return nil
	} else if kvp == nil {
		fmt.Printf("Key >%s< does not exist in registry. Returning nil as key value.", key)
		return nil
	}

	return kvp.Value
}

var consulClients = make(map[string]*api.Client)
var mutex = &sync.Mutex{}

// getConsulClientForUrl returns a Consul HTTP client specific
// for one URL. The HTTP client will pool and reuse idle
// connections to Consul.
// See github.com/hashicorp/consul/api.DefaultConfig()
func GetConsulClientForUrl(consulUrl string) (*api.Client, error) {

	mutex.Lock()
	defer mutex.Unlock()

	if consulClients[consulUrl] == nil {
		config := api.DefaultConfig()
		config.Address = consulUrl

		var err error
		consulClients[consulUrl], err = api.NewClient(config)

		if err != nil {
			fmt.Printf("Error while creating consul HTTP client: %s", err)
			return nil, err
		}
	}

	return consulClients[consulUrl], nil
}

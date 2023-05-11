package defaults

import "sync"

const (
	defaultVerbose         = false
	defaultCountTokens     = true
	defaultCache           = "false"
	defaultCallbackManager = "false"
)

var (
	mu              sync.Mutex
	verbose         = defaultVerbose
	countTokens     = defaultCountTokens
	cache           = defaultCache
	callbackManager = defaultCallbackManager
)

func SetVerbose(v bool) {
	mu.Lock()
	defer mu.Unlock()
	verbose = v
}

func GetDefaultVerbose() bool {
	mu.Lock()
	defer mu.Unlock()
	return verbose
}

func SetCache(c string) {
	mu.Lock()
	defer mu.Unlock()
	cache = c
}

func GetDefaultCache() string {
	mu.Lock()
	defer mu.Unlock()
	return cache
}

func SetCallbackManager(c string) {
	mu.Lock()
	defer mu.Unlock()
	callbackManager = c
}

func GetDefaultCallbackManager() string {
	mu.Lock()
	defer mu.Unlock()
	return callbackManager
}

func SetCountTokens(c bool) {
	mu.Lock()
	defer mu.Unlock()
	countTokens = c
}

func GetDefaultCountTokens() bool {
	mu.Lock()
	defer mu.Unlock()
	return countTokens
}

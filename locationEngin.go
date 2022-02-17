package geo2city

import "sync"

type locationEngin struct {
	sync.RWMutex
	engine map[string]*LocationParserEngine
}

func newLocationEngin() *locationEngin {
	return &locationEngin{
		engine: make(map[string]*LocationParserEngine),
	}
}

func (le *locationEngin) Load(key string) (value *LocationParserEngine, ok bool) {
	le.RLock()
	result, ok := le.engine[key]
	le.RUnlock()
	return result, ok
}

func (le *locationEngin) Delete(key string) {
	le.Lock()
	delete(le.engine, key)
	le.Unlock()
}

func (le *locationEngin) Store(key string, value *LocationParserEngine) {
	le.Lock()
	le.engine[key] = value
	le.Unlock()
}

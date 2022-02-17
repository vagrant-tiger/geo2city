package geo2city

import "sync"

type LocationEngin struct {
	sync.RWMutex
	engine map[string]*LocationParserEngine
}

func NewLocationEngin() *LocationEngin {
	return &LocationEngin{
		engine: make(map[string]*LocationParserEngine),
	}
}

func (le *LocationEngin) Load(key string) (value *LocationParserEngine, ok bool) {
	le.RLock()
	result, ok := le.engine[key]
	le.RUnlock()
	return result, ok
}

func (le *LocationEngin) Delete(key string) {
	le.Lock()
	delete(le.engine, key)
	le.Unlock()
}

func (le *LocationEngin) Store(key string, value *LocationParserEngine) {
	le.Lock()
	le.engine[key] = value
	le.Unlock()
}

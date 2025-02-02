package validators

import (
	"log"
	"sync"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

// Registry - holds all registered Validators
var Registry = NewDefaultRegistry()

func NewDefaultRegistry() defaultRegistry {
	return defaultRegistry{
		&sync.Mutex{},
		make(map[string]types.Validator),
	}
}

type defaultRegistry struct {
	*sync.Mutex
	Data map[string]types.Validator
}

// Add - update the registry in a thread-safe way. Called in init() functions
func (r *defaultRegistry) Add(v types.Validator) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.Data[v.Code]; ok {
		log.Fatalf("Validator code %v already exist.", v.Code)
	}
	r.Data[v.Code] = v
}

func (r defaultRegistry) Len() int {
	return len(r.Data)
}

func (r defaultRegistry) All() map[string]types.Validator {
	return r.Data
}

func (r defaultRegistry) Get(k string) (types.Validator, bool) {
	v, ok := r.Data[k]
	return v, ok
}

// TestRegistry - register all test structs
var TestRegistry = NewTestRegistry()

func NewTestRegistry() *testRegistry {
	return &testRegistry{
		&sync.Mutex{},
		[]types.ValidatorTest{},
	}
}

type testRegistry struct {
	*sync.Mutex
	Data []types.ValidatorTest
}

// Add - update the test registry in a thread-safe way. Called in init() functions
func (t *testRegistry) Add(v types.ValidatorTest) {
	t.Lock()
	defer t.Unlock()

	t.Data = append(t.Data, v)
}

func (t *testRegistry) All() []types.ValidatorTest {
	return t.Data
}

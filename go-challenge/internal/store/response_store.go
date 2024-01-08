package store

import (
	"encoding/json"
	"fmt"
	"github.com/cloudx-labs/challenge/internal/model/response"
	"github.com/samber/lo"
	"strings"
	"sync"
)

type ResponseStore struct {
	Store map[string]int
	mu    sync.Mutex
}

func NewResponseStore() *ResponseStore {
	return &ResponseStore{
		Store: make(map[string]int),
		mu:    sync.Mutex{},
	}
}

func (r *ResponseStore) Add(parentGroup ParentGroup) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s,%s", parentGroup.Source, parentGroup.Destination)
	r.Store[key]++
}

func (r *ResponseStore) ToResponse() ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := lo.MapToSlice(r.Store, func(k string, v int) response.Response {
		keys := strings.Split(k, ",")
		return response.NewResponse(keys[0], keys[1], v)
	})

	bytes, err := json.Marshal(&res)
	if err != nil {
		return nil, err
	}

	r.Store = make(map[string]int)

	return bytes, nil
}

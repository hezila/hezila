package hezila

package wechat

import (
	"net/http"
	"sync"

	"github.com/hezila/hezila/cache"
)

// Hezila
type Hezila struct {
	
}

// Config for user
type Config struct {
	Cache cache.Cache
}
// NewHezila init
func NewHezila(cfg *Config) *Hezila {
	
}

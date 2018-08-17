package creditd

import (
	"sync"

	"gsxt/credit"
)

type Pool struct {
	sync.Mutex
	max  int
	pool map[string]chan credit.Credit
}

// NewPool creates a new pool of Clients.
func NewPool(max int) *Pool {
	return &Pool{
		max:  max,
		pool: make(map[string]chan credit.Credit),
	}
}

func (p *Pool) getPool(name string) chan credit.Credit {
	p.Lock()
	defer p.Unlock()
	pool, ok := p.pool[name]
	if !ok {
		pool = make(chan credit.Credit, p.max)
		p.pool[name] = pool
	}
	return pool
}

// Borrow a Client from the pool.
func (p *Pool) Borrow(name string) credit.Credit {
	pool := p.getPool(name)
	var c credit.Credit
	select {
	case c = <-pool:
	default:
		c = credit.MustGet(name)
	}
	return c
}

// Return returns a Client to the pool.
func (p *Pool) Return(name string, c credit.Credit) {
	pool := p.getPool(name)
	select {
	case pool <- c:
	default:
		// let it go, let it go...
	}
}

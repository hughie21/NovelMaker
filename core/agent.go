// Description: Agent is a manager for all plugins, it can execute plugin in parallel and return the result.
// Author: Hughie21
// Date: 2024-11-21
// license that can be found in the LICENSE file.

// The program is use the agent to schedule the plugin. The return value of the plugin is stored in the result channel.
// There will be a pool of transaction that used to store the plugin that is running. The transaction will be executed
// in parallel.
package core

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

const (
	TransactionStateWaiting = iota
	TransactionStateRunning
)

// the result carry the data and error of the plugin
type Result struct {
	data interface{}
	err  error
}

type Transaction struct {
	executor Pluginer
	args     []interface{}
	result   chan Result
	wg       *sync.WaitGroup
	state    int
}

// Agent is a manager for all plugins, so there is map that store all the registered plugins.
// The context that carraied by the agent is used to cancel the plugin when the app is shutdown.
// The timeout that represent the time to wait for the plugin to find a empty transaction.
type Agent struct {
	plugins map[string]Pluginer
	ctx     context.Context
	cancel  context.CancelFunc
	timeout time.Duration
	pool    []Transaction
}

func (r *Result) Data() interface{} {
	return r.data
}

func (r *Result) Err() error {
	return r.err
}

// NewAgent create a new agent with the size of the transaction pool and the timeout for the plugin to execute
// the size of the pool and the timeout are defined at the config file.
func NewAgent(size int, timeout time.Duration) *Agent {
	agt := Agent{
		plugins: make(map[string]Pluginer),
	}
	tmpDir := filepath.Join(currentPath, "tmp")
	reader := NewEpubReader(tmpDir)
	writer := NewEpubWriter(tmpDir)

	// the writer and the reader are the default plugin that is registered.
	// the other plugin is the extension of these two plugin.
	agt.RegisterPlugin("writer", writer)
	agt.RegisterPlugin("reader", reader)

	agt.ctx, agt.cancel = context.WithCancel(context.Background())

	agt.timeout = timeout

	agt.pool = make([]Transaction, size)
	for i := 0; i < size; i++ {
		// Initialize the transaction
		agt.pool[i] = *NewTransaction()
	}
	return &agt
}

func (agt *Agent) RegisterPlugin(name string, plugin Pluginer) error {
	if _, ok := agt.plugins[name]; ok {
		return errors.New("plugin already registered")
	}
	agt.plugins[name] = plugin
	return nil
}

func (agt *Agent) Exec(name string, args ...interface{}) Result {
	plugin, ok := agt.plugins[name]
	if !ok {
		return Result{err: fmt.Errorf("plugin<%s> not found", name)}
	}

	ctx, cancel := context.WithTimeout(agt.ctx, agt.timeout*time.Second)
	defer cancel()

	// Wait for a empty transaction, if there is no empty transaction, the agent will wait for the timeout.
	for {
		for i := 0; i < len(agt.pool); i++ {
			if agt.pool[i].state == TransactionStateWaiting {
				agt.pool[i].executor = plugin
				agt.pool[i].state = TransactionStateRunning
				agt.pool[i].exec(ctx, args...)
				return agt.pool[i].wait()
			}
		}
	}
}

func (agt *Agent) Close() {
	agt.cancel()
}

func NewTransaction() *Transaction {
	return &Transaction{result: make(chan Result, 1), wg: &sync.WaitGroup{}, state: TransactionStateWaiting}
}

// wait() will wait for the plugin to finish and return the result.
func (t *Transaction) wait() Result {
	t.wg.Wait()
	t.state = TransactionStateWaiting
	// After return the result, close the result channel to release the resources.
	defer close(t.result)
	return <-t.result
}

func (t *Transaction) exec(ctx context.Context, args ...interface{}) {
	t.args = args
	t.result = make(chan Result, 1)
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		result, err := t.executor.Run(ctx, args...)
		res := Result{data: result, err: err}
		t.result <- res
	}()
}

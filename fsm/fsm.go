package fsm

import (
	"fmt"

	"github.com/hashicorp/raft"
)

// FSM is a Raft Finite State Machine (FSM).
type FSM struct {
}

// New creates a new Finite State Machine (FSM).
func New() *FSM {
	return &FSM{}
}

// Apply log is invoked once a Raft log entry is committed. It returns 
// a value which will be made available in the ApplyFuture returned by 
// Raft.Apply method if that method was called on the same Raft node as 
// the FSM.
func (f *FSM) Apply(l *raft.Log) interface{} {
	fmt.Printf("%v\n", l.Data)
	return nil
}

// Snapshot will be called duiring make snapshot
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
}



// restore from FSMSnapshot
// TODO
func (f *fsm) Restore(rc io.ReadCloser) error {
	f.logger.Printf("Restore snapshot from FSMSnapshot")
	defer rc.Close()
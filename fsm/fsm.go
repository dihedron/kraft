package fsm

import (
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/dihedron/kraft/log"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var ErrNotFound error = fmt.Errorf("item not found")

// FSM is a Raft Finite State Machine (FSM).
type FSM struct {
	mutex sync.RWMutex
	VMs   map[string]int `json:"vms,omitempty" yaml:"vms,omitempty"`
}

// New creates a new Finite State Machine (FSM).
func New() *FSM {
	return &FSM{}
}

// Get retrieves an item from the internal store, if available; it does
// not have to go through the Leader, since all Followers have the same
// data as the Leader and there is no mutation involved.
func (f *FSM) Get(key string) (int, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	if value, ok := f.VMs[key]; ok {
		return value, nil
	}
	return 0, ErrNotFound
}

// Apply log is invoked once a Raft log entry is committed. It returns
// a value which will be made available in the ApplyFuture returned by
// Raft.Apply method if that method was called on the same Raft node as
// the FSM.
func (f *FSM) Apply(l *raft.Log) interface{} {
	fmt.Printf("%v\n", l.Data)
	// var c command
	// if err := json.Unmarshal(l.Data, &c); err != nil {
	// 	panic("failed to unmarshal raft log")
	// }

	// switch strings.ToLower(c.Op) {
	// case "set":
	// 	return f.applySet(c.Key, c.Value)
	// case "delete":
	// 	return f.applyDelete(c.Key)
	// default:
	// 	panic("command type not support")
	// }
	return nil
}

// Snapshot will be called to create a snapshot of the FSM internal status;
// it must return an object that is capable of saving the FSM internal
// status while there MAY be concurrent calls to Apply; in our simplistic
// case, the FSM is capable of snapshotting itself.
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return f, nil
}

// Restore is called to erase the current FSM internal status and to load
// it with external data as per the provided reader.
func (f *FSM) Restore(reader io.ReadCloser) error {
	defer reader.Close()
	log.L.Info("restoring from snapshot...")
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.L.Error("error reading snapshot data", zap.Error(err))
		return err
	}
	log.L.Info("snapshot data read")
	m := map[string]int{}
	if err = yaml.Unmarshal(data, &m); err != nil {
		log.L.Error("error unmarshaling snapshot data from YAML", zap.Error(err))
		return err
	}
	log.L.Info("replacing data from snapshot")
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.VMs = m
	return nil
}

// Persist should write all items to the given sink and then close it;
// it might acquire a lock to prevent concurrent updates, since it this
// function could be called with concurrent calls to Apply
func (f *FSM) Persist(sink raft.SnapshotSink) error {
	defer sink.Close()
	log.L.Info("persisting snapshot...")
	snapshot, err := f.takeSnapshot()
	if err != nil {
		log.L.Error("error marshalling snapshot to YAML", zap.Error(err))
		return err
	}
	log.L.Info("snapshot taken and marshalled to YAML")
	if _, err := sink.Write(snapshot); err != nil {
		log.L.Error("error writing snapshot to sink", zap.Error(err))
		return err
	}
	log.L.Info("snapshot written out")
	return nil
}

// Release is invoked when we are finished with the snapshot; it could
// be used to release a lock that might have been acquired when starting
// to Persist (to avoid problems with concurrent calls to Apply).
func (f *FSM) Release() {
	log.L.Info("releasing lock...")
}

// takeSNapshot is a critical section around the copying of the internal
// store to YAML; it acquires a read lock, marshalls the data and then
// releases the as soon as possible.
func (f *FSM) takeSnapshot() ([]byte, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return yaml.Marshal(f.VMs)
}

package fsm

import "github.com/hashicorp/raft"

type Snapshot struct {
}

// Persist should write all items to the given sink and then close it;
// might acquire a lock to prevent concurrent updates.
func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	defer sink.Close()

	return nil
}

// Release is invoked when we are finished with the snapshot; it could
// be used to release a lock that might have been acquired when starting
// to persist.
func (s *Snapshot) Release() {
}

// type fsmSnapshot struct {
// 	db     DB
// 	logger *log.Logger
// }

// // Persist data in specific type
// // kv item serialize in google protubuf
// func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
// 	f.logger.Printf("Persist action in fsmSnapshot")
// 	defer sink.Close()

// 	ch := f.db.SnapshotItems()

// 	keyCount := 0

// 	// read kv item from channel
// 	for {
// 		buff := proto.NewBuffer([]byte{})

// 		dataItem := <-ch
// 		item := dataItem.(*KVItem)

// 		if item.IsFinished() {
// 			break
// 		}

// 		// create new protobuf item
// 		protoKVItem := &ProtoKVItem{
// 			Key:   item.key,
// 			Value: item.value,
// 		}

// 		keyCount = keyCount + 1

// 		// encode message
// 		buff.EncodeMessage(protoKVItem)

// 		if _, err := sink.Write(buff.Bytes()); err != nil {
// 			return err
// 		}
// 	}
// 	f.logger.Printf("Persist total %d keys", keyCount)

// 	return nil
// }

// func (f *fsmSnapshot) Release() {
// 	f.logger.Printf("Release action in fsmSnapshot")
// }

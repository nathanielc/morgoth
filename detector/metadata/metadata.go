// Simple storage of metadata about metrics
package metadata

import (
	//"github.com/golang/glog"
	"errors"
	"github.com/cznic/kv"
	metric "github.com/nvcook42/morgoth/metric/types"
	"path"
)

type MetadataStore struct {
	db *kv.DB
}

func New(dir, detectorID string) (*MetadataStore, error) {

	dbPath := path.Join(dir, detectorID + ".db")
	//Init KV database
	opts := kv.Options{}
	db, err := kv.Create(dbPath, &opts)
	if err != nil {
		db, err = kv.Open(dbPath, &opts)
		if err != nil {
			return nil, err
		}
	}
	if db == nil {
		return nil, errors.New("DB failed to initialize for unknown reason")
	}

	ms := new(MetadataStore)
	ms.db = db

	return ms, nil
}

func (self *MetadataStore) StoreDoc(metric metric.MetricID, doc []byte) {
	self.db.Set([]byte(metric), doc)
}

func (self *MetadataStore) GetDoc(metric metric.MetricID) []byte {
	doc, err := self.db.Get(nil, []byte(metric))
	if err != nil {
		return []byte{}
	}
	return doc
}


func (self *MetadataStore) Close() {
	self.db.Close()
}

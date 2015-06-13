// Simple storage of metadata about metrics
package metadata

import (
	//"github.com/golang/glog"
	"errors"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/cznic/kv"
	metric "github.com/nathanielc/morgoth/metric/types"
	"path"
)

type MetadataStore interface {
	StoreDoc(metric.MetricID, []byte)
	GetDoc(metric.MetricID) []byte
	Close()
}

type MetadataStoreT struct {
	db *kv.DB
}

func New(dir, detectorID string) (MetadataStore, error) {

	dbPath := path.Join(dir, detectorID+".db")
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

	ms := new(MetadataStoreT)
	ms.db = db

	return ms, nil
}

func (self *MetadataStoreT) StoreDoc(metric metric.MetricID, doc []byte) {
	self.db.Set([]byte(metric), doc)
}

func (self *MetadataStoreT) GetDoc(metric metric.MetricID) []byte {
	doc, err := self.db.Get(nil, []byte(metric))
	if err != nil {
		return []byte{}
	}
	return doc
}

func (self *MetadataStoreT) Close() {
	self.db.Close()
}

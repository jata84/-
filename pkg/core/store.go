package core

import (
	"strings"
)

type DataStore struct {
	datastore map[string]interface{}
}

func (s *DataStore) Set(key string, value interface{}) {
	s.datastore[key] = value
}

func (s *DataStore) Get(key string) interface{} {
	return s.datastore[key]
}

func (ds *DataStore) GetValueByName(name string) interface{} {
	keys := strings.Split(name, ".")
	value := ds.datastore

	for _, k := range keys {
		if val, ok := value[k]; ok {
			if valMap, ok := val.(map[string]interface{}); ok {
				value = valMap
			} else {
				return val
			}
		} else {
			return nil
		}
	}

	return value
}
func NewDataStore(data map[string]interface{}) *DataStore {

	if data == nil {
		return &DataStore{
			datastore: make(map[string]interface{}),
		}
	} else {
		return &DataStore{
			datastore: data,
		}
	}

}

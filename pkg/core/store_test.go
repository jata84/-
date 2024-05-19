package core

import (
	"testing"
)

func TestSetAndGet(t *testing.T) {
	store := NewDataStore(nil)
	store.Set("key1", "value1")
	store.Set("key2", 42)

	value1 := store.Get("key1")
	value2 := store.Get("key2")
	value3 := store.Get("key3")

	if value1 != "value1" {
		t.Errorf("Expected 'value1', but got '%v'", value1)
	}

	if value2 != 42 {
		t.Errorf("Expected 42, but got '%v'", value2)
	}

	if value3 != nil {
		t.Errorf("Expected nil, but got '%v'", value3)
	}
}

func TestDataStoreGetValueByName(t *testing.T) {
	store := NewDataStore(nil)
	store.Set("m", "Hello")
	store.Set("a", map[string]interface{}{"b": "Hello"})
	store.Set("c", map[string]interface{}{"b": map[string]interface{}{"c": "Hello"}})

	value := store.GetValueByName("a.b")
	if value != "Hello" {
		t.Errorf("Expected Hello on a key")
	}
	value2 := store.GetValueByName("a.b.c")
	if value2 != "Hello" {
		t.Errorf("Expected Hello on a.b key")
	}
	value3 := store.GetValueByName("m")
	if value3 != "Hello" {
		t.Errorf("Expected Hello on a key")
	}
}

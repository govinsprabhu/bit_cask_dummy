package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/govinsprabhu/kv_store/utils"
)

var (
	store = make(map[string]int64)
	mu    sync.RWMutex
)

func init_kvstore(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	position := int64(0)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("line:", line)
		if !strings.HasSuffix(line, "*") {
			fmt.Println("inside has sufficx line:")
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				store[key] = position
			}
		} else {
			key := strings.TrimSuffix(line, "=*")
			delete(store, key)
		}
		position += int64(len(line) + 1)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %v", err)
	}
	fmt.Printf("Initialized kv store with %d keys\n", len(store))
	return nil
}

func Get(key string) (int64, error) {
	mu.RLock()
	defer mu.RUnlock()
	value, exists := store[key]
	if !exists {
		return 0, fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

func Put(key, value string) {
	mu.Lock()
	defer mu.Unlock()
	position, err := utils.WriteKeyValueToFile("kv_store.txt", key, value)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	store[key] = position
}

func Delete(key string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := store[key]; !exists {
		return fmt.Errorf("key not found: %s", key)
	}
	utils.MarkDelete("kv_store.txt", key)
	delete(store, key)
	return nil
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}
	position, err := Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	key_value_pair, err := utils.ReadFromFileAtPosition("kv_store.txt", position)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Write([]byte(key_value_pair))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	Put(key, value)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Key %s added/updated", key)))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}
	err := Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Key %s deleted", key)))
}

func main() {
	init_kvstore("kv_store.txt")
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/put", putHandler)
	http.HandleFunc("/delete", deleteHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

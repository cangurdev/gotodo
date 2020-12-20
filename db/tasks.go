package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

type Task struct {
	Key   int
	Value Content
}
type Content struct {
	IsDone bool   `json:"IsDone"`
	Parent string `json:"Parent"`
	Text   string `json:"Text"`
	Tag    string `json:"Tag"`
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func CreateTask(task, tag string) error {

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		key := itob(int(id64))
		config := &Content{false, "main", task, tag}
		dataBytes, _ := json.Marshal(config)
		return b.Put(key, dataBytes)
	})
	if err != nil {
		return err
	}
	return nil
}
func UpdateTask(key int, task string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		v := b.Get(itob(key))
		data := unmarshall(v)
		data.Text = task
		dataBytes, _ := json.Marshal(data)
		err := b.Put(itob(key), dataBytes)
		return err
	})
}

func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			content := unmarshall(v)
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: content,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

func DoneTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		v := b.Get(itob(key))
		data := unmarshall(v)
		data.IsDone = true
		dataBytes, _ := json.Marshal(data)
		err := b.Put(itob(key), dataBytes)
		return err
	})
}

//converts integer to binary
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

//converts binary to integer
func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

//converts []byte to Content struct
func unmarshall(v []byte) Content {
	var data Content
	err := json.Unmarshal(v, &data)
	if err != nil {
		fmt.Println("kawga")
	}
	return data
}

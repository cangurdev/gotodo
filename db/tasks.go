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
	IsDone      bool   `json:"IsDone"`
	Parent      string `json:"Parent"`
	Text        string `json:"Text"`
	IsImportant bool   `json:"IsImportant"`
	Due         string `json:"Due"`
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

//CreateTask add the task to the bucket
func CreateTask(task, due, parent string, isImportant bool) error {

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		key := itob(int(id64))
		config := &Content{false, parent, task, isImportant, due}
		dataBytes, _ := json.Marshal(config)
		return b.Put(key, dataBytes)
	})
	if err != nil {
		return err
	}
	return nil
}

//UpdateTask updates the task
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

//AllTasks returns all tasks in the bucket
func AllTasks() ([]Task, error) {
	var tasks []Task
	var doneTasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			content := unmarshall(v)
			if content.IsDone {
				doneTasks = append(doneTasks, Task{
					Key:   btoi(k),
					Value: content,
				})
			} else {
				tasks = append(tasks, Task{
					Key:   btoi(k),
					Value: content,
				})
			}
		}
		tasks = append(tasks, doneTasks...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

//DeleteTask removes the task from bucket
func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

//DoTask makes true isDone property of task that given key
func DoTask(key int) error {
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
		fmt.Println("kawga", err)
	}
	return data
}

package models_test

import (
	"log"
	"os"
	"testing"

	"github.com/1000k/gj/models"
)

const dbfile = "./gj_test.db"

func setup() {
	_, err := models.ConnectDb(dbfile)
	if err != nil {
		log.Fatalf("Cannot connect database: '%v'", err)
	}
	models.ResetDb(dbfile)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestSave(t *testing.T) {
	setup()

	_, err := models.NewMessage("Me", "You", "Hello")
	if err != nil {
		t.Errorf("Failed to save a message: '%v'", err)
	}
}

func TestFindMessages(t *testing.T) {
	setup()

	models.NewMessage("Me", "You", "Hello")
	models.NewMessage("Me2", "You2", "Hello2")
	res, err := models.FindMessages()
	if len(res) != 2 {
		t.Errorf("Number of result set does not match (expected: 1, actual: %v)", len(res))
	}
	if err != nil {
		t.Errorf("Failed to find messages: '%v'", err)
	}
}

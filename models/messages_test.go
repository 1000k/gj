package models_test

import (
	"log"
	"os"
	"testing"

	"github.com/1000k/gj/models"
)

func TestMain(m *testing.M) {
	_, err := models.ConnectDb("./gj_test.db")
	if err != nil {
		log.Fatalf("Cannot connect database: '%v'", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestSave(t *testing.T) {
	_, err := models.NewMessage("Me", "You", "Hello")
	if err != nil {
		t.Errorf("Failed to save with error '%v'", err)
	}
}

func TestFind(t *testing.T) {
	t.Skip("not implemented yet")
}

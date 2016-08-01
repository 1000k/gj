package models_test

import (
	"log"
	"os"
	"reflect"
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

func TestNewMessage(t *testing.T) {
	setup()

	params := models.NewMessageParams{From: "Me", To: "You", Message: "Hello"}
	_, err := models.NewMessage(params)
	if err != nil {
		t.Errorf("Failed to save a message: '%v'", err)
	}
}

func TestFindMessages(t *testing.T) {
	setup()

	params := models.NewMessageParams{From: "Me", To: "You", Message: "Hello"}
	models.NewMessage(params)
	models.NewMessage(params)
	res, err := models.FindMessages()
	if len(res) != 2 {
		t.Errorf("Number of result set does not match (expected: 1, actual: %v)", len(res))
	}
	if err != nil {
		t.Errorf("Failed to find messages: '%v'", err)
	}
}

func TestFindRanking(t *testing.T) {
	setup()

	params := models.NewMessageParams{From: "Miles", To: "Coltrane", Message: "foobar", CreatedAt: "2016-08-01 00:00:00"}
	for i := 0; i < 5; i++ {
		models.NewMessage(params)
	}

	params2 := models.NewMessageParams{From: "Coltrane", To: "Monk", Message: "foobar", CreatedAt: "2016-08-02 00:00:00"}
	for i := 0; i < 3; i++ {
		models.NewMessage(params2)
	}

	expected := []models.RankingItem{{Name: "Coltrane", Count: 5}, {Name: "Monk", Count: 3}}

	res, err := models.FindRanking("201608")

	if err != nil {
		t.Errorf("Failed to find ranking: '%v'", err)
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("Result does not match to the expected. expected: %v, actual: %v", expected, res)
	}
}

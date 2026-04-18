package db

import (
	"fmt"
	"testing"
)

func TestStats(t *testing.T) {
	err := InitDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	CreateUser("Pudge", "123")
	CreateUser("Sniper", "456")
	AddPushups(1, 10)
	AddPushups(1, 15)
	AddPushups(2, 50)
	stats, err := GetStats(0)
	if err != nil {
		t.Fatal(err)
	}
	for _, user := range stats {
		fmt.Printf("User: %s (Color: %s)\n", user.Username, user.Color)
		for _, p := range user.Points {
			fmt.Printf("  Date: %s, Cumulative Total: %d\n", p.Date, p.Value)
		}
	}
}

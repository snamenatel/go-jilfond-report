package main

import "testing"

func TestGetMonthTranslate(t *testing.T) {
	res := getMonthTranslate("2021-12")
	expected := "дек."

	if res != expected {
		t.Errorf("res %s, expected %s", res, expected)
	}
}

func TestMinutesToCost(t *testing.T) {
	res := minutesToCost(30)
	expected := "350 руб."

	if res != expected {
		t.Errorf("res %s, expected %s", res, expected)
	}
}
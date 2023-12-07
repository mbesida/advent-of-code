package main

import (
	"testing"
)

func TestOnePair(t *testing.T) {
	rank := "32T3K"
	hand := NewHand(rank)
	if hand.tpe != OnePair {
		t.Fatal(hand)
	}
}

func TestTwoPair(t *testing.T) {
	rank := "23432"
	hand := NewHand(rank)
	if hand.tpe != TwoPair {
		t.Fatal(hand)
	}
}

func TestThreeKind(t *testing.T) {
	rank := "TTT98"
	hand := NewHand(rank)
	if hand.tpe != ThreeKind {
		t.Fatal(hand)
	}
}

func TestFullHouse(t *testing.T) {
	rank := "23332"
	hand := NewHand(rank)
	if hand.tpe != FullHouse {
		t.Fatal(hand)
	}
}

func TestFourKind(t *testing.T) {
	rank := "AA8AA"
	hand := NewHand(rank)
	if hand.tpe != FourKind {
		t.Fatal(hand)
	}
}

func TestFiveKind(t *testing.T) {
	rank := "AAAAA"
	hand := NewHand(rank)
	if hand.tpe != FiveKind {
		t.Fatal(hand)
	}
}

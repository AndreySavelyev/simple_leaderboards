package engine

import (
	"fmt"
	"testing"
)

func TestRewardForPosition(t *testing.T) {
	var ranks = [][]int{
		{1, 1000},
		{2, 500},
		{3, 250},
		{4, 50},
		{5, 50},
		{50, 50},
		{51, 0},
	}
	for _, rank := range ranks {
		testname := fmt.Sprintf("%d,%d", rank[0], rank[1])
		t.Run(testname, func(t *testing.T) {
			rew := RewardForPosition(rank[0])
			if rew != rank[1] {
				t.Errorf("got %d, expected %d", rew, rank[1])
			}
		})
	}
}

// Copyright (C) 2018 go-dappley authors
//
// This file is part of the go-dappley library.
//
// the go-dappley library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-dappley library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-dappley library.  If not, see <http://www.gnu.org/licenses/>.
//

package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const DefaultMaxProducersIfNoProducersGiven = 5
const DefaultTimeBetweenBlockIfNoneGiven = 15

func TestDynasty_NewDynasty(t *testing.T) {
	dynasty := NewDynasty([]string{}, DefaultMaxProducersIfNoProducersGiven, DefaultTimeBetweenBlockIfNoneGiven)
	assert.Empty(t, dynasty.producers)
}

func TestDynasty_AddProducer(t *testing.T) {
	tests := []struct {
		name     string
		maxPeers int
		input    string
		expected []string
	}{
		{
			name:     "ValidInput",
			maxPeers: 3,
			input:    "dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
			expected: []string{"dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf"},
		},
		{
			name:     "MinerExceedsLimit",
			maxPeers: 0,
			input:    "dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
			expected: []string{},
		},
		{
			name:     "InvalidInput",
			maxPeers: 3,
			input:    "m1",
			expected: []string{},
		},
		{
			name:     "EmptyInput",
			maxPeers: 3,
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynasty := NewDynasty([]string{}, tt.maxPeers, DefaultTimeBetweenBlockIfNoneGiven)
			dynasty.AddProducer(tt.input)
			assert.Equal(t, tt.expected, dynasty.producers)
		})
	}
}

func TestDynasty_AddMultipleProducers(t *testing.T) {
	tests := []struct {
		name     string
		maxPeers int
		input    []string
		expected []string
	}{
		{
			name:     "ValidInput",
			maxPeers: 3,
			input: []string{
				"dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
				"dG6HhzSdA5m7KqvJNszVSf8i5f4neAteSs",
				"dZ8GsrkSAiARL7ZnJLZSADzVXH4ea9EzhL"},
			expected: []string{
				"dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
				"dG6HhzSdA5m7KqvJNszVSf8i5f4neAteSs",
				"dZ8GsrkSAiARL7ZnJLZSADzVXH4ea9EzhL"},
		},
		{
			name:     "ExceedsLimit",
			maxPeers: 2,
			input: []string{
				"dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
				"dG6HhzSdA5m7KqvJNszVSf8i5f4neAteSs",
				"dZ8GsrkSAiARL7ZnJLZSADzVXH4ea9EzhL"},
			expected: []string{
				"dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
				"dG6HhzSdA5m7KqvJNszVSf8i5f4neAteSs",
			},
		},
		{
			name:     "InvalidInput",
			maxPeers: 3,
			input:    []string{"m1", "m2", "m3"},
			expected: []string{},
		},
		{
			name:     "mixedInput",
			maxPeers: 3,
			input:    []string{"m1", "dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf", "m3"},
			expected: []string{"dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf"},
		},
		{
			name:     "EmptyInput",
			maxPeers: 3,
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynasty := NewDynasty([]string{}, tt.maxPeers, DefaultTimeBetweenBlockIfNoneGiven)
			dynasty.AddMultipleProducers(tt.input)
			assert.Equal(t, tt.expected, dynasty.producers)
		})
	}
}

func TestDynasty_GetMinerIndex(t *testing.T) {
	tests := []struct {
		name             string
		initialProducers []string
		miner            string
		expected         int
	}{
		{
			name: "minerCouldBeFound",
			initialProducers: []string{
				"dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
				"dG6HhzSdA5m7KqvJNszVSf8i5f4neAteSs",
				"dZ8GsrkSAiARL7ZnJLZSADzVXH4ea9EzhL"},
			miner:    "dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
			expected: 0,
		},
		{
			name: "minerCouldNotBeFound",
			initialProducers: []string{
				"dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
				"dG6HhzSdA5m7KqvJNszVSf8i5f4neAteSs",
				"dZ8GsrkSAiARL7ZnJLZSADzVXH4ea9EzhL"},
			miner:    "dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivg",
			expected: -1,
		},
		{
			name: "EmptyInput",
			initialProducers: []string{
				"dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
				"dG6HhzSdA5m7KqvJNszVSf8i5f4neAteSs",
				"dZ8GsrkSAiARL7ZnJLZSADzVXH4ea9EzhL"},
			miner:    "",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynasty := NewDynasty(tt.initialProducers, len(tt.initialProducers), DefaultTimeBetweenBlockIfNoneGiven)
			index := dynasty.GetProducerIndex(tt.miner)
			assert.Equal(t, tt.expected, index)
		})
	}
}

func TestDynasty_IsMyTurnByIndex(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		now      int64
		expected bool
	}{
		{
			name:     "isMyTurn",
			index:    2,
			now:      105,
			expected: true,
		},
		{
			name:     "NotMyTurn",
			index:    1,
			now:      61,
			expected: false,
		},
		{
			name:     "InvalidIndexInput",
			index:    -6,
			now:      61,
			expected: false,
		},
		{
			name:     "InvalidNowInput",
			index:    2,
			now:      -1,
			expected: false,
		},
		{
			name:     "IndexInputExceedsMaxSize",
			index:    5,
			now:      44,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynasty := NewDynasty([]string{}, DefaultMaxProducersIfNoProducersGiven, DefaultTimeBetweenBlockIfNoneGiven)
			nextMintTime := dynasty.isMyTurnByIndex(tt.index, tt.now)
			assert.Equal(t, tt.expected, nextMintTime)
		})
	}
}

func TestDynasty_IsMyTurn(t *testing.T) {
	tests := []struct {
		name             string
		initialProducers []string
		producer         string
		index            int
		now              int64
		expected         bool
	}{
		{
			name: "IsMyTurn",
			initialProducers: []string{
				"121yKAXeG4cw6uaGCBYjWk9yTWmMkhcoDD",
				"1MeSBgufmzwpiJNLemUe1emxAussBnz7a7",
				"1LCn8D5W7DLV1CbKE3buuJgNJjSeoBw2ct"},
			producer: "1LCn8D5W7DLV1CbKE3buuJgNJjSeoBw2ct",
			now:      75,
			expected: true,
		},
		{
			name: "NotMyTurn",
			initialProducers: []string{
				"121yKAXeG4cw6uaGCBYjWk9yTWmMkhcoDD",
				"1MeSBgufmzwpiJNLemUe1emxAussBnz7a7",
				"1LCn8D5W7DLV1CbKE3buuJgNJjSeoBw2ct"},
			producer: "1MeSBgufmzwpiJNLemUe1emxAussBnz7a7",
			now:      61,
			expected: false,
		},
		{
			name: "EmptyInput",
			initialProducers: []string{
				"121yKAXeG4cw6uaGCBYjWk9yTWmMkhcoDD",
				"1MeSBgufmzwpiJNLemUe1emxAussBnz7a7",
				"1LCn8D5W7DLV1CbKE3buuJgNJjSeoBw2ct"},
			producer: "",
			now:      61,
			expected: false,
		},
		{
			name: "InvalidNowInput",
			initialProducers: []string{
				"121yKAXeG4cw6uaGCBYjWk9yTWmMkhcoDD",
				"1MeSBgufmzwpiJNLemUe1emxAussBnz7a7",
				"1LCn8D5W7DLV1CbKE3buuJgNJjSeoBw2ct"},
			producer: "m2",
			now:      0,
			expected: false,
		},
		{
			name: "minerNotFoundInDynasty",
			initialProducers: []string{
				"121yKAXeG4cw6uaGCBYjWk9yTWmMkhcoDD",
				"1MeSBgufmzwpiJNLemUe1emxAussBnz7a7",
				"1LCn8D5W7DLV1CbKE3buuJgNJjSeoBw2ct"},
			producer: "1LCn8D5W7DLV1CbKE3buuJgNJjSeoBw2cf",
			now:      90,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynasty := NewDynasty(tt.initialProducers, len(tt.initialProducers), DefaultTimeBetweenBlockIfNoneGiven)
			nextMintTime := dynasty.IsMyTurn(tt.producer, tt.now)
			assert.Equal(t, tt.expected, nextMintTime)
		})
	}
}

func TestDynasty_ProducerAtATime(t *testing.T) {
	producers := []string{
		"121yKAXeG4cw6uaGCBYjWk9yTWmMkhcoDD",
		"1MeSBgufmzwpiJNLemUe1emxAussBnz7a7",
		"1LCn8D5W7DLV1CbKE3buuJgNJjSeoBw2ct"}

	tests := []struct {
		name     string
		now      int64
		expected string
	}{
		{
			name:     "Normal",
			now:      62,
			expected: "1MeSBgufmzwpiJNLemUe1emxAussBnz7a7",
		},
		{
			name:     "InvalidInput",
			now:      -1,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynasty := NewDynasty(producers, len(producers), DefaultTimeBetweenBlockIfNoneGiven)
			producer := dynasty.ProducerAtATime(tt.now)
			assert.Equal(t, tt.expected, producer)
		})
	}
}

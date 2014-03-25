package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var TEST_BLOCK = []byte("The quick brown fox jumps over the lazy dog.")
var TEST_HASH = "e4d909c290d0fb1ca068ffaddf22cbd0"
var BAD_BLOCK = []byte("The magic words are squeamish ossifrage.")

// Test simple block reads.
func TestGetBlockOK(t *testing.T) {
	defer teardown()

	// Create two test Keep volumes and store a block in each of them.
	if err := setup(2); err != nil {
		t.Fatal(err)
	}
	for _, vol := range KeepVolumes {
		if err := storeTestBlock(vol, TEST_HASH, TEST_BLOCK); err != nil {
			t.Fatal(err)
		}
	}

	// Check that GetBlock returns success.
	result, err := GetBlock(TEST_HASH)
	if err != nil {
		t.Errorf("GetBlock error: %s", err)
	}
	if fmt.Sprint(result) != fmt.Sprint(TEST_BLOCK) {
		t.Errorf("expected %s, got %s", TEST_BLOCK, result)
	}
}

// Test block reads when one Keep volume is missing.
func TestGetBlockOneKeepOK(t *testing.T) {
	defer teardown()

	// Two test Keep volumes, only the second has a block.
	if err := setup(2); err != nil {
		t.Fatal(err)
	}
	if err := storeTestBlock(KeepVolumes[1], TEST_HASH, TEST_BLOCK); err != nil {
		t.Fatal(err)
	}

	// Check that GetBlock returns success.
	result, err := GetBlock(TEST_HASH)
	if err != nil {
		t.Errorf("GetBlock error: %s", err)
	}
	if fmt.Sprint(result) != fmt.Sprint(TEST_BLOCK) {
		t.Errorf("expected %s, got %s", TEST_BLOCK, result)
	}
}

// Test block read failure.
func TestGetBlockFail(t *testing.T) {
	defer teardown()

	// Create two empty test Keep volumes.
	if err := setup(2); err != nil {
		t.Fatal(err)
	}

	// Check that GetBlock returns failure.
	result, err := GetBlock(TEST_HASH)
	if err == nil {
		t.Errorf("GetBlock incorrectly returned success: ", result)
	}
}

// Test reading a corrupt block.
func TestGetBlockCorrupt(t *testing.T) {
	defer teardown()

	// Create two test Keep volumes and store a block in each of them,
	// but the hash of the block does not match the filename.
	if err := setup(2); err != nil {
		t.Fatal(err)
	}
	for _, vol := range KeepVolumes {
		if err := storeTestBlock(vol, TEST_HASH, BAD_BLOCK); err != nil {
			t.Fatal(err)
		}
	}

	// Check that GetBlock returns failure.
	result, err := GetBlock(TEST_HASH)
	if err == nil {
		t.Errorf("GetBlock incorrectly returned success: %s", result)
	}
}

// setup
//     Create KeepVolumes for testing.
func setup(nkeeps int) error {
	KeepVolumes = make([]string, 2)
	for i := range KeepVolumes {
		if dir, err := ioutil.TempDir(os.TempDir(), "keeptest"); err == nil {
			KeepVolumes[i] = dir + "/keep"
		} else {
			return err
		}
	}
	return nil
}

func teardown() {
	for _, vol := range KeepVolumes {
		os.RemoveAll(path.Dir(vol))
	}
}

func storeTestBlock(keepdir string, filename string, block []byte) error {
	blockdir := fmt.Sprintf("%s/%s", keepdir, filename[:3])
	if err := os.MkdirAll(blockdir, 0755); err != nil {
		return err
	}

	blockpath := fmt.Sprintf("%s/%s", blockdir, filename)
	if f, err := os.Create(blockpath); err == nil {
		f.Write(block)
		f.Close()
	} else {
		return err
	}

	return nil
}

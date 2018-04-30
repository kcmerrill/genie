package genie

import "testing"

func TestDir(t *testing.T) {
	if dir("/tmp") != "/tmp/" {
		t.Fatalf("dir() should add trailing slash")
	}

	if dir("/tmp/") != "/tmp/" {
		t.Fatalf("dir() should add trailing slash")
	}

}

func TestCmd(t *testing.T) {
	// test php
	if cmd("file.php") != "php" {
		t.Fatal(".php should return php")
	}

	// test python
	if cmd("file.py") != "python" {
		t.Fatal(".py should return python")
	}

	// test ruby
	if cmd("file.rb") != "ruby" {
		t.Fatal(".rb should return ruby")
	}

	// test sh
	if cmd("file.sh") != "sh" {
		t.Fatal(".sh should return sh")
	}

	// test default
	if cmd("python") != "python" {
		t.Fatal("Unknown file types should return the hint given")
	}
}

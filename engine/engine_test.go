package engine

import "testing"

func TestCRUD(t *testing.T) { // Create, Read, Update, Destroy
	nabia_db := NewNabiaDB()

	var nabia_read NabiaRecord
	var err error
	var expected []byte

	if nabia_db.Exists("A") {
		t.Error("Uninitialised database contains elements!")
	}
	//CREATE
	s := NewNabiaString("Value_A")
	nabia_db.Write("A", *s)
	if !nabia_db.Exists("A") {
		t.Error("Database is not writing items correctly!")
	}
	//READ
	nabia_read, err = nabia_db.Read("A") // TODO not testing content type
	if err != nil {
		t.Errorf("\"Read\" returns an unexpected error:\n%q", err.Error())
	}
	expected = []byte("Value_A")
	for i, e := range nabia_read.rawData {
		if e != expected[i] {
			t.Errorf("\"Read\" returns unexpected data!\nGot %q, expected %q", nabia_read, expected)
		}
	}
	//UPDATE
	s1 := NewNabiaString("Modified value") // TODO not testing another content type
	nabia_db.Write("A", *s1)
	if !nabia_db.Exists("A") {
		t.Errorf("Overwritten item doesn't exist!")
	}
	nabia_read, err = nabia_db.Read("A") // TODO not testing content-type
	if err != nil {
		t.Errorf("\"Read\" returns an unexpected error:\n%q", err.Error())
	}
	expected = []byte("Modified value")
	for i, e := range nabia_read.rawData {
		if e != expected[i] {
			t.Errorf("\"Write\" on an existing item saves unexpected data!\nGot %q, expected %q", nabia_read, expected)
		}
	}
	//DESTROY
	if !nabia_db.Exists("A") {
		t.Error("Can't destroy item because it doesn't exist!")
	}
	nabia_db.Destroy("A")
	if nabia_db.Exists("A") {
		t.Error("\"Destroy\" isn't working!\nDeleted item still exists in DB.")
	}

}

// TODO: Test concurrency

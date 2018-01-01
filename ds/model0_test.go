package ds

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	dbDir := "/tmp/TestModel0Db"
	mgr := ConnectDb(dbDir)
	col := MODEL0
	mgr.CreateCol(col)
	result := m.Run()
	mgr.DropCol(col)
	mgr.CloseDb()
	os.Exit(result)
}

func TestColumns(t *testing.T) {
	m0 := UseCol(MODEL0)
    if err := m0.Index([]string{"code"}); err != nil {
		t.Fatal("Missing column", err)
    }
    if err := m0.Index([]string{"name"}); err != nil {
		t.Fatal("Missing column", err)
    }
}

func TestFind(t *testing.T) {
	code := "100123"

	if found, err := FindModel0(code); err != nil {
		t.Fatal("Find failed:", err)
	} else if found.Code != "" {
		t.Fatal("Expected nil, found:", found)
	}

	// TODO insert test data and assert find
}

func TestInsert(t *testing.T) {
	m := Model0{
		Code: "100123",
		Name: "Entry1",
	}
	t.Log("entry:", m)
	if _, err := InsertModel0(m); err != nil {
		t.Fatal("Insert failed:", err)
	}
	if _, err := InsertModel0(m); err == nil {
		t.Fatal("Duplicate insert succeeded:", err)
	}
	// TODO fetch record and assert data
}

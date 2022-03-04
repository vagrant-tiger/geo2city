package tool

import (
	"testing"
)

func TestCreateTable(t *testing.T) {
	sc := NewMysqlConf("127.0.0.1", "3306", "user", "pass", "china_geo")
	err := sc.CreateTable()
	if err != nil {
		t.Error(err.Error())
	}
}

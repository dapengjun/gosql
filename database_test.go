package gosql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func Test_data(t *testing.T) {
	db := GetInstance()
	db.Init("mysql", "root", "root", "127.0.0.1", 3306, "dbname")
	list := db.Query("select * from tb")
	fmt.Println(len(list), list)
}

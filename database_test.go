package gosql

import (
	"fmt"
	"testing"
)

func Test_data(t *testing.T) {
	db := GetInstance()
	db.Init("mysql", "fengshen", "", "", 3306, "voip_conf")
	list := db.Query("select * from ac_conference")
	fmt.Println(list)
}

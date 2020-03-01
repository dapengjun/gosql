package gosql

import (
	"database/sql"
	"fmt"
	log "github.com/fsjun/logrotating"
	"sync"
)

type DataBase struct {
	db     *sql.DB
	dbType string
	dsn    string
}

var d *DataBase
var once sync.Once

func GetInstance() *DataBase {
	once.Do(func() {
		d = &DataBase{}
	})
	return d
}

func (database *DataBase) Init(dbType string, user string, password string, host string, port int, dbname string) {
	database.dbType = dbType
	database.dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbname)
}

func (database *DataBase) connect() bool {
	if database.db != nil {
		return true
	}
	db, err := sql.Open(database.dbType, database.dsn)
	if err != nil {
		log.Error(err)
		return false
	}
	database.db = db
	return true
}

func (database *DataBase) Insert(sql string, args ...interface{}) int64 {
	if !database.connect() {
		return 0
	}
	stmt, err := database.db.Prepare(sql)
	if err != nil {
		log.Errorln(err)
		return 0
	}
	defer stmt.Close()

	ret, err := stmt.Exec(args...)
	if err != nil {
		log.Errorln(err)
		return 0
	}
	lastInsertId, err := ret.LastInsertId()
	if nil != err {
		log.Errorln(err)
		return 0
	}
	return lastInsertId
}

func (database *DataBase) Exec(sql string, args ...interface{}) int64 {
	if !database.connect() {
		return 0
	}
	stmt, err := database.db.Prepare(sql)
	if err != nil {
		log.Errorln(err)
		return 0
	}
	defer stmt.Close()

	ret, err := stmt.Exec(args...)
	if err != nil {
		log.Errorln(err)
		return 0
	}
	rowsAffected, err := ret.RowsAffected()
	if nil != err {
		log.Errorln(err)
		return 0
	}
	return rowsAffected
}

func (database *DataBase) QueryOne(sql string, args ...interface{}) map[string]string {
	if !database.connect() {
		return nil
	}
	result := make(map[string]string, 0)
	stmt, err := database.db.Prepare(sql)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	defer rows.Close()
	if err != nil {
		log.Errorln(err)
		return nil
	}
	columns, err := rows.Columns()
	if err != nil {
		log.Errorln(err)
		return nil
	}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range columns {
			values[i] = new(*string)
		}
		err := rows.Scan(values...)
		if err != nil {
			log.Errorln(err)
			return nil
		}
		for k, v := range columns {
			val := *(values[k].(**string))
			if val == nil {
				result[v] = ""
			} else {
				result[v] = *val
			}
		}
		break
	}

	err = rows.Err()
	if err != nil {
		log.Errorln(err)
		return nil
	}
	return result
}

func (database *DataBase) Query(sql string, args ...interface{}) []map[string]string {
	if !database.connect() {
		return nil
	}
	result := make([]map[string]string, 0)
	stmt, err := database.db.Prepare(sql)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	defer rows.Close()
	if err != nil {
		log.Errorln(err)
		return nil
	}
	columns, err := rows.Columns()
	if err != nil {
		log.Errorln(err)
		return nil
	}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range columns {
			values[i] = new(*string)
		}
		err := rows.Scan(values...)
		if err != nil {
			log.Errorln(err)
			return nil
		}
		row := map[string]string{}
		for k, v := range columns {
			val := *(values[k].(**string))
			if val == nil {
				row[v] = ""
			} else {
				row[v] = *val
			}
		}
		result = append(result, row)
	}

	err = rows.Err()
	if err != nil {
		log.Errorln(err)
		return nil
	}
	return result
}

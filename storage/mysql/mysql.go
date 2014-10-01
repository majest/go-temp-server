package mysql

import (
	"database/sql"
	"fmt"
	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
)

type Db struct {
	connection *sql.DB
	table      string
}

func New(user, pass, host, port, dbname, table string) *Db {

	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname))

	log.Debugf("Connecting to db: %s:%s/%s", host, port, dbname)
	if err != nil {
		log.Errorf(err.Error())
	}

	db := &Db{}
	db.connection = conn
	db.table = table
	return db
}

func (d *Db) InsertData(data string) {

	var t = time.Now()
	var data2, data3, data4 int64
	var waga float64
	var err error

	dataArr := strings.Split(data, ";")

	if len(dataArr) != 5 {
		log.Debugf("Bad data size. Got: %s, expected 5 fields", len(dataArr))
		return
	}

	log.Debugf("Inserting data: %s", data)
	// 1.0;2;3;4;123.123.123.123
	insert := fmt.Sprintf("INSERT INTO %s(waga,data,wej1,wej2,wej3,ip) VALUES(?,?,?,?,?,?)", d.table)

	data2, err = strconv.ParseInt(dataArr[1], 10, 8)
	checkErr(err)

	data3, err = strconv.ParseInt(dataArr[2], 10, 8)
	checkErr(err)

	data4, err = strconv.ParseInt(dataArr[3], 10, 8)
	checkErr(err)

	waga, err = strconv.ParseFloat(dataArr[0], 32)
	checkErr(err)

	_, err = d.connection.Exec(insert, waga, t.Format(time.RFC3339), data2, data3, data4, dataArr[4])
	checkErr(err)

}

func checkErr(err error) {
	if err != nil {
		log.Errorf(err.Error())
	}
}

package creditd

import (
	"database/sql"
	"encoding/json"
	//"errors"
	"fmt"
	"gsxt/credit"
	"strconv"
	"strings"
	//"gsxt/gsxt/creditd"
	"time"
	//	"test/mongoDB"

	"crypto/md5"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var Statistic map[string]*Count = make(map[string]*Count)

type Count struct {
	Name  string
	Count int64
	Time  int64
}

func addstati(cname string, ctime int64) {
	_, ok := Statistic[cname]
	if ok {
		Statistic[cname].Count = Statistic[cname].Count + 1
		Statistic[cname].Time = Statistic[cname].Time + ctime
	} else {
		Statistic[cname] = &Count{Name: cname, Count: 1, Time: ctime}
	}
}

func init() {
	//user:password@tcp(localhost:5555)/dbname?charset=utf8
	//  db, _ = sql.Open("mysql", "zhengguoqiang:zgq@2017@tcp(180.76.168.121:3306)/guoqiang?charset=utf8")
	//	db, _ = sql.Open("mysql", "zhengguoqiang:zgq@2017@tcp(192.168.0.15:3306)/dc_import_append?charset=utf8")
	db, _ = sql.Open("mysql", "crawler:ark#2017@tcp(180.76.172.99:3306)/ent?charset=utf8")
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(115)
	db.Ping()

}

//更新数据
func update() {
	stmt, err := db.Prepare("UPDATE user SET user_age=?,user_sex=? WHERE user_id=?")
	checkErr(err)
	res, err := stmt.Exec(21, 2, 1)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(num)
}

//删除数据
func remove() {
	stmt, err := db.Prepare("DELETE FROM user WHERE user_id=?")
	checkErr(err)
	res, err := stmt.Exec(1)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(num)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func InsertBusinessInfo(msg []credit.Mbus) {
	tx, err := db.Begin()
	checkErr(err)
	sql := "REPLACE INTO business_info(province,create_time,name,md5,type,regno,base,investors,changes,members,branchs,licenses,mortgages,pledges,punishs,abnormals,spot_checks) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := tx.Prepare(sql)
	defer stmt.Close()
	checkErr(err)
	if len(msg) > 0 {
		for i := 0; i < len(msg); i++ {
			t := msg[i]

			_, err := stmt.Exec(t.Province, checkTime(t.CreateDate), t.Base.Name, t.Id, t.Base.Type, t.Base.RegNo, Obj2String(t.Base), Obj2String(t.Investors), Obj2String(t.Changes), Obj2String(t.Members), Obj2String(t.Branchs), Obj2String(t.Licenses), Obj2String(t.Mortgages), Obj2String(t.Pledges), Obj2String(t.Punishs), Obj2String(t.Abnormals), Obj2String(t.SpotChecks))
			checkErr(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		checkErr(err)
		//		err3 := tx.Rollback()
		//		checkErr(err3)
	}
}

func InsertEnterpriseInfo(msg []credit.Ment) {
	tx, err := db.Begin()
	checkErr(err)
	sql := "REPLACE INTO enterprise_info(province,create_time,md5,investors,changes,stock_changes,licenses,intells,punishs) VALUES(?,?,?,?,?,?,?,?,?)"
	sql1 := "REPLACE INTO report_info(province,create_time,year,md5,date,`from`,general,operation,websites,licenses,branchs,invents,guarantees,investors,stockchanges,changes) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := tx.Prepare(sql)
	stmt1, err1 := tx.Prepare(sql1)
	defer stmt.Close()
	defer stmt1.Close()
	checkErr(err)
	checkErr(err1)
	if len(msg) > 0 {
		for i := 0; i < len(msg); i++ {
			t := msg[i]
			_, err := stmt.Exec(t.Province, checkTime(t.CreateDate), t.Id, Obj2String(t.Investors), Obj2String(t.Changes), Obj2String(t.StockChanges), Obj2String(t.Licenses), Obj2String(t.Intells), Obj2String(t.Punishs))
			checkErr(err)

			for j := 0; j < len(t.Reports); j++ {
				k := t.Reports[j]
				_, err := stmt1.Exec(t.Province, checkTime(t.CreateDate), k.Year, t.Id, k.Date, k.From, Obj2String(k.General), Obj2String(k.Operation), Obj2String(k.Websites), Obj2String(k.Licenses), Obj2String(k.Branchs), Obj2String(k.InvEnts), Obj2String(k.Guarantees),
					Obj2String(k.Investors), Obj2String(k.StockChanges), Obj2String(k.Changes))
				checkErr(err)
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		checkErr(err)
		//		err3 := tx.Rollback()
		//		checkErr(err3)
	}
}

func InsertMsgInfoByV2(msg []credit.InfoV2, tprovince string) {
	var t1, t2 int64

	create_time := time.Now().Format("2006-01-02")
	Province := m1[tprovince]
	tx, err := db.Begin()
	checkErr(err)
	sql1 := "REPLACE INTO business_info(province,create_time,name,md5,type,regno,base,investors,changes,members,branchs,licenses,mortgages,pledges,punishs,abnormals,spot_checks) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	sql2 := "REPLACE INTO enterprise_info(province,create_time,md5,investors,changes,stock_changes,licenses,intells,punishs) VALUES(?,?,?,?,?,?,?,?,?)"
	sql3 := "REPLACE INTO report_info(province,create_time,year,md5,date,`from`,general,operation,websites,licenses,branchs,invents,guarantees,investors,stockchanges,changes) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt1, err1 := tx.Prepare(sql1)
	stmt2, err2 := tx.Prepare(sql2)
	stmt3, err3 := tx.Prepare(sql3)
	defer stmt1.Close()
	defer stmt2.Close()
	defer stmt3.Close()
	checkErr(err1)
	checkErr(err2)
	checkErr(err3)
	if len(msg) > 0 {
		for i := 0; i < len(msg); i++ {
			t := msg[i]
			//md5 Id 计算
			if len(t.Business.Base.RegNo) == 15 {
				data := []byte(t.Business.Base.RegNo)
				has := md5.Sum(data)
				md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

				t1 = time.Now().UnixNano()
				_, err := stmt1.Exec(Province, checkTime(create_time), t.Business.Base.Name, md5str1, t.Business.Base.Type, t.Business.Base.RegNo, Obj2String(t.Business.Base), Obj2String(t.Business.Investors), Obj2String(t.Business.Changes), Obj2String(t.Business.Members), Obj2String(t.Business.Branchs), Obj2String(t.Business.Licenses), Obj2String(t.Business.Mortgages), Obj2String(t.Business.Pledges), Obj2String(t.Business.Punishs), Obj2String(t.Business.Abnormals), Obj2String(t.Business.SpotChecks))
				checkErr(err)
				t2 = time.Now().UnixNano()
				addstati("sql1", t2-t1)

				t1 = time.Now().UnixNano()
				_, err = stmt2.Exec(Province, checkTime(create_time), md5str1, Obj2String(t.Enterprise.Investors), Obj2String(t.Enterprise.Changes), Obj2String(t.Enterprise.StockChanges), Obj2String(t.Enterprise.Licenses), Obj2String(t.Enterprise.Intells), Obj2String(t.Enterprise.Punishs))
				checkErr(err)
				t2 = time.Now().UnixNano()
				addstati("sql2", t2-t1)

				for j := 0; j < len(t.Enterprise.Reports); j++ {
					k := t.Enterprise.Reports[j]

					t1 = time.Now().UnixNano()
					_, err := stmt3.Exec(Province, checkTime(create_time), k.Year, md5str1, k.Date, k.From, Obj2String(k.General), Obj2String(k.Operation), Obj2String(k.Websites), Obj2String(k.Licenses), Obj2String(k.Branchs), Obj2String(k.InvEnts), Obj2String(k.Guarantees), Obj2String(k.Investors), Obj2String(k.StockChanges), Obj2String(k.Changes))
					checkErr(err)
					t2 = time.Now().UnixNano()
					addstati("sql3", t2-t1)

				}
			} else {
				fmt.Println("注册号不存在！")
				continue
			}
		}
	}
	t1 = time.Now().UnixNano()
	err = tx.Commit()
	t2 = time.Now().UnixNano()
	addstati("commit", t2-t1)
	if err != nil {
		checkErr(err)
		//		err3 := tx.Rollback()
		//		checkErr(err3)
	}
}

func Obj2String(msg interface{}) interface{} {
	temp, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	stemp := string(temp)
	if stemp == "[]" || stemp == "{}" {
		return nil
	} else {
		return stemp
	}
}

func checkTime(time string) interface{} {
	if time == "" {
		return nil
	} else {
		return time
	}

}

//查找关键词
func QueryKeywordMsg() (msg []string) {
	sql := "SELECT keyword,province,crawler_count,id FROM `keyword_info` WHERE province >0 and `status`<2 LIMIT 1000;"
	t1 := time.Now().UnixNano()
	rows, err := db.Query(sql)
	t2 := time.Now().UnixNano()
	addstati("queryKeywordMsg", t2-t1)
	checkErr(err)
	//创建字典
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]string, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	//查询
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		checkErr(err)
		fmt.Println("QueryKeywordMsg:", values)
		msg = append(msg, strings.Join(values, "+"))
	}
	return
}

//更新关键词的状态
func UpdateKeywordStatus(id, status, search_count string) {
	count, err := strconv.Atoi(search_count)
	checkErr(err)
	stmt, err := db.Prepare("UPDATE keyword_info SET `status`= ?,crawler_count =? WHERE id = ?")
	checkErr(err)
	res, err := stmt.Exec(status, count+1, id)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(num)
}

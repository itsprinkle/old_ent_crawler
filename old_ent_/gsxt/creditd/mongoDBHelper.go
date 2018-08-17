package creditd

import (
	//"bufio"
	//"encoding/json"
	//"gsxt/gsxt/creditd"
	//"io"
	//"os"
	//"strings"

	"crypto/md5"
	"time"
	//      "strconv"

	"errors"
	"fmt"

	"gsxt/credit"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type msg struct {
	Id         string         `bson:"_id"`
	Province   int            `bson:"province"`
	RegNo      string         `bson:"regno"`
	CreateDate string         `bson:"create_date"`
	Content    *credit.InfoV2 `bson:"content"`
}

type ent struct {
	Id           string                   `bson:"_id"`
	Province     int                      `bson:"province"`
	CreateDate   string                   `bson:"create_date"`
	Reports      []credit.ReportInfo      `bson:"reports"`
	Investors    []credit.InvestorInfo    `bson:"investors"`
	Changes      []credit.ChangeInfo      `bson:"changes"`
	StockChanges []credit.StockChangeInfo `bson:"stock_changes"`
	Licenses     []credit.LicenseInfo     `bson:"licenses"`
	Intells      []credit.IntellInfo      `bson:"intells"`
	Punishs      []credit.PunishInfo      `bson:"punishs"`
}

type bus struct {
	Id         string                 `bson:"_id"`
	Province   int                    `bson:"province"`
	CreateDate string                 `bson:"create_date"`
	Base       credit.BaseInfo        `bson:"base"`
	Investors  []credit.InvestorInfo  `bson:"investors"`
	Changes    []credit.ChangeInfo    `bson:"changes"`
	Members    []credit.MemberInfo    `bson:"members"`
	Branchs    []credit.BranchInfo    `bson:"branchs"`
	Licenses   []credit.LicenseInfo   `bson:"licenses"`
	Mortgages  []credit.MortgageInfo  `bson:"mortgages"`
	Pledges    []credit.PledgeInfo    `bson:"pledges"`
	Punishs    []credit.PunishInfo    `bson:"punishs"`
	Abnormals  []credit.AbnormalInfo  `bson:"abnormals"`
	SpotChecks []credit.SpotCheckInfo `bson:"spot_checks"`
}

const (
	URL = "mongodb://crawler:ark#2017@180.76.190.74:20001/gsxt"
)

var (
	mgoSession *mgo.Session
	dataBase   = "gsxt"
	m1         map[string]int
)

func init() {
	m1 = make(map[string]int)
	m1["anhui"] = 1
	m1["beijing"] = 2
	m1["fujian"] = 3
	m1["gansu"] = 4
	m1["guangxi"] = 5
	m1["hainan"] = 6
	m1["hebei"] = 7
	m1["heilongjiang"] = 8
	m1["henan"] = 9
	m1["hubei"] = 10
	m1["hunan"] = 11
	m1["jiangsu"] = 12
	m1["jilin"] = 13
	m1["liaoning"] = 14
	m1["ningxia"] = 15
	m1["qinghai"] = 16
	m1["shandong"] = 17
	m1["shanghai"] = 18
	m1["shanxi"] = 19
	m1["tianjin"] = 20
	m1["xinjiang"] = 21
	m1["xizang"] = 22
	m1["yunnan"] = 23
	m1["zongju"] = 24
	m1["guangdong"] = 25
	m1["chongqing"] = 26
	m1["zhejiang"] = 27
	m1["sichuan"] = 28
	m1["guizhou"] = 29
	m1["neimenggu"] = 30
	m1["xianxi"] = 31
	m1["jiangxi"] = 32
}

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法，获取collection对象
func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

func getNextID(collection string) int {
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	doc := struct{ Seq int }{}
	query := func(c *mgo.Collection) error {
		_, err := c.Find(bson.M{"collection": collection}).Apply(change, &doc)
		return err
	}
	witchCollection("provice_seq", query)
	return doc.Seq
}

func addEnt(id string, p *ent, collection string) error {
	query := func(c *mgo.Collection) error {
		_, err := c.UpsertId(id, p)
		return err
	}
	return witchCollection(collection, query)
}

func addBus(id string, p *bus, collection string) error {
	query := func(c *mgo.Collection) error {
		_, err := c.UpsertId(id, p)
		return err
	}
	return witchCollection(collection, query)
}

func AddEnt(stb *credit.InfoV2, collection, province string) error {
	if len(stb.Business.Base.RegNo) == 15 {
		data := []byte(stb.Business.Base.RegNo)
		has := md5.Sum(data)
		md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

		//      Id           string                   `bson:"_id"`
		//      Province     int                      `bson:"province"`
		//      Reports      []credit.ReportInfo      `bson:"reports"`
		//      Investors    []credit.InvestorInfo    `bson:"investors"`
		//      Changes      []credit.ChangeInfo      `bson:"changes"`
		//      StockChanges []credit.StockChangeInfo `bson:"stock_changes"`
		//      Licenses     []credit.LicenseInfo     `bson:"licenses"`
		//      Intells      []credit.IntellInfo      `bson:"intells"`
		//      Punishs      []credit.PunishInfo      `bson:"punishs"`

		var temp ent
		temp.CreateDate = time.Now().Format("2006-01-02")
		temp.Province = m1[province]
		temp.Reports = stb.Enterprise.Reports
		temp.Investors = stb.Enterprise.Investors
		temp.Changes = stb.Enterprise.Changes
		temp.StockChanges = stb.Enterprise.StockChanges
		temp.Licenses = stb.Enterprise.Licenses
		temp.Intells = stb.Enterprise.Intells
		temp.Punishs = stb.Enterprise.Punishs

		return addEnt(md5str1, &temp, collection)
	} else {
		return errors.New("注册号不存在！")
	}
}

func AddBus(stb *credit.InfoV2, collection, province string) error {
	if len(stb.Business.Base.RegNo) == 15 {
		data := []byte(stb.Business.Base.RegNo)
		has := md5.Sum(data)
		md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

		//      Id         string          `bson:"_id"`
		//      Province   int             `bson:"province"`
		//      Base       BaseInfo        `bson:"base"`
		//      Investors  []InvestorInfo  `bson:"investors"`
		//      Changes    []ChangeInfo    `bson:"changes"`
		//      Members    []MemberInfo    `bson:"members"`
		//      Branchs    []BranchInfo    `bson:"branchs"`
		//      Licenses   []LicenseInfo   `bson:"licenses"`
		//      Mortgages  []MortgageInfo  `bson:"mortgages"`
		//      Pledges    []PledgeInfo    `bson:"pledges"`
		//      Punishs    []PunishInfo    `bson:"punishs"`
		//      Abnormals  []AbnormalInfo  `bson:"abnormals"`
		//      SpotChecks []SpotCheckInfo `bson:"spot_checks"`

		var temp bus
		temp.CreateDate = time.Now().Format("2006-01-02")
		temp.Province = m1[province]
		temp.Base = stb.Business.Base
		temp.Investors = stb.Business.Investors
		temp.Changes = stb.Business.Changes
		temp.Members = stb.Business.Members
		temp.Branchs = stb.Business.Branchs
		temp.Licenses = stb.Business.Licenses
		temp.Mortgages = stb.Business.Mortgages
		temp.Pledges = stb.Business.Pledges
		temp.Punishs = stb.Business.Punishs
		temp.Abnormals = stb.Business.Abnormals
		temp.SpotChecks = stb.Business.SpotChecks

		return addBus(md5str1, &temp, collection)
	} else {
		return errors.New("注册号不存在！")
	}
}

func AddMsg(stb *credit.InfoV2, province string) error {
	var err error
	err = AddBus(stb, "business_"+province, province)
	if err != nil {
		return errors.New("插入business_失败")
	}
	err = AddEnt(stb, "enterprise_"+province, province)
	if err != nil {
		return errors.New("插入enterprise_失败")
	}
	return err
}

func GetBusiness(collection string, i int) (rows []credit.Mbus) {
	var row credit.Mbus
	query := func(c *mgo.Collection) error {
		iter := c.Find(bson.M{"create_date": bson.M{"$exists": true}}).Skip(i * 10).Limit(10).Iter() //500000
		defer iter.Close()
		for iter.Next(&row) {
			fmt.Println("row>>>>", row)
			rows = append(rows, row)
		}
		return nil
	}
	witchCollection(collection, query)
	return
}

func GetEnterprise(collection string, i int) (rows []credit.Ment) {
	var row credit.Ment
	query := func(c *mgo.Collection) error {
		iter := c.Find(bson.M{"create_date": bson.M{"$exists": true}}).Skip(i * 1000000).Limit(1000000).Iter() //500000
		defer iter.Close()
		for iter.Next(&row) {
			fmt.Println("row>>>>", row)
			rows = append(rows, row)
		}
		return nil
	}
	witchCollection(collection, query)
	return
}

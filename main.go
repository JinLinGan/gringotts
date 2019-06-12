package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//type Product struct {
//	gorm.Model
//	Code  string
//	Price uint
//}
//
//func main1() {
//	db, err := gorm.Open("sqlite3", "test.db")
//	if err != nil {
//		panic("failed to connect database")
//	}
//	defer db.Close()
//
//	// Migrate the schema
//	db.AutoMigrate(&Product{})
//
//	// 创建
//	db.Create(&Product{Code: "L1212", Price: 1000})
//
//	// 读取
//	var product Product
//	db.First(&product, 1)                   // 查询id为1的product
//	db.First(&product, "code = ?", "L1212") // 查询code为l1212的product
//
//	// 更新 - 更新product的price为2000
//	db.Model(&product).Update("Price", 2000)
//
//	// 删除 - 删除product
//	db.Delete(&product)
//}

// Host .
type Host struct {
	gorm.Model
	HostName      string
	HostInterface []*HostInterface
}

// HostInterface .
type HostInterface struct {
	gorm.Model
	HostID        uint
	HWAddr        string
	InterfaceName string
	IPAddress     string
}

func main() {
	db, err := gorm.Open("mysql", "gringotts:gringotts@tcp(127.0.0.1)/gringotts?parseTime=true")
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()
	db = db.Debug()
	db.SingularTable(true)
	db.AutoMigrate(&Host{}, &HostInterface{})
	//db.AutoMigrate(&HostInterface{})
	h := &Host{
		HostName: "testhost",
		HostInterface: []*HostInterface{
			{
				HWAddr:        "aaa",
				InterfaceName: "bbb",
			},
		},
	}
	//db.Create(h)
	//fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	//db.Model(h).Association("HostInterface").Append(&HostInterface{
	//	HWAddr:        "111",
	//	InterfaceName: "111",
	//})
	//fmt.Println("1111111")
	////db.Save(h)
	//fmt.Println("bbbbbbbbbbbbbbbbbbb")
	//db.Model(h).Association("HostInterface").Append(&HostInterface{
	//	HWAddr:        "222",
	//	InterfaceName: "222",
	//})
	//fmt.Println("222222")
	//db.Model(h).Association("HostInterface").Append(&HostInterface{
	//	HWAddr:        "333",
	//	InterfaceName: "3333",
	//})
	//fmt.Println("333333")
	////db.Save(h)
	////fmt.Println("ccccccccccccc")
	h.HostInterface = append(h.HostInterface, &HostInterface{
		HWAddr:        "1111",
		InterfaceName: "1111",
	})

	h.HostInterface = append(h.HostInterface, &HostInterface{
		HWAddr:        "222",
		InterfaceName: "222",
	})
	//fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	//spew.Dump(h)
	db.Save(h)
	////spew.Dump(h)
	//
	//db.Save(h)
}

func main2() {
	db, err := gorm.Open("mysql", "gringotts:gringotts@tcp(127.0.0.1)/gringotts?parseTime=true")
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()
	db = db.Debug()
	db.SingularTable(true)
	db.AutoMigrate(&Host{}, &HostInterface{})
	results := []Host{}
	db.Preload("HostInterface").Table("host").Joins("left join host_interface on host_interface.host_id = host.id").Where("host_interface.hw_addr = ?", "1111").Find(&results)
	//if err != nil {
	//	log.Println(err)
	//}
	//for rows.Next() {
	//	var result Host
	//	if err := db.ScanRows(rows, &result); err != nil {
	//		log.Println(err)
	//	}
	//	db.Model(&result).Association("host_interface")
	//	results = append(results, result)
	//}
	spew.Dump(results)
}

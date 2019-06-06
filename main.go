package main

import (
	"fmt"
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
}

func main() {
	db, err := gorm.Open("mysql", "gringotts:gringotts@tcp(127.0.0.1)/gringotts")
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
	db.Model(h).Association("HostInterface").Append(&HostInterface{
		HWAddr:        "111",
		InterfaceName: "111",
	})
	fmt.Println("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	spew.Dump(h)
	db.Save(h)
	//spew.Dump(h)
	h.HostInterface = append(h.HostInterface, &HostInterface{
		HWAddr:        "222",
		InterfaceName: "222",
	})
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	spew.Dump(h)
	db.Save(h)
	//spew.Dump(h)
}

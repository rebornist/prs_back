package configs

import (
	"encoding/json"
	"fmt"

	"prs/customTypes"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb() *gorm.DB {

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	var DB customTypes.Database
	getInfo, err := GetServiceInfo("database")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(getInfo, &DB)

	// gorm DB 접속
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul", DB.Host, DB.User, DB.Password, DB.Name, DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	return db
}

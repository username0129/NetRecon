package test

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

type User struct {
	gorm.Model
	Name  string
	Age   int
	Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	gorm.Model
	Name string
}

func initDB() *gorm.DB {

	db, err := gorm.Open(mysql.Open("root:root@tcp(192.168.80.130:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect db")
	}

	// 自动迁移模式
	db.AutoMigrate(&User{}, &Role{})

	return db
}

func createUser(db *gorm.DB, name string, age int) {
	user := User{
		Name: name,
		Age:  age,
		Roles: []Role{
			{Name: "Admin"},
			{Name: "User"},
		},
	}
	db.Create(&user)
}

func queryUsers(db *gorm.DB, userId uint) {
	var user User
	db.Preload("Roles").First(&user, userId)
	fmt.Printf("User: %v\n", user.Name)
	for _, role := range user.Roles {
		fmt.Printf("Role: %v\n", role.Name)
	}
}

func TestGorm(t *testing.T) {
	db := initDB()

	// 插入示例数据
	//createUser(db, "John Doe", 30)

	// 查询并打印用户数据
	queryUsers(db, 9)
}

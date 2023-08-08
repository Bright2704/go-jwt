package orm

import (
	"gorm.io/gorm" //สร้างตาราง
)



type User struct {
	gorm.Model  //ตัวโมเดลสำเร็จรูป มี id, create_at, update_at, deleted_at
	Username	string
	Password	string
	Fullname	string
	Avatar		string
}
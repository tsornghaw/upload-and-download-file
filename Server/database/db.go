package database

import (
	"fmt"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormDatabase struct {
	dsn    string
	dbtype string
	db     *gorm.DB
}

type DatabaseInterface interface {
	CombineSqlConnectionStrings(username string, password string, addr string, port int, dbname string)
	InitDatabase() error
	CreateTable() error
	GetAllData(data interface{}, conds ...interface{}) (interface{}, error)
	GetCorrespondingData(data interface{}, conds ...interface{}) (interface{}, error)
	CreateData(data interface{}) error
	UpdateData(data interface{}) error
	DeleteData(data interface{}, id interface{}) error
}

/**
 * ------------------------------------------------------------------------------------------------------------------------
 * 0. Create type User struct {}
 * 1. Combine sql connection strings
 *		- CombineSqlConnectionStrings(...)
 * 2. Connect to the database
 *		- InitDatabase()
 * 3. Create table - AutoMigrate will create the table if it doesn't exist
 *		- CreateTable()
 * 4. Insert model information
 *		- InsertData(users interface{})
 * 5. Action
 *		5.1 Update data
 *				- UpdateData
 *		5.2 Get data
 *				- GetUsernameByID
 * ------------------------------------------------------------------------------------------------------------------------
 */

func NewGormDatabase(username string, password string, addr string, port int, dbname string, dbtype string) (*GormDatabase, error) {
	gd := &GormDatabase{
		dsn:    "",
		dbtype: dbtype,
		db:     nil,
	}

	gd.CombineSqlConnectionStrings(username, password, addr, port, dbname)

	err := gd.InitDatabase(dbtype)
	if err != nil {
		return nil, err
	}

	return gd, err
}

func (gd *GormDatabase) InitDatabase(dbtype string) error {
	var err error

	switch dbtype {
	case "mysql", "TiDB":
		gd.db, err = gorm.Open(mysql.Open(gd.dsn), &gorm.Config{})
	case "postgres":
		gd.db, err = gorm.Open(postgres.Open(gd.dsn), &gorm.Config{})
	case "sqlserver":
		gd.db, err = gorm.Open(sqlserver.Open(gd.dsn), &gorm.Config{})
	case "sqlite":
		gd.db, err = gorm.Open(sqlite.Open(gd.dsn), &gorm.Config{})
	case "clickhouse":
		gd.db, err = gorm.Open(clickhouse.Open(gd.dsn), &gorm.Config{})
	}

	return err
}

func (gd *GormDatabase) Close() error {
	if gd.db != nil {

		sqlDB, err := gd.db.DB()

		if err != nil {
			return err
		}

		return sqlDB.Close()
	}

	return nil
}

func (gd *GormDatabase) CreateTable(tables ...interface{}) error {
	// AutoMigrate will create tables if they don't exist
	return gd.db.Debug().AutoMigrate(tables...)
}

/**
 * Assert the returned interface{} as a User struct to access the retrieved data.
 * 	- If the assertion is successful, we can use the User struct to access the fields.
 *	- Otherwise, we handle the error or invalid data type as needed.
 *		user, ok := data.(User)
 *		if !ok {
 *			fmt.Println("Error: Invalid data type")
 *			return
 *		}
 */

func (gd *GormDatabase) GetAll(data interface{}, conds ...interface{}) error {
	// Data here should be like this: []user
	return gd.db.Find(data, conds).Error
}

func (gd *GormDatabase) GetCorresponding(data interface{}, query interface{}, args ...interface{}) error {
	// Important: make sure the value we pass in is an pointer. (&user)
	// Retrieving objects with primary key
	// return gd.db.First(data, conds).Error
	return gd.db.Where(query, args).Find(data).Error
}

func (gd *GormDatabase) Create(data interface{}) error {
	// Important: make sure the value we pass in is an pointer. (&user)
	return gd.db.Create(data).Error
}

func (gd *GormDatabase) Update(data interface{}) error {
	// Important: make sure the value we pass in is an pointer. (&user)
	return gd.db.Save(data).Error
}

func (gd *GormDatabase) Upsert(data interface{}) error {
	// Important: make sure the value we pass in is an pointer. (&user)
	return gd.db.Debug().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(data).Error
}

func (gd *GormDatabase) Delete(data interface{}, conds ...interface{}) error {
	// Important: make sure the value we pass in is an pointer.
	// Data here should be like this: &user{}
	return gd.db.Delete(data, conds).Error
}

func (gd *GormDatabase) CombineSqlConnectionStrings(username string, password string, addr string, port int, dbname string) {
	switch gd.dbtype {
	case "mysql":
		// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		gd.dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", username, password, addr, port, dbname)
	case "postgres":
		// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		gd.dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", addr, username, password, dbname, port)
	case "sqlserver":
		// dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
		gd.dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", username, password, addr, port, dbname)
	case "TiDB":
		// Not for sure : dsn := "root:@tcp(127.0.0.1:4000)/test?charset=utf8mb4"
		gd.dsn = fmt.Sprintf("%s:@tcp(%s:%d)/%s?charset=utf8mb4", username, addr, port, dbname)
	case "sqlite":
		// Not for sure : An SQLite database is normally stored in a single ordinary disk file.
		gd.dsn = "gorm.db"
	case "clickhouse":
		// dsn := "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
		gd.dsn = fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s&read_timeout=10&write_timeout=20", addr, port, dbname, username, password)
	}
}

// func Testing(gd *GormDatabase) {
// 	var user1 []User
// 	var user2 User
// 	var user3 User

// 	gd.GetAll(&user1)
// 	fmt.Printf("GetAllData: %v\n", user1)

// 	gd.GetCorresponding(&user2, "ID = ?", "2")

// 	user3 = User{ID: 3, Name: "Tester3", Email: "tester3emailaddr@email.com", Password: "12345678", RoomId: 3}
// 	gd.Update(&user3)

// 	gd.Delete(&user1, "ID = ?", "3")
// }

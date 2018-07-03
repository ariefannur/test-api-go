package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

type UserReview struct {
	Id        int64  `db:"id" json:"id" value:"AUTO_INCREMENT"`
	OrderId   int64  `db:"order_id" json:"order_id"`
	ProductId int64  `db:"product_id" json:"product_id"`
	UserId    int64  `db:"user_id" json:"user_id"`
	Rating    int8   `db:"rating" json:"rating"`
	Review    string `db:"review" json:"review"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
}

var dbmap = initDb()

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:@/testdb")
	checkErr(err, "sql.Open failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(UserReview{}, "user_review").SetKeys(true, "id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}

}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(Cors())

	v1 := router.Group("api/v1")
	{
		v1.POST("/add", AddReview)
		v1.GET("/", GetAllReview)
		v1.DELETE("/delete/:id", DeleteReview)
		v1.PUT("/update/:id", UpdateReview)
	}

	return router
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	router := SetupRouter()
	router.Run(":8080")

}

func GetAllReview(c *gin.Context) {

	var userReview []UserReview
	_, err := dbmap.Select(&userReview, "SELECT * FROM user_review")

	if err == nil {
		c.JSON(200, userReview)
	} else {
		c.JSON(404, gin.H{"code": 404, "message": "empty data"})
	}

}

func AddReview(c *gin.Context) {

	var user UserReview
	c.BindJSON(&user)

	if user.OrderId != 0 && user.ProductId != 0 && user.UserId != 0 {

		tm := time.Now().Unix()

		if insert, err := dbmap.Exec(`INSERT INTO user_review (order_id, product_id, user_id, rating, review, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, user.OrderId, user.ProductId, user.UserId, user.Rating, user.Review, tm, tm); insert != nil {

			c.JSON(200, gin.H{"code": 200, "message": "success review"})
		} else {
			c.JSON(200, gin.H{"code": 404, "message": "failed review " + err.Error()})
		}
	}

}

func DeleteReview(c *gin.Context) {

	id := c.Params.ByName("id")

	d, err := dbmap.Exec("delete from user_review where id=?", id)
	fmt.Println(d)
	if err == nil {
		c.JSON(200, gin.H{"id #" + id: "deleted"})
	}
}

func UpdateReview(c *gin.Context) {

	id := c.Params.ByName("id")

	var user UserReview
	c.BindJSON(&user)
	tm := time.Now().Unix()

	if d, err := dbmap.Exec(`UPDATE user_review SET order_id = ?, product_id=? , user_id = ? , rating=?, review=?, updated_at=? WHERE id = ? `, user.OrderId, user.ProductId, user.UserId, user.Rating, user.Review, tm, id); d != nil {

		c.JSON(200, gin.H{"code": 200, "message": "success review"})
	} else {
		c.JSON(200, gin.H{"code": 404, "message": "failed update review " + err.Error()})
	}

}

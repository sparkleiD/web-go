package main

import (
	"database/sql"
	"path/filepath"

	//"encoding/json"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
)

type User struct {
	username string
	password string
	sex      string
	email    string
	iconurl  string
	level    string
}

func errhandle(e any, option ...any) {
	if e != nil {
		log.Fatal(e)
		if len(option) != 0 {
			for _, el := range option {
				log.Fatal(el)
			}
		}
	}
}

func tokenInit(key any, method jwt.SigningMethod, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	return token.SignedString(key)
}

func main() {
	router := gin.Default() //新建路由对象

	router.LoadHTMLGlob("templates/*.html")
	router.Static("/icon", "static/icon")
	router.Static("/image", "static/image")
	router.Static("/js", "static/js")
	router.Static("/css", "static/css") //路径映射
	//Static(站点路径,实际路径)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/infocheck", func(c *gin.Context) {
		c.HTML(http.StatusOK, "infocheck.html", gin.H{})
	})

	router.POST("/ajax/userinfo", func(c *gin.Context) {
		email := c.PostForm("email")
		db, sqlerr1 := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/aremember")
		db.Ping()
		defer db.Close()
		if sqlerr1 != nil {
			fmt.Println("DBconn fail!")
			log.Fatal(sqlerr1)
		} else {
			fmt.Println("DBconn success!")
		}
		sqlsent, sqlerr2 := db.Prepare("SELECT email FROM users WHERE email = ?")
		if sqlerr2 != nil {
			log.Fatal(sqlerr2)
		}

		sqlcheckresp, sqlerr3 := sqlsent.Query(email)
		if sqlerr3 != nil {
			log.Fatal(sqlerr3)
		}
		var user User
		for sqlcheckresp.Next() {
			sqlerr4 := sqlcheckresp.Scan(&user.email)
			if sqlerr4 != nil {
				log.Fatal(sqlerr4)
			}
			log.Println(user.email)
		}
		if user.email == "" {
			c.String(http.StatusOK, "true")
		} else {
			c.String(http.StatusOK, "false")
		}
	})

	router.POST("/signup", func(c *gin.Context) {
		email := c.PostForm("email")
		username := c.PostForm("username")
		password := c.PostForm("password")
		sex := c.PostForm("sex")

		icon, err := c.FormFile("icon")
		if err != nil {
			log.Fatal("图片获取失败：", err)
		}
		iconext := filepath.Ext(icon.Filename)
		fileName := email + iconext
		filePath := "icon/" + fileName
		if err := c.SaveUploadedFile(icon, "static/"+filePath); err != nil {
			log.Fatal("保存文件失败：", err)
		}

		fmt.Println("email:", email)
		fmt.Println("password:", password)

		db, sqlerr1 := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/aremember")
		db.Ping()
		defer db.Close()
		if sqlerr1 != nil {
			fmt.Println("DBconn fail!")
			log.Fatal(sqlerr1)
		} else {
			fmt.Println("DBconn success!")
		}
		sqlsent, sqlerr2 := db.Prepare("INSERT INTO users VALUES (? , ? , ? , ? , ? ,'δ')")
		if sqlerr2 != nil {
			log.Fatal(sqlerr2)
		}

		_, sqlerr3 := sqlsent.Exec(username, password, sex, email, filePath)
		if sqlerr3 != nil {
			log.Fatal(sqlerr3)
			c.HTML(http.StatusFailedDependency, "result.html", gin.H{
				"result": "数据插入错误！",
			})
		} else {
			c.HTML(http.StatusOK, "result.html", gin.H{
				"result": "您的数据已登记，欢迎加入ARE！",
			})
		}
	})

	router.POST("/login", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		fmt.Println("email:", email)
		fmt.Println("password:", password)

		db, sqlerr1 := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/aremember")
		db.Ping()
		defer db.Close()
		if sqlerr1 != nil {
			fmt.Println("DBconn fail!")
			log.Fatal(sqlerr1)
		} else {
			fmt.Println("DBconn success!")
		}
		sqlsent, sqlerr2 := db.Prepare("SELECT * FROM users WHERE email = ?")
		if sqlerr2 != nil {
			log.Fatal(sqlerr2)
		}

		sqlcheckresp, sqlerr3 := sqlsent.Query(email)
		if sqlerr3 != nil {
			log.Fatal(sqlerr3)
		}
		var user User
		for sqlcheckresp.Next() {
			sqlerr4 := sqlcheckresp.Scan(&user.username, &user.password, &user.sex, &user.email, &user.iconurl, &user.level)
			if sqlerr4 != nil {
				log.Fatal(sqlerr4)
			}
			log.Println(user.username, user.password, user.sex, user.email, user.iconurl, user.level)
		}

		if user.username != "" {
			if user.password == password {
				jwtKey := make([]byte, 32) // 生成32字节（256位）的密钥
				_, err := rand.Read(jwtKey)
				errhandle(err, "密钥生成错误")
				jwtStr, err := tokenInit(jwtKey, jwt.SigningMethodHS256, jwt.MapClaims{
					"iss": "arefd.com",
					"aud": email,
					"exp": time.Now().Add(time.Hour * 720).UnixMilli(),
					"iat": time.Now().UnixMilli(),
				})
				if err != nil {
					log.Fatal(err)
					log.Fatal("token生成失败")
				} else {
					c.Header("Token", jwtStr)
				}

				c.HTML(http.StatusOK, "result.html", gin.H{
					"result": "欢迎回来！",
				})
			} else {
				c.HTML(http.StatusOK, "result.html", gin.H{
					"result": "密码错误！",
				})
			}
		} else {
			c.HTML(http.StatusOK, "result.html", gin.H{
				"result": "请先注册账号！",
			})
		}
	})

	router.Run(":9090") //开始服务
}

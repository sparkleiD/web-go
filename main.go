package main

import (
	"database/sql"
	"path/filepath"

	//"encoding/json"
	//"crypto/rand"
	"errors"
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
		if len(option) != 0 {
			for _, el := range option {
				log.Println(el)
			}
		}
		panic(e)
	}
}

func ParseJwtWithClaims(key any, jwtStr string) (jwt.MapClaims, error) {
	mc := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtStr, mc, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	// 校验 Claims 对象是否有效，基于 exp（过期时间），nbf（不早于），iat（签发时间）等进行判断（如果有这些声明的话）。
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return mc, nil
}

func main() {
	//定义jwt密钥(暂时)
	jwtKey := []byte{147, 177, 39, 92, 144, 145, 21, 252, 239, 187, 17, 39, 46, 207, 26, 112, 131, 66, 9, 141, 112, 83, 239, 187, 166, 237, 7, 245, 35, 176, 174, 210}

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
		Token := c.GetHeader("Authorization")
		log.Println(Token)
		//此处逻辑表达式非为最简，但因为一些问题只能暂时这么写
		if Token != "" && Token != "null" {
			claims, err := ParseJwtWithClaims(jwtKey, Token)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(claims["aud"])
			}
			log.Println("准备返回 result.html")
			c.HTML(http.StatusOK, "result.html", gin.H{
				"result": "欢迎回来!",
			})
		} else {
			log.Println("进入 else 分支，返回 infocheck.html")
			c.HTML(http.StatusOK, "infocheck.html", gin.H{})
		}
	})

	router.POST("/ajax/userinfo", func(c *gin.Context) {
		email := c.PostForm("email")
		db, sqlerr1 := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/aremember")
		db.Ping()
		defer db.Close()
		if sqlerr1 != nil {
			fmt.Println("DBconn fail!")
			panic(sqlerr1)
		} else {
			fmt.Println("DBconn success!")
		}
		sqlsent, sqlerr2 := db.Prepare("SELECT email FROM users WHERE email = ?")
		if sqlerr2 != nil {
			panic(sqlerr2)
		}

		sqlcheckresp, sqlerr3 := sqlsent.Query(email)
		if sqlerr3 != nil {
			panic(sqlerr3)
		}
		var user User
		for sqlcheckresp.Next() {
			sqlerr4 := sqlcheckresp.Scan(&user.email)
			if sqlerr4 != nil {
				panic(sqlerr4)
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
			panic(err)
		}
		iconext := filepath.Ext(icon.Filename)
		fileName := email + iconext
		filePath := "icon/" + fileName
		if err := c.SaveUploadedFile(icon, "static/"+filePath); err != nil {
			panic(err)
		}

		fmt.Println("email:", email)
		fmt.Println("password:", password)

		db, sqlerr1 := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/aremember")
		db.Ping()
		defer db.Close()
		if sqlerr1 != nil {
			fmt.Println("DBconn fail!")
			panic(sqlerr1)
		} else {
			fmt.Println("DBconn success!")
		}
		sqlsent, sqlerr2 := db.Prepare("INSERT INTO users VALUES (? , ? , ? , ? , ? ,'δ')")
		if sqlerr2 != nil {
			panic(sqlerr2)
		}

		_, sqlerr3 := sqlsent.Exec(username, password, sex, email, filePath)
		if sqlerr3 != nil {
			log.Println(sqlerr3)
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
			panic(sqlerr1)
		} else {
			fmt.Println("DBconn success!")
		}
		sqlsent, sqlerr2 := db.Prepare("SELECT * FROM users WHERE email = ?")
		if sqlerr2 != nil {
			panic(sqlerr2)
		}

		sqlcheckresp, sqlerr3 := sqlsent.Query(email)
		errhandle(sqlerr3)
		var user User
		for sqlcheckresp.Next() {
			sqlerr4 := sqlcheckresp.Scan(&user.username, &user.password, &user.sex, &user.email, &user.iconurl, &user.level)
			if sqlerr4 != nil {
				panic(sqlerr4)
			}
			log.Println(user.username, user.password, user.sex, user.email, user.iconurl, user.level)
		}

		if user.username != "" {
			if user.password == password {
				/*jwtKey := make([]byte, 32)
				// 生成32字节（256位）的密钥
				_, err := rand.Read(jwtKey)
				errhandle(err, "密钥生成错误")*/
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"iss": "arefd.com",
					"aud": user.username,
					"exp": time.Now().Add(time.Hour * 720).UnixMilli(),
					//"iat": time.Now().UnixMilli(),
				})
				jwtStr, err := token.SignedString(jwtKey)
				if err != nil {
					log.Println("token生成失败:", err)
				} else {
					c.Header("Authorization", jwtStr)
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

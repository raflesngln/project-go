package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	_ "strconv"
	"github.com/gin-gonic/gin"
)

type Users struct {
	NamaDepan    string `json:"nama_depan"`
	NamaBelakang string `json:"nama_belakang"`
}

func encodeUSER(data []Users) []byte {
	dataJSON, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}
	return dataJSON
}

func decodeUSER(data []byte) []Users {
	var usr []Users
	err := json.Unmarshal(data, &usr)
	if err != nil {
		log.Fatal(err)
	}
	return usr
}

func main() {
	dataStruct()
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Selamat datang Selamat belajar GO")
	})
	router.Static("/assets", "./assets")   //memanggil file assets agar bisa dipangggil oleh funsi lain
	router.GET("/user/", User)             // http://localhost:8000/user
	router.GET("/user/:name/", UserDetail) // http://localhost:8000/user/rafles/?address=jakarta%20barat

	/* ROUTER GROUP */
	v1 := router.Group("/v1")
	{
		v1.GET("/user", UserGroup) //http://localhost:8000/v1/user
	}
	v2 := router.Group("/v2")
	{
		v2.GET("/user", NewUser) //http://localhost:8000/v2/user
	}
	web := router.Group("/web")
	{
		router.LoadHTMLGlob("web/*")
		web.GET("/", HomePage) //http://localhost:8000/v2/user
		web.GET("/profile", Profile)
		web.GET("/about", About)
		web.GET("/encode_json", encodeJSON)
		web.GET("/decode_json", decodeJSON)
	}

	router.Run(":8000")
}

func User(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello User Routessss",
	})
}
func UserDetail(c *gin.Context) {
	// Dengan PATH
	name := c.Param("name")
	// Dengan Query String
	address := c.Query("address")
	c.JSON(200, gin.H{
		"message": name,
		"Alamat":  address,
	})
	// c.String(http.StatusOK, "Hello  %s Alamat %s", name, address)
}
func UserGroup(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("Dari V1"))
}
func NewUser(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("Dari V2"))
}
func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"judul": "Response dengan Output HTML",
	})
}
func Profile(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{
		"judul": "Profile Paging",
	})
}
func About(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"judul": "About Paging",
	})
}

func encodeJSON(c *gin.Context) {
	usr := []Users{{
		NamaDepan:    "Mawar",
		NamaBelakang: "Merah",
	},
		{
			NamaDepan:    "Budi",
			NamaBelakang: "Angga",
		},
	}
	dataJSON := encodeUSER(usr)

	c.HTML(http.StatusOK, "encode_json.html", gin.H{
		"judul": "Data JSON Encode",
		"data":  string(dataJSON),
	})
}
func decodeJSON(c *gin.Context) {
	usrJSON := `[{"nama_depan":"Rafles","nama_belakang":"Nainggolan"},
	{"nama_depan":"Mawar","nama_belakang":"Merah"}]`

	dataUsr := decodeUSER([]byte(usrJSON))
	listdata := []string{}
	for key, value := range dataUsr {
		// fmt.Printf("%d : Nama Lengkap %v %v\n", i, v.NamaDepan, v.NamaBelakang)
		listdata = append(listdata, fmt.Sprintf("%d, %s, %s, %s", key, value.NamaDepan, value.NamaBelakang, "\n"))
	}

	c.HTML(http.StatusOK, "encode_json.html", gin.H{
		"judul": "Data JSON Decode",
		"list":  dataUsr,
		"data":  strings.Join(listdata, ";"),
	})
}


func aritMatika(a int, b int, status string) (string, string) {
		var hasil int
		if status == "jumlah" {
			hasil=a+b
		} else if status =="bagi" {
			hasil=a/b
		} else {
			hasil=a-b
		}
	return " Hasil aritmatika adalah "+strconv.Itoa(hasil),
		" Jenis aritmatika adalah "+status
}


type person struct {
	name string
	age  int
}
var allStudents = []struct {
	person
	grade int
}{
	{person: person{"wick", 28}, grade: 7},
	{person: person{"ethan", 25}, grade: 5},
	{person: person{"bond", 21}, grade: 8},
}

func dataStruct(){
	for _, student := range allStudents {
		fmt.Println(student.name, "age is", student.age)
	}
}


type student struct {
	person struct {
		name string
		age  int
	}
	grade   int
	hobbies []string
}
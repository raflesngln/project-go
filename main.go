package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	//"path/filepath"
	"strconv"
	_ "strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"
)
var client *redis.Client

func init() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
//
//
//type Account struct {
//	ID uint64            `json:"id"`
//	Username string `json:"username"`
//	Password string `json:"password"`
//}
////A sample user account get from DB
//var accounts = Account{
//	ID:             1,
//	Username: "raflesngln",
//	Password: "12345",
//}
//func Login(c *gin.Context) {
//	var u Account
//	if err := c.ShouldBindJSON(&u); err != nil {
//		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
//		return
//	}
//	//compare the user from the request, with the one we defined:
//	if accounts.Username != u.Username || accounts.Password != u.Password {
//		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
//		return
//	}
//	token, err := CreateToken(accounts.ID,accounts.Username)
//	if err != nil {
//		c.JSON(http.StatusUnprocessableEntity, err.Error())
//		return
//	}
//	c.JSON(http.StatusOK, token)
//}
//func CreateToken(userid uint64,username string) (string, error) {
//	var err error
//	//Creating Access Token
//	os.Setenv("ACCESS_SECRET", "nainggolan") //this should be in an env file
//	atClaims := jwt.MapClaims{}
//	atClaims["authorized"] = true
//	atClaims["user_id"] = userid
//	atClaims["user_nm"] = username
//	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
//	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
//	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
//	if err != nil {
//		return "", err
//	}
//	return token, nil
//}

/* ==================================================== */

type Users struct {
	NamaDepan    string `json:"nama_depan"`
	NamaBelakang string `json:"nama_belakang"`
}

// Middleware Login
func MyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.Request.URL.Query()
		var person = val
		token := c.Request.FormValue("api_token")

		//fmt.Printf("%v\n", person["username"])
		fmt.Printf("%v\n", person)
		fmt.Println("Im a Middleware!" + token)
		c.Next()
		//fmt.Println("ini ADALAH MIDDLEWARE",val)
	}

}
func MyMiddlewarePerRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Middleware Using in Per Routes")
		c.Next()
	}

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
	router := gin.Default()
	router.Static("/assets", "./assets") //memanggil file assets agar bisa dipangggil oleh funsi lain
	router.Use(MyMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Selamat datang Selamat belajar GO")
	})
	/*  FOR TOKEN MANAGE  */
	router.POST("/login", Login)
	router.POST("/todo", CreateTodo)
	router.POST("/logout", Logout)
	router.POST("/token/refresh", Refresh)

	router.GET("/form_upload", Form_upload)
	router.POST("/upload_file", UploadFile)


	router.GET("/product/", ProductPage) // http://localhost:8000/product

	router.GET("/user/:name/", UserDetail) // http://localhost:8000/user/rafles/?address=jakarta%20barat
	router.GET("/user/", User)             // http://localhost:8000/user
	router.POST("user", MyMiddlewarePerRoute(), insertUser)        // http://localhost:8000/user/rafles/?address=jakarta%20barat
	router.POST("form_post", MyMiddlewarePerRoute(), form_post)        // multipart form data


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

	router.StaticFS("/file", http.Dir("public"))
	router.Run(":8000")
}

func User(c *gin.Context) {
	val := c.Request.URL.Query()
	name := c.Param("name")
	password := c.Param("password")
	c.JSON(200, gin.H{
		"datas": val,
		"name": name,
		"password": password,
	})
}

type USERSDATA struct {
	USER     string `json:"user" binding:"required"`
	PASSWORD string `json:"password" binding:"required"`
}

func insertUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	//c.String(http.StatusOK, "Hello %s", name)

	c.JSON(200, gin.H{
		"username": username,
		"password": password,
	})
}
func form_post(c *gin.Context) {
	message := c.PostForm("message")
	gambar, _ := c.FormFile("gambar")

	// Upload the file to specific dst.
	// c.SaveUploadedFile(file, dst)

	c.JSON(200, gin.H{
		"status":  "posted",
		"message": message,
		"gambar":gambar.Filename,
	})
}

func UserDetail(c *gin.Context) {
	// Dengan PATH
	name := c.Param("name")
	// Dengan Query String
	address := c.Query("address")
	c.JSON(200, gin.H{
		"Nama": name,
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

func ProductPage(c *gin.Context) {
	prod := []Product{{
		NamaDepan:    "Didik",
		NamaBelakang: "Prabowo",
	},
		{
			NamaDepan:    "Charly",
			NamaBelakang: "Van Houten",
		},
	}
	c.JSON(200, gin.H{
		"data": prod,
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

type Product struct {
	NamaDepan    string `json:"nama_depan"`
	NamaBelakang string `json:"nama_belakang"`
}

func encodeProduct(data []Product) []byte {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return dataJSON
}



/*======================== FOR GENERATE TOKEN ==================================*/
type Account struct {
	ID uint64            `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
//A sample user account get from DB
var accounts = Account{
	ID:             1,
	Username: "raflesngln",
	Password: "12345",
}


func Login(c *gin.Context) {
	var u Account
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	if accounts.Username != u.Username || accounts.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	ts, err := CreateToken(accounts.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := CreateAuth(accounts.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

type AccessDetails struct {
	AccessUuid string
	UserId   uint64
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    uint64
	RtExpires    uint64
}


func CreateToken(userid uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = uint64(time.Now().Add(time.Minute * 15).Unix())
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = uint64(time.Now().Add(time.Hour * 24 * 7).Unix())
	td.RefreshUuid = td.AccessUuid + "++" + strconv.Itoa(int(userid))


	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}


func CreateAuth(userid uint64, td *TokenDetails) error {
	at := time.Unix(int64(td.AtExpires), 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(int64(td.RtExpires), 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

type Todo struct {
	UserID uint64 `json:"user_id"`
	Title string `json:"title"`
}


func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		var userId, err = strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			accessUuid,
			userId,
		}, nil
	}
	return nil, err
}


func FetchAuth(authD *AccessDetails) (uint64, error) {
	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	if authD.UserId != userID {
		return 0, errors.New("unauthorized")
	}
	return userID, nil
}

func CreateTodo(c *gin.Context) {
	var td Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	//Extract the access token metadata
	metadata, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userid, err := FetchAuth(metadata)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	td.UserID = userid
	//you can proceed to save the Todo to a database
	//but we will just return it to the caller:

	c.JSON(http.StatusCreated, td)
}

func DeleteAuth(givenUuid string) (int64,error) {
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}


func Logout(c *gin.Context) {
	metadata, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	delErr := DeleteTokens(metadata)
	if delErr != nil {
		c.JSON(http.StatusUnauthorized, delErr.Error())
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}


func Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		fmt.Println("the error: ", err)
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := CreateToken(userId)
		if  createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		saveErr := CreateAuth(userId, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}

func DeleteTokens(authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.AccessUuid, authD.UserId)
	//delete access token
	deletedAt, err := client.Del(authD.AccessUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := client.Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	out, err := os.Create("public/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	filepath := "http://localhost:8000/public/" + filename
	c.JSON(http.StatusOK, gin.H{"filepath": filepath})
}

func Form_upload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload_file.html", gin.H{
		"judul": "Profile Paging",
	})
}
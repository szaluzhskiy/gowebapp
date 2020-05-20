package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	Username string `json: username`
	Password string `json: password`
}

var router *gin.Engine
var client *resty.Client = resty.New()
var messagesCh chan []byte = make(chan []byte, 1000)

var cadenceAddress = "http://127.0.0.1:8092/result"

func main() {
	router = gin.New()
	initializeRoutes()
	router.Run(":8090")
}

func initializeRoutes() {
	router.POST("/api", handleVerification)
	router.OPTIONS("/api", handleVerification)
	router.GET("/api", handleGet)
	router.PUT("/api/cadence/async", handleCadenceActivityAsyncCompletion)
}

func handleGet(c *gin.Context) {
	message, _ := c.GetQuery("m")
	c.String(http.StatusOK, "Get works! you sent: "+message)
}

func handleCadenceActivityAsyncCompletion(c *gin.Context) {
	body := c.Request.Body
	token, _ := ioutil.ReadAll(body)
	messagesCh <- token
	go processAsyncResponse(messagesCh)
	c.String(http.StatusOK, "Get async cadence activity completion request")
}

func handleVerification(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		// setup headers
		c.Header("Allow", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "origin, content-type, accept")
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	} else if c.Request.Method == "POST" {
		var u User
		c.BindJSON(&u)
		c.JSON(http.StatusOK, gin.H{
			"user": u.Username,
			"pass": u.Password,
		})
	}
}

func processAsyncResponse(messagesCh chan []byte) {
	token, _ := <-messagesCh
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(1000)
	time.Sleep(time.Duration(r) * time.Millisecond)
	resp, err := client.R().
		SetBody([]byte(token)).
		EnableTrace().
		Put(cadenceAddress)
	if resp.StatusCode() == 200 {
		s := string(token)
		fmt.Println("Success response from time %v reponse %v token %v", time.Now(), resp, s)
	} else {
		fmt.Println("Error response from remote %v", err)
	}
}

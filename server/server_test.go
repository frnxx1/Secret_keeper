package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"server/pkg/key"
	"server/pkg/storage"

	"github.com/gin-gonic/gin"
)

var keeper = storage.GetDummyKeeper()

func handleTestRequest(w *httptest.ResponseRecorder, r *http.Request) {
	keyBuilder := key.GetKeyBuilder()
	router := getRouter(keyBuilder, keeper)
	router.ServeHTTP(w, r)

}

func TestIndexPageCase(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestSaveMessage(t *testing.T) {
	testMessage := "foo"
	postData := strings.NewReader(fmt.Sprintf("message=%s", testMessage))
	request, _ := http.NewRequest("POST", "/", postData)
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 200 {
		t.Error("save is not 200")
	}

	keyBuilder := key.GetKeyBuilder()
	key, _ := keyBuilder.Get()
	savedMessage, _ := keeper.Get(key)
	if savedMessage != testMessage {
		t.Error("message was not saved")
	}

	result := w.Result()
	defer result.Body.Close()
	data, _ := ioutil.ReadAll(result.Body)

	if !strings.Contains(string(data), key) {
		t.Error("result page without key")
	}

}

func TestReadMessage(t *testing.T) {
	keyBuilder := key.GetKeyBuilder()
	key, _ := keyBuilder.Get()
	testMessage := "foo"
	keeper.Set(key, testMessage, 0)

	request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)

	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 200 {
		t.Error("read message is not 200")
	}

	result := w.Result()
	defer result.Body.Close()
	data, _ := ioutil.ReadAll(result.Body)

	if !strings.Contains(string(data), testMessage) {
		t.Error("result page without key")
	}

	_, err := keeper.Get(key)
	if err == nil {
		t.Error("Keeper value must be empty")
	}

}

func TestReadMessageNotFound(t *testing.T) {
	keyBuilder := key.GetKeyBuilder()
	key, _ := keyBuilder.Get()

	request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)

	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 200 {
		t.Error("empty message must be 404")
	}
}

func TestOneMessage(t *testing.T) {
	dummyKeeper := storage.GetDummyKeeper()
	keyBuilder := key.GetKeyBuilder()
	key, _ := keyBuilder.Get()
	testMessage := "foo"
	dummyKeeper.Set(key, testMessage, 0)

	router := getRouter(keyBuilder, dummyKeeper)
	resultChannal := make(chan int, 2)

	go func(key string, c chan int, router *gin.Engine) {
		request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, request)
		resultChannal <- w.Code
	}(key, resultChannal, router)

	go func(key string, c chan int, router *gin.Engine) {
		request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, request)
		resultChannal <- w.Code
	}(key, resultChannal, router)

	firstCode := <-resultChannal
	secondCode := <-resultChannal

	if firstCode+secondCode != (200 + 404) {
		t.Error("one answer must be 404")
	}

}

func TestCheckLenMessage(t *testing.T) {
	testMessage := ""

	for i := 0; i < 1025; i++ {
		testMessage += "1"
	}
	postData := strings.NewReader(fmt.Sprintf("message=%s", testMessage))
	request, _ := http.NewRequest("POST", "/", postData)
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 411 {
		t.Error("save is not 200")
	}

}

func TestCheckTTLMax(t *testing.T) {
	testMessage := "foo"
	ttl := 999999
	postData := strings.NewReader(fmt.Sprintf("message=%s&ttl=%d", testMessage, ttl))
	request, _ := http.NewRequest("POST", "/", postData)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 411 {
		t.Error("Must be error, because message is too long, but received", w.Code)
	}

}


func TestCheckTTLMin(t *testing.T) {
	testMessage := "foo"
	
	postData := strings.NewReader(fmt.Sprintf("message=%s&ttl=%d", testMessage, MIN_TTL-1))
	request, _ := http.NewRequest("POST", "/", postData)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 411 {
		t.Error("Must be error, because message is too long, but received", w.Code)
	}
	

}

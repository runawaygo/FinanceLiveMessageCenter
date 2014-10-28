package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LiveMessage struct {
	Content string
	Time    string
	Market  string
	Person  string
	Type    string
}

const durationTime = 200 * time.Second

var ticker = time.NewTicker(durationTime)
var quit = make(chan struct{})

func getLiveMessageCollection() []LiveMessage {
	now := time.Now().Add(-durationTime)
	url := "http://b.1stool.com/index.php?m=index&c=index&a=bcajax&type=all&btype=all&lasttime=" + strconv.FormatInt(now.Unix(), 10)
	fmt.Println(url)
	resp, _ := http.Get(url)
	content, _ := ioutil.ReadAll(resp.Body)
	contents := strings.Split(string(content), "~|")

	liveMessageCollection := []LiveMessage{}
	for i := 1; i < len(contents); i += 5 {
		lm := &LiveMessage{
			Content: contents[i],
			Time:    contents[i+1],
			Market:  contents[i+2],
			Person:  contents[i+3],
			Type:    contents[i+4],
		}
		liveMessageCollection = append(liveMessageCollection, *lm)
	}
	return liveMessageCollection
}

func loopMessage() {
	for {
		select {
		case <-ticker.C:
			jsonDatas, _ := json.Marshal(getLiveMessageCollection())
			h.broadcast <- jsonDatas
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

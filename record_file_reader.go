package main

import (
	"encoding/json"
	"log"
	"strings"
	"strconv"
	"os"
	//"io"
	"io/ioutil"
	"bufio"
	"fmt"
	"net/http"
	"bytes"
)
import camcloud "github.com/kzcabstone/camcloud"

func check(e error) {
    if e != nil {
    	log.Fatal(e)
    }
}

func record_file_reader_routine(filename string, result_channel chan string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	counter := 0
	for scanner.Scan() {
		counter = counter + 1
		line := scanner.Text()
		words := strings.Fields(line)

		if len(words) < 13 {
			log.Printf("Line error, expect 13 or 14 fields, got %d", len(words))
			continue
		}

		plate := ""
		if len(words) > 13 {
			plate = strings.Join(words[13:], ",")
		}
		x, _ := strconv.ParseInt(words[4], 10, 0)
		y, _ := strconv.ParseInt(words[5], 10, 0)
		w, _ := strconv.ParseInt(words[6], 10, 0)
		h, _ := strconv.ParseInt(words[7], 10, 0)
		id, _ := strconv.ParseInt(words[3], 10, 64)
		type_score, _ := strconv.ParseFloat(words[11], 64)
		obj_score, _ := strconv.ParseFloat(words[9], 64)
		
		record := camcloud.CamRecord{
			Ts: time.Now(),
			Cam: words[0],
			Software: "v1",
			Object_id: id,
			Object_x: int(x),
			Object_y: int(y),
			Object_w: int(w),
			Object_h: int(h),
			Object_label: words[8],
			Object_score: obj_score,
			Vehicle_type: words[10],
			Vehicle_type_score: type_score,
			Vehicle_color: words[12],
			Vehicle_plate: plate,
		}

		jsonstr, err := json.Marshal(record)
		check(err)

		doHttpPost(jsonstr)	
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %s", err)
		result_channel <- "Error: " + err.Error()
		return
	}

	result_channel <- "ok"
}

func doHttpPost(jsonstr []byte) {
	req, err := http.NewRequest("POST", "http://localhost:8073/u", bytes.NewBuffer(jsonstr))
    req.Header.Set("X-Custom-Header", "hhe")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    
    check(err)
    defer resp.Body.Close()

    if resp.Status != "200 OK" {
	    fmt.Println("response Status:", resp.Status)
	    fmt.Println("response Headers:", resp.Header)
	    body, _ := ioutil.ReadAll(resp.Body)
	    fmt.Println("response Body:", string(body))
	}
}
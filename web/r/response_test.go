package r_test

import (
	"encoding/json"
	"github.com/gogoclouds/gogo/internal/server/response"
	"github.com/gogoclouds/gogo/web/r"
	"log"
	"testing"
)

func TestResp(t *testing.T) {
	list := []string{"a", "b", "c"}
	meta := map[string]int{
		"a": 1,
	}
	page := r.NewPageMeta(list, 1, 2, 10, meta)
	resp := response.SuccessData(*page)
	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("json.Marshal", err)
	}
	log.Printf("%s", string(bytes))

	p := &response.RespData[r.PageMetaResp[string, map[string]int]]{}
	err = json.Unmarshal(bytes, p)
	if err != nil {
		log.Println("json.Unmarshal", err)
	}
	log.Printf("%+v", *p)
}

func TestPage(t *testing.T) {
	list := []string{"a", "b", "c"}
	meta := map[string]int{
		"a": 1,
	}
	page := r.NewPageMeta(list, 1, 2, 10, meta)
	bytes, err := json.Marshal(page)
	if err != nil {
		log.Println("json.Marshal", err)
	}
	log.Printf("%s", string(bytes))

	p := &r.PageMetaResp[string, map[string]int]{}
	err = json.Unmarshal(bytes, p)
	if err != nil {
		log.Println("json.Unmarshal", err)
	}
	log.Printf("%+v", *p)
}
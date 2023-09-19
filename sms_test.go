package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestSms(t *testing.T) {
	jsonStr := []byte("{\n  \"mob\": \"15046332824\",\n  \"msg\": \"尊敬的 李华 您已成功预订中国联航KN5737上海虹桥T2\\t-缅甸丹老机场，05月29日07:25起飞11:25到达,行李额0KG,请您提前120分钟到达机场,祝您一路平安！\"\n}")
	xx, _ := http.NewRequest("POST", "http://messagefat.chengdd.cn/Msg/SendSms", bytes.NewBuffer(jsonStr))
	xx.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(xx)
	fmt.Println(resp.StatusCode)
}

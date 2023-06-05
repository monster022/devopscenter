package model

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestDemo(t *testing.T) {
	//message := "admin:Mzl123456"
	//encoded := base64.StdEncoding.EncodeToString([]byte(message))
	//fmt.Println(encoded)
	wd, _ := os.Getwd()
	output := exec.Command("cat " + wd + "\\..\\config.ini")
	fmt.Println(output)
}

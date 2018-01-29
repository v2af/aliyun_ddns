package lib

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func GetIP() (ip string, err error) {
	cmd := exec.Command("curl", "ifconfig.me/ip")
	output, err := cmd.Output()
	if err != nil {
		log.Println("get ip failed :", err)
		return
	}

	reg := regexp.MustCompile(`(?m)^.*[\d]*\.[\d]*\.[\d]*\.[\d]*.*$`)
	str := reg.FindAllString(string(output), -1)
	ip = strings.TrimSpace(str[0])
	return
}

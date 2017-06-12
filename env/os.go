package env

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func IsRoot() (flag bool) {
	// also can use os/user user.Current() function.
	if 0 == os.Getuid() {
		flag = true
	}
	return
}

func Timestamp() string {
	// return string(time.Now().Unix())
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func Hostname() string {
	return "hostname=" + readProcFile("/proc/sys/kernel/hostname")
}

func Osrelease() string {
	return "osrelease=" + readProcFile("/proc/sys/kernel/osrelease")
}
func Ostype() string {
	return "ostype=" + readProcFile("/proc/sys/kernel/ostype")
}

func readProcFile(filename string) (ret string) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		return ""
	}
	defer f.Close()
	contentByte, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	ret = strings.Replace(string(contentByte), "\n", "", -1)
	return
}

/*
	for _, v := range os.Environ() { //获取全部系统环境变量 获取的是 key=val 的[]string
				str := strings.Split(v, "=")
						log.Println(str[0], str[1])
								// log.Println(os.Environ)
									}
*/

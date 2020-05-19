package pkg

import (
	"fmt"
	"strconv"
)

var (
	Version       = "000001"
	StringVersion = strVersion(Version)
	NumVersion, _ = strconv.Atoi(Version)
)

func strVersion(ver string) (strVer string) {
	strVer = string(ver[0]) + "." + string(ver[1]) + "." + string(ver[2]) + "." + string(ver[3]) + string(ver[4]) + string(ver[5])
	return strVer
}

func PrintVrsion() {
	fmt.Printf(" _____    _ _                                                  \n|  __ \\  | | |     Pdk v%s                              \n| |__) |_| | | __  Copyright © 2020 Tim_Paik <timpaik@163.com> \n|  ___/ _| | |/ /                                              \n| |  | (_| |   <   Redistributed under the terms of GNU GPL.   \n|_|   \\__,_|_|\\_\\                                              \n", StringVersion)
}

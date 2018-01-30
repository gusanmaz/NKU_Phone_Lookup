/* This terminal program makes it easy to lookup personnel at Namik Kemal University*/

package main

import (
	"github.com/urfave/cli"
	"os/exec"
	"fmt"
	"os"
	"strings"
	"github.com/anaskhan96/soup"
)

var cmdTmpl string = "curl -s 'http://rehber.nku.edu.tr/iptelefonlistesiarama.php' -H 'Host: rehber.nku.edu.tr' -H " +
	"'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:54.0) Gecko/20100101 Firefox/55.0' -H 'Accept: */*' -H " +
		"'Accept-Language: en-US,en;q=0.5' --compressed -H 'Referer: http://rehber.nku.edu.tr/' " +
			"-H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'X-Requested-With: " +
				"XMLHttpRequest' -H 'Cookie: _ga=GA1.3.546134864.1513155902; _gid=GA1.3.111753471.1517302490' -H " +
					"'Connection: keep-alive' --data 'birim=DEPARTMENT&ad=NAME'"


func main() {
	app := cli.NewApp()
	app.Name = "NKU Telephone Directory"
	app.Usage = "Look up phone numbers of NKU personnel by name or/and department by using appropriate flags"
	app.Version = "0.1"
	app.Author = "Guvenc Usanmaz"
	app.Email  = "gusanmaz <ett> gM@|l C_0M"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "name, n",
			Value:       "",
			Usage:       "Name of the personnel to query",
		},
		cli.StringFlag{
			Name:        "department, d",
			Value:       "",
			Usage:       "Name of the department to query",
		},
	}

	app.Action = func(c *cli.Context) error {
		var name = c.String("name")
		var department = c.String("department")
		commandStr := strings.Replace(cmdTmpl, "DEPARTMENT", department,1)
		commandStr = strings.Replace(commandStr, "NAME", name,1)

		cmdOutByte, _ := exec.Command("/bin/sh", "-c", commandStr).Output()
		cmdOutStr := string(cmdOutByte)

		htmlDoc := soup.HTMLParse(cmdOutStr)
		elems := htmlDoc.FindAll("td")
		elemsLen  := len(elems)
		rowCnt := (elemsLen / 4)

		seperatorLen := 0
		for id := 0; id < rowCnt; id++ {
			nameLen 	  := len(elems[4 * id + 0].Text())
			departmentLen := len(elems[4 * id + 1].Text())
			if nameLen > seperatorLen{
				seperatorLen = nameLen
			}
			if departmentLen > seperatorLen{
				seperatorLen = departmentLen
			}
		}
		seperatorLen += 20
		seperator := strings.Repeat("-", seperatorLen)

		fmt.Println("NKU PHONE LOOKUP\n")
		for id := 0; id < rowCnt; id++{
			name := elems[4 * id + 0].Text()
			department := elems[4 * id + 1].Text()
			extension  := elems[4 * id + 2].Text()
			phone      := elems[4 * id + 3].Text()
			fmt.Println(seperator)
			fmt.Println("Name:\t\t", name)
			fmt.Println("Department:\t", department)
			fmt.Println("Extension:\t", extension)
			fmt.Println("Phone:\t\t", phone)
		}
		fmt.Println(seperator)

		return nil
	}

	app.Run(os.Args)
}




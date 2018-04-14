package command

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/jiro4989/diff2xlsx/internal/config"
	"github.com/loadoff/excl"
)

const (
	rowIndex    = 2 // excelの値を追加する開始行番号
	columnIndex = 2 // excelのcellの列番号
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func CmdWrite(c *cli.Context) {
	outPath := c.String("out-file-path")
	noAttrFlag := c.Bool("no-attribute")

	// 出力先ファイルパスは必須なんでなかったら終了
	if outPath == "" {
		fmt.Fprintf(os.Stderr, "Must set \"out-file-path\" option. See --help.")
		os.Exit(2)
	}

	if err := config.InitConfigFile(); err != nil {
		log.Fatal(err)
	}
	var styleConf config.StyleConfig
	if _, err := toml.DecodeFile(config.StyleFilePath, &styleConf); err != nil {
		log.Fatal(err)
	}

	w, err := excl.Create()
	defer w.Save(outPath)
	if err != nil {
		log.Fatal(err)
	}

	s, err := w.OpenSheet("Sheet1")
	defer s.Close()
	if err != nil {
		log.Fatal(err)
	}

	// 標準入力をExcelにひたすら書き込む
	sc := bufio.NewScanner(os.Stdin)
	i := rowIndex
	for sc.Scan() {
		t := sc.Text()
		// TABがあるとexcelの列がズレるので回避
		t = strings.Replace(t, "\t", styleConf.Tab, -1)

		r := s.GetRow(i)
		c := r.GetCell(columnIndex)
		c.SetString(t)

		c.SetFont(excl.Font{Size: styleConf.Font.Size, Name: styleConf.Font.Name})
		c.SetBorder(excl.Border{
			Left:  &excl.BorderSetting{Style: "hair", Color: "000000"},
			Right: &excl.BorderSetting{Style: "hair", Color: "000000"},
		})

		// 装飾なしフラグがtrueの時は設定しない
		if !noAttrFlag {
			switch {
			case strings.HasPrefix(t, "+++"), strings.HasPrefix(t, "---"):
				c.SetFont(excl.Font{Bold: true})
				c.SetBackgroundColor("FFFFFF")
			case strings.HasPrefix(t, "@@"):
				c.SetBackgroundColor("0000FF")
			case strings.HasPrefix(t, "+"):
				c.SetBackgroundColor("00FF00")
			case strings.HasPrefix(t, "-"):
				c.SetBackgroundColor("FF0000")
			default:
				c.SetBackgroundColor("FFFFFF")
			}
		}

		i++
	}

	tc := s.GetRow(rowIndex).GetCell(columnIndex)
	bc := s.GetRow(i - 1).GetCell(columnIndex)

	tc.SetBorder(excl.Border{
		Left:  &excl.BorderSetting{Style: "hair", Color: "000000"},
		Right: &excl.BorderSetting{Style: "hair", Color: "000000"},
		Top:   &excl.BorderSetting{Style: "hair", Color: "000000"},
	})

	bc.SetBorder(excl.Border{
		Left:   &excl.BorderSetting{Style: "hair", Color: "000000"},
		Right:  &excl.BorderSetting{Style: "hair", Color: "000000"},
		Bottom: &excl.BorderSetting{Style: "hair", Color: "000000"},
	})

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
}

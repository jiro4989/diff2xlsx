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
	s.SetColWidth(float64(styleConf.Width), 2)

	// 標準入力をExcelにひたすら書き込む
	sc := bufio.NewScanner(os.Stdin)
	i := rowIndex
	for sc.Scan() {
		t := sc.Text()
		// TABがあるとexcelの列がズレるので回避
		t = strings.Replace(t, "\t", styleConf.Tab, -1)

		r := s.GetRow(i)
		r.SetHeight(float64(styleConf.Font.Size))

		c := r.GetCell(columnIndex)
		c.SetString(t)

		c.SetFont(excl.Font{Size: styleConf.Font.Size, Name: styleConf.Font.Name})
		c.SetBorder(excl.Border{
			Left:  &excl.BorderSetting{Style: styleConf.Border.Style, Color: styleConf.Border.Color},
			Right: &excl.BorderSetting{Style: styleConf.Border.Style, Color: styleConf.Border.Color},
		})

		// 装飾なしフラグがtrueの時は設定しない
		if !noAttrFlag {
			switch {
			case strings.HasPrefix(t, "+++"), strings.HasPrefix(t, "---"):
				c.SetFont(excl.Font{Size: styleConf.Font.Size, Bold: true})
				c.SetBackgroundColor("FFFFFF")
			case strings.HasPrefix(t, "-"):
				c.SetBackgroundColor(styleConf.DiffBackgroundColor.Remove)
			case strings.HasPrefix(t, "+"):
				c.SetBackgroundColor(styleConf.DiffBackgroundColor.Add)
			case strings.HasPrefix(t, "@@"):
				c.SetBackgroundColor(styleConf.DiffBackgroundColor.Range)
			default:
				c.SetBackgroundColor("FFFFFF")
			}
		}

		i++
	}

	tc := s.GetRow(rowIndex).GetCell(columnIndex)
	bc := s.GetRow(i - 1).GetCell(columnIndex)

	tc.SetBorder(excl.Border{
		Left:  &excl.BorderSetting{Style: styleConf.Border.Style, Color: styleConf.Border.Color},
		Right: &excl.BorderSetting{Style: styleConf.Border.Style, Color: styleConf.Border.Color},
		Top:   &excl.BorderSetting{Style: styleConf.Border.Style, Color: styleConf.Border.Color},
	})

	bc.SetBorder(excl.Border{
		Left:   &excl.BorderSetting{Style: styleConf.Border.Style, Color: styleConf.Border.Color},
		Right:  &excl.BorderSetting{Style: styleConf.Border.Style, Color: styleConf.Border.Color},
		Bottom: &excl.BorderSetting{Style: styleConf.Border.Style, Color: styleConf.Border.Color},
	})

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
}

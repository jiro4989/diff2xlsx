// configはtoml設定ファイルパッケージです。
package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/jiro4989/diff2xlsx/internal/version"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// excelのスタイルの設定ファイルのパス
var StyleFilePath = fmt.Sprintf("%s/.config/%s/style.toml", getEnvHome(), version.Name)

// StyleConfigはtoml設定ファイルのすべての設定を保持します。
type StyleConfig struct {
	Tab                 string
	Width               int
	Font                Font
	DiffBackgroundColor DiffBackgroundColor
	Border              Border
}

// Fontはすべてのセルのフォントを扱います。
type Font struct {
	Size int
	Name string
}

// DiffBackgroundColorはdiffテキストの背景色を扱います。
type DiffBackgroundColor struct {
	Remove string
	Add    string
	Range  string
}

// Borderは枠線の設定を扱います。
type Border struct {
	Style string
	Color string
}

// InitConfigFileはtomlファイルが存在しない時に初期ファイルを生成する。
func InitConfigFile() error {
	_, err := os.Stat(StyleFilePath)
	if err == nil {
		return nil
	}

	d := filepath.Dir(StyleFilePath)
	if err := os.MkdirAll(d, os.ModePerm); err != nil {
		log.Println(err)
		return err
	}

	w, err := os.Create(StyleFilePath)
	defer w.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	enc := toml.NewEncoder(w)

	conf := StyleConfig{
		Tab:   "    ",
		Width: 100,
		Font: Font{
			Size: 12,
			Name: "monospace",
		},
		DiffBackgroundColor: DiffBackgroundColor{
			Remove: "FF0000",
			Add:    "00FF00",
			Range:  "0000FF",
		},
		Border: Border{
			Style: "hair",
			Color: "00000",
		},
	}
	if err := enc.Encode(conf); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// getEnvHomeはホームディレクトリ文字列を返す。
func getEnvHome() string {
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		// WindowsでHOME環境変数が定義されていない場合
		home = os.Getenv("USERPROFILE")
	}
	return home
}

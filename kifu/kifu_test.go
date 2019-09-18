package kifu

import (
	//	attr "github.com/smugmug/godynamo/types/attributevalue"
	loader "KifuLibrary-Logic/kifLoader"
	f "fmt"
	"os"
	"testing"
)

func __kifu_test() {
	f.Printf("")
}

func Test_putKifu(t *testing.T) {
	file, err := os.Open("ippan.kif")

	if err != nil {
		t.Errorf("ファイルの読み込みに失敗しました。")
		return
	}

	var kif *loader.KifFile
	kif, err = loader.LoadKifFile(file)

	if kif == nil {
		t.Errorf("ファイルの読み込みに失敗しました。")
		return
	}
	if err != nil {
		t.Errorf("kifファイルの変換に失敗しました。")
		return
	}
	_, err = PutKifu(kif, "test")

	if err != nil {
		t.Errorf("DBへの書き込みにしっぱしいました。")
		return
	}
}

func Test_GetCountKifu(t *testing.T) {
	t.SkipNow()
	count, err := GetKifuCount()
	if err != nil || count == 0 {
		t.Errorf("DBエラー:%s", err.Error())
	}

	//f.Printf("current Count = %d\n", count)
}

func Test_getKifu(t *testing.T) {
	t.SkipNow()
	routes := GetKifu("99999")

	for _, route := range routes {
		f.Printf("%s / %s / %s \n", route.Prev, route.Move.ToJpnCode(), route.Current)
	}
}

func Test_getKifuListByUser(t *testing.T) {
	t.SkipNow()
	ids, err := GetKifuListByUser("uploader")
	if err != nil {
		t.Error(err.Error())
	}
	f.Printf("%v", ids)
}

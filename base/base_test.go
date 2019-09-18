package base

import (
	"testing"
)

func reset() {
	attrList := make(map[string]interface{})

	attrList["Table_Name"] = "Test"
	attrList["Seq"] = 0

	PutItem("KifCloud-Sequences", attrList)
}

func Test_GetSequence(t *testing.T) {
	reset()
	seq, err := GetSequence("Test")

	if err != nil {
		t.Errorf("シーケンスの取得に失敗しました。err = %s", err.Error())
	}
	if seq != 0 {
		t.Errorf("seqの値が期待値と異なっています。 = %d", seq)
	}
}

func Test_CountSequence(t *testing.T) {
	reset()

	err := CountSequence("Test")

	if err != nil {
		t.Errorf("シーケンスの加算に失敗しました。error＝[ %s ]", err.Error())
		return
	}

	var seq int
	seq, err = GetSequence("Test")

	if err != nil {
		t.Errorf("シーケンスの加算に失敗しました。error＝[ %s ]", err.Error())
		return
	}

	if seq != 1 {
		t.Errorf("seqの値が期待値と異なっています。 = %d", seq)
	}

}

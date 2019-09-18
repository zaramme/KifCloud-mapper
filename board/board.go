package board

import (
	b "KifCloud-Mapper/base"
	"encoding/json"
	"errors"
	"fmt"
	attr "github.com/smugmug/godynamo/types/attributevalue"
)

// update
// type BoardItem struct {
// 	Rsh  string
// 	Item map[strin]
// }

type BoardItem struct {
	Item map[string]attr.AttributeValue `json:"Item"`
}

type ItemValue struct {
	Text map[string]attr.AttributeValue `json:"text"`
}

type StringSet struct {
	S string
}

const (
	BOARD_TABLE_NAME = "KifCloud-Board"
	BOARD_KEY_NAME   = "RSH"
)

func GetBoardItem(rsh string) BoardItem { //*BoardItem {
	response, _, err := b.GetItem(BOARD_TABLE_NAME, BOARD_KEY_NAME, rsh)
	if err != nil {
		return BoardItem{}
	}

	fmt.Println(response)
	var v BoardItem //responseSet
	err = json.Unmarshal([]byte(response), &v)
	if err != nil {
		return BoardItem{}
	}

	return v
}

func PutBoardItem(attrList map[string]interface{}) error {

	if ok := attrList[BOARD_KEY_NAME]; ok == nil {
		return errors.New("ハッシュ値が設定されていません")
	}
	_, _, err := b.PutItem(BOARD_TABLE_NAME, attrList)

	return err
}

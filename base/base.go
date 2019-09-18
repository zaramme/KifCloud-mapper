package base

import (
	//"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	//ep "github.com/smugmug/godynamo/endpoint"
	"encoding/json"
	"github.com/smugmug/godynamo/endpoint"
	getitem "github.com/smugmug/godynamo/endpoints/get_item"
	putitem "github.com/smugmug/godynamo/endpoints/put_item"
	"github.com/smugmug/godynamo/endpoints/query"
	attr "github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/condition"
	"strconv"
)

type sequence struct {
	Item map[string]attr.AttributeValue `json:"Item"`
}

type ItemValue struct {
	TableName attr.AttributeValue `json:Table_Name`
	Seq       attr.AttributeValue `json:Seq`
}

func (this sequence) toInt() (n int, err error) {
	return strconv.Atoi(this.Item["Seq"].N)
}

// リクエスト前の共通処理
func willRequest() {
	conf_file.Read()
	if conf.Vals.Initialized == false {
		panic("the conf.Vals global conf struct has not been initialized")
	}
}

// リクエスト後の共通処理
func didRequest() {

}

func requestWrapper(fc func(ep *endpoint.Endpoint)) (response string, err error) {

	willRequest()

	var ep endpoint.Endpoint

	fc(&ep)

	body, _, err := ep.EndpointReq()

	didRequest()

	return body, err
}

func PutItem(tableName string, attrList map[string]interface{}) (response string, code int, err error) {

	willRequest()
	pi := putitem.NewPutItem()
	pi.TableName = tableName

	// 属性値の設定
	for k, v := range attrList {
		pi.Item[k] = attr.NewAttributeValue()

		switch ty := v.(type) {
		case string:
			//fmt.Printf("[%s]はstring型です。\n", k)
			pi.Item[k].InsertS(ty)
		case int:
			//fmt.Printf("[%s]はint型です。\n", k)
			pi.Item[k].InsertN_float64(float64(ty))
		case bool:
			//fmt.Printf("[%s]はbool型です。\n", k)
			pi.Item[k].InsertBOOL(v.(bool))
		}
	}
	body, code, err := pi.EndpointReq()
	//	fmt.Printf("%v\n%v\n%v\n", body, code, err)

	didRequest()

	return body, code, err
}

func GetItem(tableName, keyTitle, keyValue string) (response string, code int, err error) {

	willRequest()

	gi := getitem.NewGetItem()
	gi.TableName = tableName
	gi.Key[keyTitle] = attr.NewAttributeValue()
	gi.Key[keyTitle].InsertS(keyValue)

	body, code, err := gi.EndpointReq()

	didRequest()

	return body, code, err
}

func GetSequence(tableName string) (seq int, err error) {
	willRequest()

	body, _, err := GetItem("KifCloud-Sequences", "Table_Name", tableName)

	didRequest()

	if err != nil {
		return 0, err
	}

	fmt.Printf("%s\n", body)

	var v sequence
	err = json.Unmarshal([]byte(body), &v)

	fmt.Printf("%s\n", v)
	if err != nil {
		return 0, err
	}

	n, convErr := v.toInt()

	if err != nil {
		return 0, convErr
	}
	return n, nil
}
func GetNextSequence(tableName string) (seq int, err error) {
	seq, err = GetSequence(tableName)

	if err != nil {
		return 0, err
	}
	seq++
	return seq, nil
}

func CountSequence(tableName string) error {
	//現在のシーケンスを取得
	currentSeq, err := GetSequence(tableName)
	if err != nil {
		return err
	}
	updatedSeq := currentSeq + 1

	attrList := make(map[string]interface{})
	attrList["Table_Name"] = tableName
	attrList["Seq"] = updatedSeq

	_, _, err = PutItem("KifCloud-Sequences", attrList)

	if err != nil {
		return err
	}

	return nil
}

func Query(tableName string, conditions condition.Conditions) (response string, code int, err error) {
	willRequest()

	q := query.NewQuery()
	q.TableName = tableName
	for key, condition := range conditions {
		q.KeyConditions[key] = condition
	}

	body, code, err := q.EndpointReq()

	didRequest()

	return body, code, err
}

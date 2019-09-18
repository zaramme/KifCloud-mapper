package kifu

import (
	b "KifCloud-Mapper/base"
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/endpoints/query"
	attr "github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/condition"
	jki "github.com/zaramme/KifCloud-Logic/JsonConverter/JsonKifuInfo"
	"github.com/zaramme/KifCloud-Logic/kifLoader"
	m "github.com/zaramme/KifCloud-Logic/move"
	r "github.com/zaramme/KifCloud-Logic/routes"
	"strconv"
)

type KifuJson struct {
	Count int                      `json:"Count"`
	Items []attr.AttributeValueMap `json:"Items"`
}

type Itemsvalue struct {
	Move     attr.AttributeValue `json:"Move"`
	Rsh_Prev attr.AttributeValue `json:"Rsh_Prev"`
}

func GetKifu(KifuID string) r.Routes {
	conditions := condition.NewConditions()

	KifuIDConf := condition.NewCondition()
	KifuIdAttr := attr.NewAttributeValue()
	KifuIdAttr.InsertN(KifuID)
	KifuIDConf.AttributeValueList = make([]*attr.AttributeValue, 1)
	KifuIDConf.AttributeValueList[0] = KifuIdAttr
	KifuIDConf.ComparisonOperator = query.OP_EQ

	conditions["KifuID"] = KifuIDConf

	body, _, err := b.Query("KifCloud-Kifu-Route", conditions)

	if err != nil {
		fmt.Print("エラー")
		return nil
	}

	// jsonオブジェクトへ変換
	var v KifuJson
	jsonErr := json.Unmarshal([]byte(body), &v)
	if jsonErr != nil {
		fmt.Print(jsonErr)
		return nil
	}

	// // Routesオブジェクトへ変換
	KifuRoutes := r.NewRoutes()
	for _, route := range v.Items {

		move := m.NewMoveFromMoveCode(route["Move"].S)
		set := r.Set{
			Prev:    route["Rsh_Prev"].S,
			Current: route["Rsh_Current"].S,
			Move:    move,
		}

		KifuRoutes = append(KifuRoutes, set)
	}

	return KifuRoutes

}

func GetKifuListByUser(UserID string) (r *jki.Shell, err error) {
	////////////////////////////////////////////////////////
	// クエリ発行
	conditions := condition.NewConditions()

	UserIDConf := condition.NewCondition()
	UserIdAttr := attr.NewAttributeValue()
	UserIdAttr.InsertS(UserID)
	UserIDConf.AttributeValueList = make([]*attr.AttributeValue, 1)
	UserIDConf.AttributeValueList[0] = UserIdAttr
	UserIDConf.ComparisonOperator = query.OP_EQ

	conditions["UserID"] = UserIDConf

	body, _, err := b.Query("KifCloud-Kifu-Info", conditions)

	if err != nil {
		return nil, err
	}
	////////////////////////////////////////////////////////
	// jsonオブジェクトへ変換
	var v KifuJson
	jsonErr := json.Unmarshal([]byte(body), &v)
	if jsonErr != nil {
		fmt.Print(jsonErr)
		return
	}

	infoList := jki.NewKifuInfo(v.Count)

	seq := 0
	for _, kifu := range v.Items {
		info := new(jki.KifuInfo)
		if userID, ok := kifu["UserID"]; ok {
			info.UserID = userID.S
		}
		if vKifuID, ok := kifu["KifuID"]; ok {
			kifuID, err := strconv.Atoi(vKifuID.N)
			if err != nil {
				return nil, fmt.Errorf("不正な棋譜IDを検出しました。(id = %s)", vKifuID.S)
			}
			info.KifuID = kifuID
		}
		if blackPlayer, ok := kifu["先手"]; ok {
			info.UserID = blackPlayer.S
		}
		if whitePlayer, ok := kifu["後手"]; ok {
			info.WhitePlayer = whitePlayer.S
		}
		fmt.Printf("[%d] = %v\n", seq, info)
		infoList.KifuInfos[seq] = info
		seq += 1
	}
	return infoList, nil
}

func GetKifuCount() (count int, err error) {

	count, err = b.GetSequence("KifuID")

	return count, err
}

func PutKifu(kif *kifLoader.KifFile, userID string) (kifFile *kifLoader.KifFile, err error) {

	routes, err := r.NewRoutesFromKifuFile(kif)
	if err != nil {
		return nil, err
	}

	attrList := make(map[string]interface{})

	kifuID := 0
	kifuID, err = b.GetNextSequence("KifuID")

	if err != nil {
		return nil, err
	}

	fmt.Printf("next Sequence = %d\n", kifuID)
	attrList["UserID"] = userID
	attrList["KifuID"] = kifuID

	for key, value := range kif.Info {
		attrList[key] = value
	}

	_, _, err = b.PutItem("KifCloud-Kifu-Info", attrList)

	if err != nil {
		return nil, err
	}

	for seq, value := range routes {
		attrList := make(map[string]interface{})
		attrList["KifuID"] = kifuID
		attrList["Seq"] = seq
		attrList["Rsh_Prev"] = value.Prev
		attrList["Move"] = value.Move.ToMoveCode()
		attrList["Rsh_Current"] = value.Current

		_, _, err := b.PutItem("KifCloud-Kifu-Route", attrList)
		if err != nil {
			break
		}
	}

	b.CountSequence("KifuID")

	return kif, nil
}

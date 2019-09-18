package board

import (
	"fmt"
	"testing"
)

func Test_PutBoardItem(t *testing.T) {
	// attrList := map[string]interface{}{
	// 	"RSH":  string("1234567890"),
	// 	"text": string("これはテストです。"),
	// }
	// PutBoardItem(attrList)
}

func Test_GetBoardItem(t *testing.T) {
	response := GetBoardItem("Icflk2hGLCEL2UtPMR4e0ohgUgpT0AEsXcnWE0")

	//		if len(response) > 0 {
	fmt.Printf("\n\n%s\n", response)
	//		}

}

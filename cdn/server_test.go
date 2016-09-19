package cdn

import (
	"testing"
)

func TestRunMain(t *testing.T){

	Rummain()
}

func TestTmain(t *testing.T)  {
	Tmain()
}

func TestSaveData(t *testing.T){
	SaveData()
	//SaveFileByte()

	fileId :="57dfa368b315f1227b2213a5"
	OpenSavedFile(fileId )
	// 57dec1b3b315f13eb095d0b8
}

func TestSaveFileByte(t *testing.T) {
	SaveFileByte()
}

package l10n

import "testing"
func checkHasErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}


func TestCreateWithCreate(t *testing.T) {
	product := Product{Code: "CreateWithCreate"}
	checkHasErr(t, dbGlobal.Create(&product).Error)
	checkHasErr(t, dbCN.Create(&product).Error)
	checkHasErr(t, dbEN.Create(&product).Error)
}

package admin_test

import (
	"ems/admin"
	. "ems/admin/tests/dummy"
	"ems/core"
	"ems/core/resource"
	"reflect"
	"testing"
)

func TestTextInput(t *testing.T) {
	user := Admin.AddResource(&User{})

	meta := user.GetMeta("Name")

	if meta.Label != "Name" {
		t.Error("default label not set")
	}

	if meta.GetFieldName() != "Name" {
		t.Error("default Alias is not same as field Name")
	}

	if meta.Type != "string" {
		t.Error("default Type is not string")
	}
}

//测试基本的meta, 即每个resource都没有手动指定Meta, 它getMeta()自动根据model struct的字段来生成
func TestDefaultMetaType(t *testing.T) {
	var (
		user        = Admin.AddResource(&User{})
		booleanMeta = user.GetMeta("Active")
		timeMeta    = user.GetMeta("RegisteredAt")
		numberMeta  = user.GetMeta("Age")
		fileMeta    = user.GetMeta("Avatar")
	)
	if booleanMeta.Type != "checkbox" {
		t.Error("boolean field doesn't set as checkbox")
	}

	if timeMeta.Type != "datetime" {
		t.Error("time field doesn't set as datetime")
	}

	if numberMeta.Type != "number" {
		t.Error("number field doesn't set as number")
	}

	//调用了meta.config()中调用了media/base.go的ConfigureMetaBeforeInitialize方法
	if fileMeta.Type != "file" {
		t.Error("file field doesn't set as file")
	}

}

func TestRelationFieldMetaType(t *testing.T) {
	userRecord := &User{}
	db.Create(userRecord)
	user := Admin.AddResource(&User{})
	userProfileMeta := user.GetMeta("Profile") //user has one profile
	if userProfileMeta.Type != "single_edit" {
		t.Error("has_one relation doesn't generate single_edit type meta")
	}

	userAddressesMeta := user.GetMeta("Addresses") //has many

	if userAddressesMeta.Type != "collection_edit" {
		t.Error("has_many relation doesn't generate collection_edit type meta")
	}

	userLanguagesMeta := user.GetMeta("Languages") //many to many

	if userLanguagesMeta.Type != "select_many" {
		t.Error("many_to_many relation doesn't generate select_many type meta")
	}

}

//通过meta从resource中获了以一个值
func TestGetStringMetaValue(t *testing.T) {
	user := Admin.AddResource(&User{})
	stringMeta := user.GetMeta("Name")

	UserName := "user name"

	userRecord := &User{Name: UserName}

	db.Create(&userRecord)

	value := stringMeta.GetValuer()(userRecord, &core.Context{Config: &core.Config{DB: db}})

	if value.(string) != UserName {
		t.Error("resource's value doesn't get")
	}
}

func TestGetStructMetaValue(t *testing.T) {
	user := Admin.AddResource(&User{})
	structMeta := user.GetMeta("CreditCard") //belongs_to

	creditCard := CreditCard{
		Number: "123456",
		Issuer: "bank",
	}
	userRecord := &User{CreditCard: creditCard}
	db.Create(&userRecord)

	//value是createCard对像
	value := structMeta.GetValuer()(userRecord, &core.Context{Config: &core.Config{DB: db}})
	creditCardValue := reflect.Indirect(reflect.ValueOf(value))

	if creditCardValue.FieldByName("Number").String() != "123456" || creditCardValue.FieldByName("Issuer").String() != "bank" {
		t.Error("struct field value doesn't get")
	}
}

func TestGetSliceMetaValue(t *testing.T) {
	user := Admin.AddResource(&User{})
	sliceMeta := user.GetMeta("Addresses") //has many
	address1 := &Address{Address1: "an address"}
	address2 := &Address{Address1: "another address"}
	userRecord := &User{Addresses: []Address{*address1, *address2}}
	db.Create(&userRecord)
	value := sliceMeta.GetValuer()(userRecord, &core.Context{Config: &core.Config{DB: db}})
	addresses := reflect.Indirect(reflect.ValueOf(value))

	if addresses.Index(0).FieldByName("Address1").String() != "an address" || addresses.Index(1).FieldByName("Address1").String() != "another address" {
		t.Error("slice field value doesn't get")
	}
}

func TestStringMetaSetter(t *testing.T) {
	user := Admin.AddResource(&User{})
	meta := user.GetMeta("Name")

	UserName := "new name"
	userRecord := &User{}

	db.Create(&userRecord)

	metaValue := &resource.MetaValue{
		Name:  "User.Name",
		Value: UserName,
		Meta:  meta,
	}

	meta.GetSetter()(userRecord, metaValue, &core.Context{Config: &core.Config{DB: db}})
	if userRecord.Name != UserName {
		t.Error("resource's value doesn't set")
	}

}

func TestNestedField(t *testing.T) {
	profileModel := Profile{
		Name:  "Ems",
		Sex:   "Female",
		Phone: Phone{Num: "1024"},
	}

	userModel := &User{Profile: profileModel}
	db.Create(userModel)

	user := Admin.AddResource(&User{})
	profileNameMeta := &admin.Meta{Name: "Profile.Name"}
	user.Meta(profileNameMeta)
	profileSexMeta := &admin.Meta{Name: "Profile.Sex"}
	user.Meta(profileSexMeta)
	phoneNumMeta := &admin.Meta{Name: "Profile.Phone.Num"}
	user.Meta(phoneNumMeta)

	userModel.Profile = Profile{}
	valx := phoneNumMeta.GetValuer()(userModel, &core.Context{Config: &core.Config{DB: db}})
	if val, ok := valx.(string); !ok || val != profileModel.Phone.Num {
		t.Errorf("Profile.Phone.Num: got %q; expect %q", val, profileModel.Phone.Num)
	}
	if userModel.Profile.Name != profileModel.Name {
		t.Errorf("Profile.Name: got %q; expect %q", userModel.Profile.Name, profileModel.Name)
	}
	if userModel.Profile.Sex != profileModel.Sex {
		t.Errorf("Profile.Sex: got %q; expect %q", userModel.Profile.Sex, profileModel.Sex)
	}
	if userModel.Profile.Phone.Num != profileModel.Phone.Num {
		t.Errorf("Profile.Phone.Num: got %q; expect %q", userModel.Profile.Phone.Num, profileModel.Phone.Num)
	}

	mvs := &resource.MetaValues{
		Values: []*resource.MetaValue{
			{
				Name:  "Profile.Name",
				Value: "Qor III",
				Meta:  profileNameMeta,
			},
			{
				Name:  "Profile.Sex",
				Value: "Male",
				Meta:  profileSexMeta,
			},
			{
				Name:  "Profile.Phone.Num",
				Value: "2048",
				Meta:  phoneNumMeta,
			},
		},
	}
	profileNameMeta.GetSetter()(userModel, mvs.Values[0], &core.Context{Config: &core.Config{DB: db}})
	if userModel.Profile.Name != mvs.Values[0].Value {
		t.Errorf("Profile.Name: got %q; expect %q", userModel.Profile.Name, mvs.Values[0].Value)
	}
	profileSexMeta.GetSetter()(userModel, mvs.Values[1], &core.Context{Config: &core.Config{DB: db}})
	if userModel.Profile.Sex != mvs.Values[1].Value {
		t.Errorf("Profile.Sex: got %q; expect %q", userModel.Profile.Sex, mvs.Values[1].Value)
	}
	phoneNumMeta.GetSetter()(userModel, mvs.Values[2], &core.Context{Config: &core.Config{DB: db}})
	if userModel.Profile.Phone.Num != mvs.Values[2].Value {
		t.Errorf("Profile.Phone.Num: got %q; expect %q", userModel.Profile.Phone.Num, mvs.Values[2].Value)
	}

}

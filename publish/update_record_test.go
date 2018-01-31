package publish

import "testing"

func TestUpdateStructFromDraft(t *testing.T) {
	name := "update_product_from_draft"
	newName := name + "_v2"
	product := Product{Name: name, Color: Color{Name: name}}
	//在production和_draft表中创建两条记录， 两条记录的publish_status都为false或者空
	pbprod.Create(&product)

	//更新_draft表中的数据,
	pbdraft.Model(&product).Update("name", newName)

	pbdraft.First(&product, product.ID)

	if !product.PublishStatus {
		t.Errorf("在更新后，product的publish status应该为DIRTY")
	}

	if pbprod.First(&Product{}, "name=?", name).RecordNotFound() {
		t.Errorf("在draft中更新，不应该颢响到prodution")
	}

	if pbdraft.First(&Product{}, "name=?", newName).RecordNotFound() {
		t.Errorf("在draft表中的product应该为新的名字")
	}

	if pbdraft.Model(&product).Related(&product.Color);product.Color.Name != name {
		t.Errorf("should be able to find related struct")
	}
}


func TestUpdateStructFromProduction(t *testing.T) {
	name := "update_product_from_production"
	newName := name + "_v2"
	product := Product{Name: name, Color: Color{Name: name}}
	pbprod.Create(&product)
	pbprod.Model(&product).Update("name", newName)

	if product.PublishStatus {
		t.Errorf("Product's publish status should be PUBLISHED when updated from production db")
	}


}

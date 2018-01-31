package publish

import "testing"

func TestDeleteStructFromDraft(t *testing.T) {
	name := "delete_product_from_draft"
	product := Product{Name: name, Color: Color{Name: name}}
	pbprod.Create(&product)
	pbdraft.Delete(&product)

	pbdraft.Unscoped().First(&product, product.ID)

	//当从draft db中删除一条记录时，应该将它的publish_status修改为DIRTY
	if !product.PublishStatus {
		t.Errorf("Product's publish status should be DIRTY when deleted from draft db")
	}

	if pbprod.First(&Product{}, "name=?", name).RecordNotFound() {
		t.Errorf("record 应该在production db中被删除了")
	}

	if !pbdraft.First(&Product{}, "name=?", name).RecordNotFound() {
		t.Errorf("record should be soft deleted in draft db")
	}

}


func TestDeleteStructFromProduction(t *testing.T) {
	name := "delete_product_from_production"
	product := Product{Name: name, Color: Color{Name: name}}
	pbprod.Create(&product)
	pbprod.Delete(&product)
	//找到不包含被安全删除的记录deleted_at为null
	if !pbprod.First(&Product{}, "name = ?", name).RecordNotFound() {
		t.Errorf("record should be soft deleted in production db")
	}


	//找到所有记录，包含被安全删除的记录
	if pbprod.Unscoped().First(&Product{}, "name=?", name).RecordNotFound() {
		t.Errorf("record should be soft deleted in production db")
	}

	//如果是从production中删除一条记录，它也会删除draft中的数据
	if !pbdraft.First(&Product{}, "name = ?", name).RecordNotFound() {
		t.Errorf("record should be soft deleted in draft db")
	}

	if pbdraft.Unscoped().First(&Product{}, "name = ?", name).RecordNotFound() {
		t.Errorf("record should be soft deleted in draft db")
	}
	//如果是从production中删除一条记录，它也会删除draft中的数据, 但不会像draft模式那样，还需要修改draft表中的publish_status为true(DIRTY)

	pbdraft.Unscoped().First(&product, product.ID)
	if product.PublishStatus {
		t.Errorf("Product's publish status should be PUBLISHED when deleted from production db")
	}

	pbprod.Unscoped().First(&product, product.ID)
	if product.PublishStatus {
		t.Errorf("Product's publish status should be PUBLISHED when deleted from production db")
	}
}
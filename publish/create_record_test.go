package publish

import (
	"testing"
)

func TestCreateStructForDraft(t *testing.T) {
	name := "create_product from draft"
	pbdraft.Create(&Product{Name: name, Color: Color{Name: name}})

	if !pbprod.First(&Product{}, "name=?", name).RecordNotFound() {
		t.Errorf("此记录不能在production db中找到")
	}

	if pbdraft.First(&Product{}, "name=?", name).RecordNotFound() {
		t.Errorf("此记录应该存在于pbdraft")
	}

	if pbprod.Table("colors").First(&Color{}, "name = ?", name).Error != nil {
		t.Errorf("color 应该保存在production")
	}

	if pbprod.Table("colors_draft").First(&Color{}, "name = ?", name).Error == nil {
		t.Errorf("colors_draft 不应该保存记录，应该production它不是mang2many的关系, colors_draft表不应该出现")
	}

	var product Product
	pbdraft.First(&product, "name=?", name)

	if !product.PublishStatus {
		t.Errorf("从draft db 中创建的 Product的publish_state列应该为DIRTY (true),")
	}

	//默认只找到product表的信息，即只包含颜色id, 而Color为空
	if product.ColorId == 0 {
		t.Errorf("ColorID 应该可以从product中找回")
	}
	if product.Color.Name != "" {
		t.Error("当前Color信息还没有获取，需要通过Related获取")
	}
	pbdraft.Model(&product).Related(&product.Color)

	if product.Color.Name != name {
		t.Errorf("通过product的ColorID找到color相关的信息")
	}
}

func TestCreateStructForProduction(t *testing.T) {
	name := "create product form production"
	pbprod.Create(&Product{Name: name, Color: Color{Name: name}})
	if pbprod.First(&Product{}, "name = ?", name).RecordNotFound() {
		t.Errorf("record should be found in production db")
	}

	if pbdraft.First(&Product{}, "name = ?", name).RecordNotFound() {
		t.Errorf("record should be found in draft db")
	}

	if pbprod.Table("colors").First(&Color{}, "name = ?", name).Error != nil {
		t.Errorf("color should be saved")
	}


	var product Product
	pbprod.First(&product, "name = ?", name)

	if product.PublishStatus {
		t.Errorf("通过pbprod创建的记录，它的publish_status应该为  PUBLISHED （false)")
	}

	if pbprod.Model(&product).Related(&product.Color); product.Color.Name != name {
		t.Errorf("should be able to find related struct")
	}
}

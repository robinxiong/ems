package publish

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
)

//创建发布单个资源的类型，它实现EventInterface接口
type createResourcePublishInterface struct {
}

func (createResourcePublishInterface) Publish(db *gorm.DB, event PublishEventInterface) error {
	if event, ok := event.(*PublishEvent); ok {
		var product Product
		db.Set("publish:draft_mode", true).First(&product, event.Argument)

		pb.Publish(&product)
	}
	return nil
}
func (createResourcePublishInterface) Discard(db *gorm.DB, event PublishEventInterface) error {
	if event, ok := event.(*PublishEvent); ok {
		var product Product
		db.Set("publish:draft_mode", true).First(&product, event.Argument)
		pb.Discard(&product)
	}
	return nil
}

type publishAllResourcesInterface struct {
}

func (publishAllResourcesInterface) Publish(db *gorm.DB, event PublishEventInterface) error {
	return nil
}

func (publishAllResourcesInterface) Discard(db *gorm.DB, event PublishEventInterface) error {
	return nil
}

func init() {
	RegisterEvent("create_product", createResourcePublishInterface{})
	RegisterEvent("publish_all_resources", publishAllResourcesInterface{})
}

func TestCreateNewEvent(t *testing.T) {
	product1 := Product{Name: "event_1", Brand: Brand{Name: "event_1_brand"}}
	pbdraft.Set("publish:publish_event", true).Save(&product1)
	//创建一个发布事件model,  保存到publish_events表中, 通过名字查找在init中注册的createResourcePublishInterface
	event := PublishEvent{Name: "create_product", Argument: fmt.Sprintf("%v", product1.ID)}
	db.Save(&event)

	if !pbprod.First(&Product{}, "name = ?", product1.Name).RecordNotFound() {
		t.Errorf("通过event创建创建的资源只在draft db, 还没有发布到production db中")
	}

	var productDraft Product
	if pbdraft.First(&productDraft, "name = ?", product1.Name).RecordNotFound() {
		t.Errorf("created resource in draft db with event should exist in draft db")
	}

	if productDraft.PublishStatus == DIRTY {
		t.Errorf("在publish event发布之前，_draft中产品的publish_status不能为DIRTY")
	}

	//查找publish_event表中的记录
	var publishEvent PublishEvent
	if pbdraft.First(&publishEvent, "name = ?", "create_product").Error != nil {
		t.Errorf("created resource in draft db with event should create the event in db")
	}

	if !pbprod.First(&Product{}, "name = ?", product1.Name).RecordNotFound() {
		t.Errorf("product should not be published to production db before publish event")
	}

	publishEvent.Publish(db)

	if pbprod.First(&Product{}, "name = ?", product1.Name).RecordNotFound() {
		t.Errorf("product should be published to production db after publish event")
	}
}

func TestCreateProductWithPublishAllEvent(t *testing.T) {
	product1 := Product{Name: "event_1"}
	event := &PublishEvent{Name: "publish_all_resources", Argument: "products"}
	pbdraft.Set("publish:publish_event", event).Save(&product1)
}

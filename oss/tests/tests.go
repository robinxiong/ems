package tests

import (
	"ems/oss"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"os"
	"io/ioutil"
)

func TestAll(storage oss.StorageInterface, t *testing.T) {
	randomPath := strings.Replace(time.Now().Format("20060102150506.000"), ".", "", -1)
	fmt.Printf("testing file in %v\n", filepath.Join(storage.GetEndpoint(), randomPath)) //filesystem的端点为 "/"

	fileName := "/" + filepath.Join(randomPath, "sample.txt")
	fileName2 := "/" + filepath.Join(randomPath, "sample2", "sample.txt")
	exceptObjects := 2
	sampleFile, _ := filepath.Abs("../tests/sample.txt")

	//测试创建文件
	//打开测试的sample.txt文件
	if file, err := os.Open(sampleFile); err == nil {
		if object, err := storage.Put(fileName, file); err != nil {
			t.Errorf("No error should happen when save sample file, but got %v", err)
		} else if object.Path == "" || object.StorageInterface == nil {
			t.Errorf("返回的资源对像，必须包含必要的信息")
		}
	} else {
		t.Errorf("No error should happen when opem sample file, but got %v", err)
	}

	//保存到sample2文件夹
	if file, err := os.Open(sampleFile); err == nil {
		if object, err := storage.Put(fileName2, file); err != nil {
			t.Errorf("No error should happen when save sample file, but got %v", err)
		} else if object.Path == "" || object.StorageInterface == nil {
			t.Errorf("returned object should necessary information")
		}
	} else {
		t.Errorf("No error should happen when opem sample file, but got %v", err)
	}

	// Get file
	if file, err := storage.Get(fileName); err != nil {
		t.Errorf("No error should happen when get sample file, but got %v", err)
	} else {
		if buffer, err := ioutil.ReadAll(file); err != nil {
			t.Errorf("No error should happen when read downloaded file, but got %v", err)
		} else if string(buffer) == "sample" {
			t.Errorf("Downloaded file should contain correct content, but got %v", string(buffer))
		}
	}


	// List
	if objects, err := storage.List(randomPath); err != nil {
		t.Errorf("No error should happen when list objects, but got %v", err)
	} else if len(objects) != exceptObjects {
		t.Errorf("Should found %v objects, but got %v", exceptObjects, len(objects))
	} else {
		var found1, found2 bool
		for _, object := range objects {
			if object.Path == fileName {
				found1 = true
			}

			if object.Path == fileName2 {
				found2 = true
			}
		}

		if !found1 {
			t.Errorf("Should found uploaded file %v", fileName)
		}

		if !found2 {
			t.Errorf("Should found uploaded file %v", fileName2)
		}
	}

	// Delete
	if err := storage.Delete(fileName); err != nil {
		t.Errorf("No error should happen when delete sample file, but got %v", err)
	}

	// Get file after delete
	if _, err := storage.Get(fileName); err == nil {
		t.Errorf("There should be an error when get deleted sample file")
	}

	// Get file after delete
	if _, err := storage.Get(fileName2); err != nil {
		t.Errorf("Sample file 2 should no been deleted")
	}


}
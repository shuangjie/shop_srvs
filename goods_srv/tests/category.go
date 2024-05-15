package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"os"
)

func TestGetCategoryList() {
	rsp, err := goodsClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.JsonData)

	// 便于观察，将 rsp.JsonData 写入 JSON 文件
	//err = writeJSONToFile("category-test.json", []byte(rsp.JsonData))
	//if err != nil {
	//	log.Fatalf("Error writing JSON data to file: %v", err)
	//}
}

// 将 JSON 数据写入文件
func writeJSONToFile(filename string, jsonData []byte) error {
	// 将 JSON 数据写入文件
	err := os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	// 输出提示信息
	log.Printf("JSON data has been written to file: %s", filename)
	return nil
}

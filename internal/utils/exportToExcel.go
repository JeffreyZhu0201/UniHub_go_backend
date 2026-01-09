package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExportToExcel 将结构体切片导出为 Excel 文件并保存到 resources/Export 目录
// data: 结构体切片 (例如 []model.User)
// filePrefix: 生成文件的前缀名
// 返回: 生成文件的相对路径 (例如 "resources/Export/users_123456789.xlsx"), error
func ExportToExcel(data interface{}, filePrefix string) (string, error) {
	sliceVal := reflect.ValueOf(data)
	if sliceVal.Kind() != reflect.Slice {
		return "", fmt.Errorf("data must be a slice")
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 默认 Sheet1
	sheetName := "Sheet1"
	// Create a new sheet.
	index, err := f.NewSheet(sheetName)
	if err != nil {
		// handle potential error or ignore if it just says it exists
	}

	if sliceVal.Len() >= 0 { // Allow empty slice to generate headers if type is known
		// 获取元素类型
		elemType := sliceVal.Type().Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}

		if elemType.Kind() != reflect.Struct {
			// If not a struct (e.g. slice of ints), we can't easily make headers.
			// Just skipping logic or return error.
			// Assuming usage for DTO/DB Models.
			return "", fmt.Errorf("elements must be structs")
		}

		// 写入表头
		numFields := elemType.NumField()
		for i := 0; i < numFields; i++ {
			field := elemType.Field(i)
			header := field.Name
			// 优先使用 json tag 作为表头
			if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
				header = tag
			} else if tag := field.Tag.Get("gorm"); tag != "" && tag != "-" {
				// 简单的尝试获取 gorm column name, 实际解析比较复杂，这里仅简单 fallback
				header = field.Name
			}

			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			_ = f.SetCellValue(sheetName, cell, header)
		}

		// 写入数据
		for i := 0; i < sliceVal.Len(); i++ {
			item := sliceVal.Index(i)
			if item.Kind() == reflect.Ptr {
				item = item.Elem()
			}
			for j := 0; j < numFields; j++ {
				val := item.Field(j).Interface()

				// 处理特定类型格式化 (如 Time)
				if t, ok := val.(time.Time); ok {
					val = t.Format("2006-01-02 15:04:05")
				}
				// GORM Model 中的 DeletedAt, CreatedAt 等

				cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
				_ = f.SetCellValue(sheetName, cell, val)
			}
		}
	}

	f.SetActiveSheet(index)

	// 确保目录存在
	exportDir := filepath.Join("resources", "Export")
	if err := os.MkdirAll(exportDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// 生成文件名
	filename := fmt.Sprintf("%s_%d.xlsx", filePrefix, time.Now().UnixMilli())
	filePath := filepath.Join(exportDir, filename)

	if err := f.SaveAs(filePath); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return filePath, nil // 返回相对路径
}

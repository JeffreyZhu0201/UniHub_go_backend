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

	if sliceVal.Len() > 0 {
		// 检查第一个元素以确定处理逻辑 (支持 Struct 和 Map)
		firstVal := sliceVal.Index(0)
		for firstVal.Kind() == reflect.Interface || firstVal.Kind() == reflect.Ptr {
			firstVal = firstVal.Elem()
		}

		if firstVal.Kind() == reflect.Map {
			// --- Map 处理逻辑 ---
			keys := firstVal.MapKeys()
			// 写入表头
			for i, key := range keys {
				cell, _ := excelize.CoordinatesToCellName(i+1, 1)
				_ = f.SetCellValue(sheetName, cell, key.String())
			}
			// 写入数据
			for i := 0; i < sliceVal.Len(); i++ {
				item := sliceVal.Index(i)
				for item.Kind() == reflect.Interface || item.Kind() == reflect.Ptr {
					item = item.Elem()
				}
				for j, key := range keys {
					val := item.MapIndex(key)
					if val.IsValid() {
						cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
						// 格式化时间
						if t, ok := val.Interface().(time.Time); ok {
							_ = f.SetCellValue(sheetName, cell, t.Format("2006-01-02 15:04:05"))
						} else {
							_ = f.SetCellValue(sheetName, cell, val.Interface())
						}
					}
				}
			}
		} else if firstVal.Kind() == reflect.Struct {
			// --- Struct 处理逻辑 (反射运行时类型) ---
			elemType := firstVal.Type()
			numFields := elemType.NumField()

			// 写入表头
			for i := 0; i < numFields; i++ {
				field := elemType.Field(i)
				header := field.Name
				if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
					header = tag
				} else if tag := field.Tag.Get("gorm"); tag != "" && tag != "-" {
					header = field.Name
				}
				cell, _ := excelize.CoordinatesToCellName(i+1, 1)
				_ = f.SetCellValue(sheetName, cell, header)
			}

			// 写入数据
			for i := 0; i < sliceVal.Len(); i++ {
				item := sliceVal.Index(i)
				for item.Kind() == reflect.Interface || item.Kind() == reflect.Ptr {
					item = item.Elem()
				}
				for j := 0; j < numFields; j++ {
					val := item.Field(j).Interface()
					if t, ok := val.(time.Time); ok {
						val = t.Format("2006-01-02 15:04:05")
					}
					cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
					_ = f.SetCellValue(sheetName, cell, val)
				}
			}
		} else {
			return "", fmt.Errorf("unsupported element type: %v", firstVal.Kind())
		}

	} else {
		// 尝试处理空 Slice (仅当明确是 Struct Slice 时)
		elemType := sliceVal.Type().Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if elemType.Kind() == reflect.Struct {
			numFields := elemType.NumField()
			for i := 0; i < numFields; i++ {
				field := elemType.Field(i)
				header := field.Name
				if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
					header = tag
				}
				cell, _ := excelize.CoordinatesToCellName(i+1, 1)
				_ = f.SetCellValue(sheetName, cell, header)
			}
		}
		// 如果是空的 []interface{}，无法确定表头，生成空文件
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

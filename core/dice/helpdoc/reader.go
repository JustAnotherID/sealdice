package helpdoc

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func Read(group string, path string) ([]HelpTextItem, error) {
	fileExt := filepath.Ext(path)
	switch fileExt {
	case ".json":
		return readJsonHelpDoc(group, path)
	case ".xlsx":
		// 梨骰帮助文件
		return readXlsxHelpDoc(group, path)
	default:
		return nil, fmt.Errorf("unsupported help doc format: `%s`", fileExt)
	}
}

func readJsonHelpDoc(group string, path string) ([]HelpTextItem, error) {
	var result []HelpTextItem

	data := HelpDocFormat{}
	pack, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(pack, &data)
	if err != nil {
		return result, err
	}

	for k, v := range data.Helpdoc {
		result = append(result, HelpTextItem{
			Group:       group,
			From:        path,
			Title:       k,
			Content:     v,
			PackageName: data.Mod,
		})
	}
	return result, nil
}

func readXlsxHelpDoc(group string, path string) ([]HelpTextItem, error) {
	var result []HelpTextItem

	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	var errs []string
	for index, s := range f.GetSheetList() {
		rows, err := f.GetRows(s)
		if err == nil {
			var synonymCount int
			for i, row := range rows {
				if i == 0 {
					synonymCount, err = validateXlsxHeaders(row)
					if err == nil {
						// 跳过第一行
						continue
					} else {
						errs = append(errs, fmt.Sprintf("%s sheet %d(zero-based): %s", path, index, err))
						break
					}
				}
				if len(row) < 3 {
					continue
				}
				key := row[0]
				for j := range synonymCount {
					if len(row[1+j]) > 0 {
						key += "/" + row[1+j]
					}
				}
				content := row[synonymCount+1]

				result = append(result, HelpTextItem{
					Group:       group,
					From:        path,
					Title:       key,
					Content:     content,
					PackageName: s,
				})
			}
		}
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}
	return result, nil
}

// validateXlsxHeaders 验证 xlsx 格式 helpdoc 的表头是否是 Key Synonym（可能有多列） Content Description Catalogue Tag
func validateXlsxHeaders(headers []string) (int, error) {
	if len(headers) < 3 {
		return 0, errors.New("helpdoc格式错误，缺少必须列 Key Synonym Content")
	}

	var (
		index    int
		expected string
	)
	var synonymCount int
	expected = "key"
out:
	for index < len(headers) {
		// 放宽同义词大小写校验
		header := strings.ToLower(headers[index])
		switch expected {
		case "key":
			if header != "key" {
				return 0, fmt.Errorf("helpdoc表头格式错误，第%d列表头必须是Key，当前为%s", index+1, header)
			}
			expected = "synonym"
			index++
		case "synonym":
			if header != "synonym" {
				return 0, fmt.Errorf("helpdoc表头格式错误，第%d列表头必须是Synonym，当前为%s", index+1, header)
			}
			expected = "content"
			index++
			synonymCount++
		case "content":
			if header == "" || header == "synonym" {
				// 有多列同义词
				index++
				synonymCount++
				continue
			} else if header != "content" {
				return 0, fmt.Errorf("helpdoc表头格式错误，第%d列表头必须是为空白（表示同义词列）或者Content，当前为%s", index+1, header)
			}
			expected = "description"
			index++
		case "description":
			if header != "description" {
				return 0, fmt.Errorf("helpdoc表头格式错误，第%d列表头必须是Description，当前为%s", index+1, header)
			}
			expected = "catalogue"
			index++
		case "catalogue":
			if header != "catalogue" {
				return 0, fmt.Errorf("helpdoc表头格式错误，第%d列表头必须是Catalogue，当前为%s", index+1, header)
			}
			expected = "tag"
			index++
		case "tag":
			if header != "tag" {
				return 0, fmt.Errorf("helpdoc表头格式错误，第%d列表头必须是Tag", index+1)
			}
			break out
		default:
			return 0, fmt.Errorf("错误的表头校验状态，当前等待表头%s，实际获得%s", expected, header)
		}
	}
	return synonymCount, nil
}

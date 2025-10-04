package dice

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"sealdice-core/dice/helpdoc"
	"sealdice-core/dice/helpdoc/engine"
	"sealdice-core/logger"

	"gopkg.in/yaml.v3"

	nanoid "github.com/matoous/go-nanoid/v2"
)

type HelpManager struct {
	lock sync.RWMutex

	CurID        uint64
	EngineType   helpdoc.EngineType
	LoadingFn    string
	HelpDocTree  []*helpdoc.HelpDoc
	GroupAliases map[string]string
	// SearchEngine
	searchEngine helpdoc.SearchEngine

	Config *helpdoc.HelpConfig
}

func (m *HelpManager) Load(internalCmdMap CmdMapCls, extList []*ExtInfo) {
	log := zap.S().Named(logger.LogKeyHelpDoc)
	m.loadSearchEngine()

	_ = m.AddItem(helpdoc.HelpTextItem{
		Group: helpdoc.HelpBuiltinGroup,
		Title: "骰点",
		Content: `.help 骰点：
 .r  //丢一个100面骰
.r d10 //丢一个10面骰(数字可改)
.r 3d6 //丢3个6面骰(数字可改)
.ra 侦查 //侦查技能检定
.ra 侦查+10 //技能临时加值检定
.ra 3#p 射击 // 连续射击三次`,
		PackageName: "帮助",
	})

	_ = m.AddItem(helpdoc.HelpTextItem{
		Group: helpdoc.HelpBuiltinGroup,
		Title: "扩展",
		Content: `.help 扩展：
扩展功能可以让你开关部分指令。
例如你希望你的骰子是纯TRPG骰，那么可以通过.ext xxx off关闭一系列娱乐模块。
或者目前正在进行dnd5e游戏，你可以通过如下指令开关dnd特化扩展。COC亦然。
注意一点，不同扩展允许存在同名指令，例如dnd和coc都有st和rc，但他们本质上不是同一个指令，并不通用，还请注意。

.ext coc7 on // 打开coc7版扩展
.ext dnd5e off // 关闭dnd5版扩展

.ext dnd5e on // 打开dnd5版扩展
.ext coc7 off // 关闭coc7版扩展
`,
		PackageName: "帮助",
	})

	_ = m.AddItem(helpdoc.HelpTextItem{
		Group: helpdoc.HelpBuiltinGroup,
		Title: "跑团",
		Content: `.help 跑团：
.st 力量50 //载入技能/属性
.coc // coc7版人物做成
.dnd // dnd5版任务做成
.pc new <角色名> // 创建角色并自动绑卡，无角色名则为当前
.pc tag <角色名> // 当前群绑卡/解除绑卡(不填角色名)
.pc save <角色名> // 保存角色[不绑卡时需要手动保存]，无角色名则为当前
.pc load <角色名> // 加载角色[不绑卡]，无角色名则为当前
.pc list //列出当前角色
.pc del <角色名> //删除角色
.setcoc 2 //设置为coc2版房规
.nn 张三 //将自己的角色名设置为张三
`,
		PackageName: "帮助",
	})

	m.HelpDocTree = make([]*helpdoc.HelpDoc, 0)
	entries, err := os.ReadDir("data/helpdoc")
	if err != nil {
		log.Errorf("unable to read helpdoc folder: %v", err)
	}
	start := time.Now() // 获取当前时间
	totalEntries := len(entries)
	for i, entry := range entries {
		progress := float64(i+1) / float64(totalEntries) * 100
		log.Infof("处理用户定义帮助文档组[文件夹]: 当前帮助文档加载进度: %s %.2f%% (%d/%d)", entry.Name(), progress, i+1, totalEntries)
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if filepath.Base(entry.Name()) == helpdoc.HelpConfigFilename {
			m.loadHelpConfig()
			continue
		}
		var child helpdoc.HelpDoc
		child.Key = generateHelpDocKey()
		child.Name = entry.Name()
		child.Path = path.Join("data/helpdoc", entry.Name())
		child.IsDir = entry.IsDir()
		if child.IsDir {
			child.Group = entry.Name()
			child.Type = "dir"
			child.Children = make([]*helpdoc.HelpDoc, 0)
		} else {
			child.Group = "default"
			child.Type = filepath.Ext(child.Path)
		}
		buildHelpDocTree(&child, func(d *helpdoc.HelpDoc) {
			if !d.IsDir {
				ok := m.loadHelpDoc(d.Group, d.Path)
				// TODO: Batch过大好像不会释放……
				err = m.AddItemApply(false)
				if ok && err == nil {
					d.LoadStatus = helpdoc.Loaded
				} else {
					d.LoadStatus = helpdoc.LoadError
				}
			}
		})
		m.HelpDocTree = append(m.HelpDocTree, &child)
	}
	err = m.AddItemApply(false)
	if err != nil {
		log.Errorf("加载用户自定义帮助文档出现异常!: %v", err)
	}
	log.Infof("用户定义的帮助文档组已加载完成!")
	log.Infof("正在处理指令相关（含插件）帮助文档组")
	err = m.addInternalCmdHelp(internalCmdMap)
	if err != nil {
		log.Errorf("加载内置指令帮助文档出现异常: %v", err)
	}
	err = m.AddItemApply(false)
	if err != nil {
		log.Errorf("加载内置指令帮助文档出现异常: %v", err)
	}
	err = m.addExternalCmdHelp(extList)
	if err != nil {
		log.Errorf("加载插件指令帮助文档出现异常: %v", err)
	}
	err = m.AddItemApply(true)
	if err != nil {
		log.Errorf("加载插件指令帮助文档出现异常: %v", err)
	}
	log.Infof("指令相关（含插件）帮助文档组已加载完成!")
	m.CurID = m.searchEngine.GetTotalID()
	elapsed := time.Since(start) // 计算执行时间
	log.Infof("帮助文档加载完毕，共耗费时间: %s 共计加载条目:%d\n", elapsed, m.CurID)
}

func (m *HelpManager) loadSearchEngine() {
	if runtime.GOARCH == "arm64" {
		m.EngineType = helpdoc.Bleve
	}
	// 删除旧版本数据，这里先不改，先集中精力测试BleveSearch
	indexDir := "./data/_index"
	_ = os.RemoveAll(indexDir)
	indexDir = "./_help_cache"
	_ = os.RemoveAll(indexDir)
	switch m.EngineType {
	case helpdoc.Bleve:
	case helpdoc.Clover:
		bleveEngine, err := engine.NewBleveEngine()
		if err != nil {
			log := zap.S().Named(logger.LogKeyHelpDoc)
			log.Errorf("初始化帮助文档失败，帮助文档不可用!")
			return
		}
		m.searchEngine = bleveEngine
	default:
		panic("unhandled default case")
	}
}

func (m *HelpManager) loadHelpConfig() {
	data, err := os.ReadFile(filepath.Join("./data/helpdoc", helpdoc.HelpConfigFilename))
	if err != nil {
		panic(err)
	}
	var config helpdoc.HelpConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	m.Config = &config
	m.refreshHelpGroupAliases(config)
}

func (m *HelpManager) refreshHelpGroupAliases(config helpdoc.HelpConfig) {
	// 先清空旧的别名
	m.GroupAliases = make(map[string]string)
	for group, aliases := range config.Aliases {
		if len(aliases) > 0 {
			for _, alias := range aliases {
				m.GroupAliases[alias] = group
			}
		}
	}
}

func (m *HelpManager) Close() {
	// 关闭Bucket，并删除所有数据
	// TODO:暂时先不动删除逻辑
	m.searchEngine.Close()
	_ = os.RemoveAll("./_help_cache")
}

func (m *HelpManager) SaveHelpConfig(config *helpdoc.HelpConfig) error {
	m.Config = config
	m.refreshHelpGroupAliases(*config)

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join("./data/helpdoc", helpdoc.HelpConfigFilename), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (m *HelpManager) loadHelpDoc(group string, path string) bool {
	log := zap.S().Named(logger.LogKeyHelpDoc)
	m.lock.Lock()
	defer m.lock.Unlock()

	m.LoadingFn = path
	items, err := helpdoc.Read(group, path)
	if err != nil {
		log.Error("HelpManager.loadHelpDoc: ", err)
		return false
	}
	for _, item := range items {
		_ = m.AddItem(item)
	}
	return true
}

func (m *HelpManager) addCmdMap(packageName string, cmdMap CmdMapCls) error {
	log := zap.S().Named(logger.LogKeyHelpDoc)
	for k, v := range cmdMap {
		content := v.Help
		if content == "" {
			content = v.ShortHelp
		}
		err := m.AddItem(helpdoc.HelpTextItem{
			Group:       helpdoc.HelpBuiltinGroup,
			Title:       k,
			Content:     content,
			PackageName: packageName,
		})
		if err != nil {
			log.Errorf("AddCmdMapItem err:%v", err)
			return err
		}
	}
	return nil
}

func (m *HelpManager) addInternalCmdHelp(cmdMap CmdMapCls) error {
	err := m.addCmdMap("核心指令", cmdMap)
	if err != nil {
		return err
	}
	return nil
}

func (m *HelpManager) addExternalCmdHelp(ext []*ExtInfo) error {
	for _, i := range ext {
		err := m.AddItem(helpdoc.HelpTextItem{
			Group:       helpdoc.HelpBuiltinGroup,
			Title:       i.Name,
			Content:     i.GetDescText(i),
			PackageName: "扩展模块",
		})
		if err != nil {
			return err
		}
		err = m.addCmdMap(i.Name, i.CmdMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *HelpManager) AddItem(item helpdoc.HelpTextItem) error {
	_, err := m.searchEngine.AddItem(item)
	return err
}

func (m *HelpManager) AddItemApply(end bool) error {
	err := m.searchEngine.Apply(end)
	if err != nil {
		return err
	}
	return nil
}

func (m *HelpManager) Search(ctx *MsgContext, text string, titleOnly bool, pageSize, pageNum int, group string) (res helpdoc.SearchResult, total, pageStart, pageEnd int, err error) {
	return m.searchEngine.Search(ctx.Group.HelpPackages, text, titleOnly, pageSize, pageNum, group)
}

func (m *HelpManager) GetSuffixText() string {
	return fmt.Sprintf("(本次搜索由%s完成)", m.searchEngine.Name())
}

func (m *HelpManager) GetPrefixText() string {
	return fmt.Sprintf("[%s] ", m.searchEngine.Name())
}

func (m *HelpManager) GetContent(item *helpdoc.HelpTextItem, depth int) string {
	if depth > 7 {
		return "{递归层数过多，不予显示}"
	}
	txt := item.Content
	re := regexp.MustCompile(`\{[^}\n]+\}`)
	matched := re.FindAllStringSubmatchIndex(txt, -1)
	if len(matched) == 0 {
		return txt
	}

	result := strings.Builder{}
	formattedIdx := 0
	for _, i := range matched {
		left := i[0]
		right := i[1]

		if left != 0 && txt[left-1] == '\\' {
			result.WriteString(txt[formattedIdx : left-1])
			if right > 1 && txt[right-2] == '\\' {
				result.WriteString(txt[left : right-2])
				result.WriteByte('}')
			} else {
				result.WriteString(txt[left:right])
			}
			formattedIdx = right
			continue
		}

		result.WriteString(txt[formattedIdx:left])
		formattedIdx = right
		name := txt[left+1 : right-1]
		// 搜索TitleOnly，严格匹配Title的情形
		// 如果查询到对应数据，那么就调用m.GetContent
		valueResult, err := m.searchEngine.GetHelpTextItemByTermTitle(name)
		if err != nil {
			result.WriteByte('{')
			result.WriteString(name)
			result.WriteString(" - 未能找到}")
		} else {
			result.WriteString(m.GetContent(valueResult, depth+1))
		}
	}
	result.WriteString(txt[formattedIdx:])
	return result.String()
}

func generateHelpDocKey() string {
	key, _ := nanoid.Generate("0123456789abcdef", 16)
	return key
}

// 修改 buildHelpDocTree 函数签名，添加进度参数
func buildHelpDocTree(node *helpdoc.HelpDoc, fn func(d *helpdoc.HelpDoc)) {
	// 收集所有节点
	allNodes := []*helpdoc.HelpDoc{node}

	for i := 0; i < len(allNodes); i++ {
		current := allNodes[i]

		p, err := os.Stat(current.Path)
		if err != nil {
			continue
		}

		if !p.IsDir() {
			continue
		}

		subs, err := os.ReadDir(current.Path)
		if err != nil {
			continue
		}

		current.Children = make([]*helpdoc.HelpDoc, 0)

		for _, sub := range subs {
			if strings.HasPrefix(sub.Name(), ".") {
				continue
			}

			var child helpdoc.HelpDoc
			child.Key = generateHelpDocKey()
			child.Name = sub.Name()
			child.Path = path.Join(current.Path, sub.Name())
			child.Group = current.Group
			child.IsDir = sub.IsDir()

			if sub.IsDir() {
				child.Type = "dir"
				child.Children = make([]*helpdoc.HelpDoc, 0)
			} else {
				child.Type = filepath.Ext(sub.Name())
			}

			allNodes = append(allNodes, &child)
			current.Children = append(current.Children, &child)
		}
	}
	for _, current := range allNodes {
		// 调用处理函数
		fn(current)
	}
}

func (m *HelpManager) UploadHelpDoc(src io.Reader, group string, name string) error {
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	group = strings.ReplaceAll(group, "/", "_")
	group = strings.ReplaceAll(group, "\\", "_")
	if group == "default" {
		// 默认组直接上传到helpdoc文件夹中
		group = ""
	}

	dirPath := filepath.Join("./data/helpdoc", group)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dirPath, name)
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	var groupExists bool
	for _, groupDir := range m.HelpDocTree {
		if groupDir.Name == group {
			groupExists = true
			groupDir.Deleted = false

			var fileExists bool
			for _, child := range groupDir.Children {
				if child.Name == name && filepath.Clean(child.Path) == filepath.Clean(filePath) && !child.Deleted {
					if child.LoadStatus == helpdoc.Unload {
						child.Deleted = false
						fileExists = true
					} else {
						child.Deleted = true
					}
				}
			}
			if !fileExists {
				groupDir.Children = append(groupDir.Children, &helpdoc.HelpDoc{
					Key:   generateHelpDocKey(),
					Name:  name,
					Path:  filePath,
					Group: group,
					Type:  filepath.Ext(filePath),
				})
			}
			break
		}
	}
	if !groupExists {
		if group != "" {
			newGroupDir := helpdoc.HelpDoc{
				Key:      generateHelpDocKey(),
				Name:     group,
				Path:     dirPath,
				Group:    group,
				Type:     "dir",
				IsDir:    true,
				Children: make([]*helpdoc.HelpDoc, 0),
			}
			newGroupDir.Children = append(newGroupDir.Children, &helpdoc.HelpDoc{
				Key:   generateHelpDocKey(),
				Name:  name,
				Path:  filePath,
				Group: group,
				Type:  filepath.Ext(filePath),
			})
			m.HelpDocTree = append(m.HelpDocTree, &newGroupDir)
		} else {
			m.HelpDocTree = append(m.HelpDocTree, &helpdoc.HelpDoc{
				Key:   generateHelpDocKey(),
				Name:  name,
				Path:  filePath,
				Group: "default",
				Type:  filepath.Ext(filePath),
			})
		}
	}

	return nil
}

func (m *HelpManager) DeleteHelpDoc(keys []string) error {
	keySet := make(map[string]bool)
	for _, key := range keys {
		keySet[key] = true
	}

	for _, node := range m.HelpDocTree {
		err := traverseHelpDocTree(node, func(d *helpdoc.HelpDoc) error {
			if !d.IsDir {
				_, ok := keySet[d.Key]
				if ok {
					_, err := os.Stat(d.Path)
					if !os.IsNotExist(err) {
						err := os.Remove(d.Path)
						if err != nil {
							return err
						}
					}
					d.Deleted = true
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		_, ok := keySet[node.Key]
		if ok {
			_, err := os.Stat(node.Path)
			if !os.IsNotExist(err) {
				err := os.RemoveAll(node.Path)
				if err != nil {
					return err
				}
			}
			node.Deleted = true
		}
	}
	return nil
}

func traverseHelpDocTree(root *helpdoc.HelpDoc, fn func(node *helpdoc.HelpDoc) error) error {
	if root == nil {
		return nil
	}
	err := fn(root)
	if err != nil {
		return err
	}

	if len(root.Children) == 0 {
		return nil
	}

	for _, child := range root.Children {
		err := traverseHelpDocTree(child, fn)
		if err != nil {
			return err
		}
	}
	return nil
}

type HelpTextVo struct {
	ID          int    `json:"id"`
	Group       string `json:"group"`
	From        string `json:"from"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	PackageName string `json:"packageName"`
	KeyWords    string `json:"keyWords"`
}

type HelpTextVos []HelpTextVo

func (h HelpTextVos) Len() int {
	return len(h)
}

func (h HelpTextVos) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h HelpTextVos) Less(i, j int) bool {
	return h[i].ID < h[j].ID
}

func (m *HelpManager) GetHelpItemPage(pageNum, pageSize int, id, group, from, title string) (int, HelpTextVos) {
	if pageNum <= 0 || pageSize <= 0 {
		return 0, HelpTextVos{}
	}

	// 如果ID不为空
	if id != "" {
		// 加载对应ID的数据
		item, err := m.searchEngine.GetItemByID(id)
		// 若成功
		if err == nil {
			// 返回这条数据
			vo := HelpTextVo{
				Group:       item.Group,
				From:        item.From,
				Title:       item.Title,
				Content:     item.Content,
				PackageName: item.PackageName,
				KeyWords:    item.KeyWords,
			}
			vo.ID, _ = strconv.Atoi(id)
			return 1, HelpTextVos{vo}
		}
		return 0, HelpTextVos{}
	}
	// ID为空的情形，分页查询数据
	total, result, err := m.searchEngine.PaginateDocuments(pageSize, pageNum, group, from, title)
	if err != nil {
		return 0, nil
	}
	var items = make(HelpTextVos, 0)
	for _, item := range result {
		vo := HelpTextVo{
			Group:       item.Group,
			From:        item.From,
			Title:       item.Title,
			Content:     item.Content,
			PackageName: item.PackageName,
			KeyWords:    item.KeyWords,
		}
		vo.ID, _ = strconv.Atoi(id)
		items = append(items, vo)
	}
	return int(total), items
}

// SetDefaultHelpGroup 设置群默认搜索分组
func (group *GroupInfo) SetDefaultHelpGroup(target string) {
	group.DefaultHelpGroup = target
}

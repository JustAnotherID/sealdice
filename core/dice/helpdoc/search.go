package helpdoc

type MatchResult struct {
	ID     string         `json:"id"`
	Fields map[string]any `json:"fields"`
	Score  float64        `json:"score"`
}

type Fields struct {
}

type MatchCollection []*MatchResult

// SearchResult Copied from bleve
type SearchResult struct {
	Hits    MatchCollection
	Total   uint64
	HitBest bool
}

// SearchEngine TODO: 进一步优化结构，封装成通用的搜索
type SearchEngine interface {
	Name() string
	// Init 初始化搜索引擎
	Init() error
	// Close 关闭搜索引擎
	Close()
	// AddItem 添加文档条目，返回添加文档的ID
	AddItem(item HelpTextItem) (string, error)
	// Apply 提交文档条目
	Apply(end bool) error
	// Search 搜索文档条目
	Search(helpPackages []string, text string, titleOnly bool, pageSize, pageNum int, group string) (SearchResult, int, int, int, error)
	// GetHelpTextItemByTermTitle 精确查询title，用于嵌套获取数据的情形
	GetHelpTextItemByTermTitle(title string) (*HelpTextItem, error)
	// GetItemByID 通过ID获取Item数据的方案
	GetItemByID(id string) (*HelpTextItem, error)
	// GetTotalID 获取当前ID总数，注意，ID必须是顺序排列的
	GetTotalID() uint64
	// PaginateDocuments 分页获取数据
	PaginateDocuments(pageSize, pageNum int, group, from, title string) (uint64, []*HelpTextItem, error)
}

type EngineType int

const (
	Bleve       EngineType = iota // 0
	Clover                        // 1
	MeiliSearch                   // 2
)

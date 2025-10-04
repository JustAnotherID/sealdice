package helpdoc

const HelpBuiltinGroup = "builtin"

const (
	Unload int = iota
	Loaded
	LoadError
)

type HelpDoc struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Group      string `json:"group"`
	Type       string `json:"type"`
	IsDir      bool   `json:"isDir"`
	LoadStatus int    `json:"loadStatus"`
	Deleted    bool   `json:"deleted"`

	Children []*HelpDoc `json:"children"`
}

type HelpTextItems []*HelpTextItem

func (e HelpTextItems) String(i int) string {
	return e[i].Title
}

func (e HelpTextItems) Len() int {
	return len(e)
}

const HelpConfigFilename = "help_config.yaml"

type HelpConfig struct {
	Aliases map[string][]string `json:"aliases" yaml:"aliases"`
}

type HelpDocFormat struct {
	Mod     string            `json:"mod"`
	Author  string            `json:"author"`
	Brief   string            `json:"brief"`
	Comment string            `json:"comment"`
	Helpdoc map[string]string `json:"helpdoc"`
}

type HelpTextItem struct {
	Group       string
	From        string
	Title       string
	Content     string
	PackageName string
	// 这俩玩意有用？
	KeyWords   string
	RelatedExt []string
}

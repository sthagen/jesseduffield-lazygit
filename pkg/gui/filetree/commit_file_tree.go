package filetree

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type ICommitFileTree interface {
	ITree[models.CommitFile]

	Get(index int) *CommitFileNode
	GetFile(path string) *models.CommitFile
	GetAllItems() []*CommitFileNode
	GetAllFiles() []*models.CommitFile
	GetRoot() *CommitFileNode
}

type CommitFileTree struct {
	getFiles       func() []*models.CommitFile
	tree           *Node[models.CommitFile]
	showTree       bool
	common         *common.Common
	collapsedPaths *CollapsedPaths
}

func (self *CommitFileTree) CollapseAll() {
	dirPaths := lo.FilterMap(self.GetAllItems(), func(file *CommitFileNode, index int) (string, bool) {
		return file.path, !file.IsFile()
	})

	for _, path := range dirPaths {
		self.collapsedPaths.Collapse(path)
	}
}

func (self *CommitFileTree) ExpandAll() {
	self.collapsedPaths.ExpandAll()
}

var _ ICommitFileTree = &CommitFileTree{}

func NewCommitFileTree(getFiles func() []*models.CommitFile, common *common.Common, showTree bool) *CommitFileTree {
	return &CommitFileTree{
		getFiles:       getFiles,
		common:         common,
		showTree:       showTree,
		collapsedPaths: NewCollapsedPaths(),
	}
}

func (self *CommitFileTree) ExpandToPath(path string) {
	self.collapsedPaths.ExpandToPath(path)
}

func (self *CommitFileTree) ToggleShowTree() {
	self.showTree = !self.showTree
	self.SetTree()
}

func (self *CommitFileTree) Get(index int) *CommitFileNode {
	// need to traverse the three depth first until we get to the index.
	return NewCommitFileNode(self.tree.GetNodeAtIndex(index+1, self.collapsedPaths)) // ignoring root
}

func (self *CommitFileTree) GetIndexForPath(path string) (int, bool) {
	index, found := self.tree.GetIndexForPath(path, self.collapsedPaths)
	return index - 1, found
}

func (self *CommitFileTree) GetAllItems() []*CommitFileNode {
	if self.tree == nil {
		return nil
	}

	// ignoring root
	return lo.Map(self.tree.Flatten(self.collapsedPaths)[1:], func(node *Node[models.CommitFile], _ int) *CommitFileNode {
		return NewCommitFileNode(node)
	})
}

func (self *CommitFileTree) Len() int {
	return self.tree.Size(self.collapsedPaths) - 1 // ignoring root
}

func (self *CommitFileTree) GetItem(index int) types.HasUrn {
	// Unimplemented because we don't yet need to show inlines statuses in commit file views
	return nil
}

func (self *CommitFileTree) GetAllFiles() []*models.CommitFile {
	return self.getFiles()
}

func (self *CommitFileTree) SetTree() {
	showRootItem := self.common.UserConfig().Gui.ShowRootItemInFileTree
	if self.showTree {
		self.tree = BuildTreeFromCommitFiles(self.getFiles(), showRootItem)
	} else {
		self.tree = BuildFlatTreeFromCommitFiles(self.getFiles(), showRootItem)
	}
}

func (self *CommitFileTree) IsCollapsed(path string) bool {
	return self.collapsedPaths.IsCollapsed(path)
}

func (self *CommitFileTree) ToggleCollapsed(path string) {
	self.collapsedPaths.ToggleCollapsed(path)
}

func (self *CommitFileTree) GetRoot() *CommitFileNode {
	return NewCommitFileNode(self.tree)
}

func (self *CommitFileTree) CollapsedPaths() *CollapsedPaths {
	return self.collapsedPaths
}

func (self *CommitFileTree) GetFile(path string) *models.CommitFile {
	for _, file := range self.getFiles() {
		if file.Path == path {
			return file
		}
	}

	return nil
}

func (self *CommitFileTree) InTreeMode() bool {
	return self.showTree
}

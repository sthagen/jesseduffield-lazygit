package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubCommitsController struct {
	baseController
	*ListControllerTrait[*models.Commit]
	c *ControllerCommon
}

var _ types.IController = &SubCommitsController{}

func NewSubCommitsController(
	c *ControllerCommon,
) *SubCommitsController {
	return &SubCommitsController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().SubCommits,
			c.Contexts().SubCommits.GetSelected,
			c.Contexts().SubCommits.GetSelectedItems,
		),
		c: c,
	}
}

func (self *SubCommitsController) Context() types.Context {
	return self.context()
}

func (self *SubCommitsController) context() *context.SubCommitsContext {
	return self.c.Contexts().SubCommits
}

func (self *SubCommitsController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			commit := self.context().GetSelected()
			var task types.UpdateTask
			if commit == nil {
				task = types.NewRenderStringTask("No commits")
			} else {
				refRange := self.context().GetSelectedRefRangeForDiffFiles()
				task = self.c.Helpers().Diff.GetUpdateTaskForRenderingCommitsDiff(commit, refRange)
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title:    "Commit",
					SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
					Task:     task,
				},
			})
		})
	}
}

func (self *SubCommitsController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {
		context := self.context()
		if context.GetSelectedLineIdx() > COMMIT_THRESHOLD && context.GetLimitCommits() {
			context.SetLimitCommits(false)
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SUB_COMMITS}})
		}
	}
}

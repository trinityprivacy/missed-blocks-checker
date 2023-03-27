package events

import (
	configPkg "main/pkg/config"
	"main/pkg/types"
)

type ValidatorGroupChanged struct {
	Validator               *types.Validator
	MissedBlocksBefore      int64
	MissedBlocksAfter       int64
	MissedBlocksGroupBefore *configPkg.MissedBlocksGroup
	MissedBlocksGroupAfter  *configPkg.MissedBlocksGroup
}

func (e *ValidatorGroupChanged) GetDescription() string {
	// increasing
	if e.MissedBlocksGroupBefore.Start < e.MissedBlocksGroupAfter.Start {
		return e.MissedBlocksGroupAfter.DescStart
	}

	// decreasing
	return e.MissedBlocksGroupAfter.DescEnd
}

func (e *ValidatorGroupChanged) GetEmoji() string {
	// increasing
	if e.MissedBlocksGroupBefore.Start < e.MissedBlocksGroupAfter.Start {
		return e.MissedBlocksGroupAfter.EmojiStart
	}

	// decreasing
	return e.MissedBlocksGroupAfter.EmojiEnd
}
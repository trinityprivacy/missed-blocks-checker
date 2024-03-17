package events_test

import (
	"main/pkg/constants"
	"main/pkg/events"
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatorLeftSignatoryBase(t *testing.T) {
	t.Parallel()

	entry := events.ValidatorLeftSignatory{Validator: &types.Validator{Moniker: "test"}}

	assert.Equal(t, constants.EventValidatorLeftSignatory, entry.Type())
	assert.Equal(t, "test", entry.GetValidator().Moniker)
}

func TestValidatorLeftSignatoryFormatHTML(t *testing.T) {
	t.Parallel()

	entry := events.ValidatorLeftSignatory{Validator: &types.Validator{Moniker: "test"}}
	renderData := types.ReportEventRenderData{Notifiers: "notifier1 notifier2", ValidatorLink: "<link>"}
	rendered := entry.Render(constants.FormatTypeHTML, renderData)
	assert.Equal(
		t,
		"<strong>👋 <link> is now not required to sign blocks</strong> notifier1 notifier2",
		rendered,
	)
}

func TestValidatorLeftSignatoryFormatMarkdown(t *testing.T) {
	t.Parallel()

	entry := events.ValidatorLeftSignatory{Validator: &types.Validator{Moniker: "test"}}
	renderData := types.ReportEventRenderData{Notifiers: "notifier1 notifier2", ValidatorLink: "<link>"}
	rendered := entry.Render(constants.FormatTypeMarkdown, renderData)
	assert.Equal(
		t,
		"**👋 <link> is now not required to sign blocks** notifier1 notifier2",
		rendered,
	)
}

func TestValidatorLeftSignatoryFormatUnsupported(t *testing.T) {
	t.Parallel()

	entry := events.ValidatorLeftSignatory{Validator: &types.Validator{Moniker: "test"}}
	renderData := types.ReportEventRenderData{Notifiers: "notifier1 notifier2", ValidatorLink: "<link>"}
	rendered := entry.Render(constants.FormatTypeTest, renderData)
	assert.Equal(
		t,
		"Unsupported format type: test",
		rendered,
	)
}

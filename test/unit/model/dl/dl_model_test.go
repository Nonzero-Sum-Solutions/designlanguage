package dl

import (
	"testing"

	"github.com/Nonzero-Sum-Solutions/designlanguage/model"
	"github.com/Nonzero-Sum-Solutions/designlanguage/test/unit"
	"github.com/stretchr/testify/require"
)

func TestGetBaseComponents(t *testing.T) {
	plane := unit.NewPlane(t)
	rectangle := unit.NewRectangle(t)
	spinner := unit.NewSpinner(t)
	baseComponents := model.GetBaseComponents([]model.Component{plane, rectangle, spinner})

	require.Len(t, baseComponents, 1)
	require.Equal(t, plane, baseComponents[0])
}

func TestGetEntities(t *testing.T) {
	plane := unit.NewPlane(t)
	rectangle := unit.NewRectangle(t)
	spinner := unit.NewSpinner(t)
	entities := model.GetEntities([]model.Component{plane, rectangle, spinner})

	require.Len(t, entities, 1)
	require.Equal(t, rectangle, entities[0])
}

func TestGetObjects(t *testing.T) {
	plane := unit.NewPlane(t)
	spinner := unit.NewSpinner(t)
	rectangle := unit.NewRectangle(t)
	objects := model.GetObjects([]model.Component{plane, rectangle, spinner})

	require.Len(t, objects, 1)
	require.Equal(t, spinner, objects[0])
}

func TestNewDesign(t *testing.T) {
	plane := unit.NewPlane(t)
	spinner := unit.NewSpinner(t)
	rectangle := unit.NewRectangle(t)
	labelledValue := unit.NewLabelledValue(t)
	design := model.NewDesign("", "", "Geometry", []model.Component{plane, rectangle, labelledValue, spinner})

	require.Len(t, design.AllComponents(), 4)
	require.Len(t, design.BaseComponents(), 1)
	require.Len(t, design.Entities(), 2)
	require.Len(t, design.Objects(), 1)
}

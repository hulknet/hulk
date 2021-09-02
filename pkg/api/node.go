package api

import (
	"net/http"

	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/host"
	"github.com/kotfalya/hulk/pkg/node"
	"github.com/labstack/echo/v4"
)

func (r *Rest) nodesList(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, host.NewNodeListModel(r.host))
}

func (r *Rest) nodesItem(ctx echo.Context) error {
	nodeId, err := crypto.FromHex(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	n, ok := r.host.FindNode(nodeId)
	if !ok {
		return ctx.NoContent(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, node.NewItemModel(n))
}

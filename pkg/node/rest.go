package node

type ItemModel struct {
	ID string `json:"id"`
}

func NewItemModel(n *Node) *ItemModel {
	return &ItemModel{
		ID: n.id.Hex(),
	}
}

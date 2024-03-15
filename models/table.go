package models

//enum TableAction {
//	  "add", "remove", "update", "clear", "checkout",

type TableAction string

const (
	AddItem    TableAction = "addItem"
	RemoveItem TableAction = "removeItem"
	Clear      TableAction = "clear"
	Checkout   TableAction = "checkout"
)

type TableMessage struct {
	Action  TableAction `json:"action"`
	ItemID  string      `json:"itemId"`
	Message string      `json:"message"`
}

type TableStateContainer struct {
	TableNumber int         `json:"tableNumber"`
	Open        bool        `json:"open"`
	Items       []OrderItem `json:"items"`
}

//func (t *TableStateContainer) handleAction(action TableAction, itemID string) {
//	switch action {
//	case AddItem:
//		t.addItem(itemID)
//	case RemoveItem:
//		t.removeItem(itemID)
//	case Clear:
//		t.clear()
//	case Checkout:
//		t.checkout()
//	}
//}

//func (t *TableStateContainer) addItem(id string) {
//	t.Items = append(t.Items, OrderItem{ItemID: primitive.ObjectID{id}, Quantity: 1})
//}

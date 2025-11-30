// Package item model
package item

type Item struct {
	Title       string
	Description string
	IsSelected  bool
}

// NewItem creates a new Item instance
func NewItem(title, description string, isSelected bool) Item {
	return Item{Title: title, Description: description, IsSelected: isSelected}
}

// PrintItem returns item description
func (i Item) PrintItem() string {
	return i.Description
}

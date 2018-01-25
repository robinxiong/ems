package sorting

import "github.com/jinzhu/gorm"
type sortingInterface interface {
	GetPosition() int
	SetPosition(int)
}
type sortingDescInterface interface {
	GetPosition() int
	SetPosition(int)
	SortingDesc()
}
// Sorting ascending mode
type Sorting struct {
	Position int `sql:"DEFAULT:NULL"`
}

// GetPosition get current position
func (position Sorting) GetPosition() int {
	return position.Position
}

// SetPosition set position, only set field value, won't save
func (position *Sorting) SetPosition(pos int) {
	position.Position = pos
}

// SortingDESC descending mode
type SortingDESC struct {
	Sorting
}



// SortingDesc make your model sorting desc by default
func (SortingDESC) SortingDesc() {}

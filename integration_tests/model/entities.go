package model

type Author struct {
	ID   string `json:"id" gorm:"primary_key" sql:"type:CHAR(36)"`
	Name string `json:"name" sql:"type:CHAR(36)"`
}

type Book struct {
	ID         string `json:"id" gorm:"primary_key" sql:"type:CHAR(36)"`
	AuthorID   string `json:"authorId" sql:"type:CHAR(36)"`
	BookTypeID string `json:"bookTypeId" sql:"type:CHAR(36)"`
	Status     string `json:"status"`
}

type BookType struct {
	ID   string `json:"id" gorm:"primary_key" sql:"type:CHAR(36)"`
	Name string `json:"name"`
}

type BookCategory struct {
	ID   int    `json:"id" gorm:"primary_key" sql:"type:int(11)"`
	Name string `json:"name"`
}

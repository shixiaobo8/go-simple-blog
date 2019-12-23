package models

import (
	"go-blog/db"
	"time"
)

func init() {
	db.Db.AutoMigrate(&Tags{})
}

type Tags struct {
	TagId       int		`gorm:"primary_key"`
	CreatedAt	time.Time
	UpdateAt 	time.Time
	Name		string  `gorm:"type:varchar(20);not null;default:\"\";unique_index"`
}

func (t *Tags) TableName() string {
	return "tags"
}

func (t *Tags) GetTagByName(name []string) ([]Tags, error)  {
	rows, err := db.Db.Model(t).Select("tag_id, name").Where("name IN (?)", name).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list = make([]Tags, 0)

	for rows.Next() {
		r := Tags{}
		_ = rows.Scan(&r.TagId, &r.Name)
		list = append(list, r)
	}
	return list, nil
}

func (t *Tags) GetTagMapByName(name []string) ([]map[string]interface{}, error){
	rows, err := db.Db.Model(t).Select("tag_id, name").Where("name IN (?)", name).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return RowsToMap(rows), nil
}
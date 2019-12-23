package models

import (
	"fmt"
	"go-blog/db"
	"time"
)

func init() {
	db.Db.AutoMigrate(&PostTags{})
	db.Db.Model(&PostTags{}).AddUniqueIndex("idx_relation", "pid", "tid")
}

type PostTags struct {
	Pid		  int	`gorm:"not null;default:0"`
	Tid       int	`gorm:"not null;default:0"`
	CreatedAt	time.Time
}

func (p *PostTags) TableName() string {
	return "post_tags"
}

// 根据文章id查询标签 结果是map
func (t *PostTags) GetTags(pid int) ([]map[string]interface{}, error) {
	postTagTable := t.TableName()
	tagsTable := (&Tags{}).TableName()

	rows, err := db.Db.Table(postTagTable).Select("tag_id, name").Joins(fmt.Sprintf("left join %s on %s = %s", tagsTable, postTagTable + ".tid", tagsTable + ".tag_id")).Where(fmt.Sprintf("%s.pid = %d", postTagTable, pid)).Rows()
	if err != nil {
		return nil, err
	}

	list := RowsToMap(rows)
	return list, nil
}

// 根据文章id查询标签
func (t *PostTags) GetPostTags(pid uint) []PostTags {
	var list []PostTags
	db.Db.Model(t).Where("pid = ?", pid).Find(&list)
	return list
}
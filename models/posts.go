package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go-blog/db"
	"time"
	"unicode/utf8"
)

func init() {
	db.Db.AutoMigrate(&Posts{})
	r, err := db.Db.Raw("SHOW FULL COLUMNS FROM posts").Rows()
	if err != nil {
		panic(err)
	}
	defer r.Close()
	list := RowsToMap(r)
	for _, v := range list {
		if string(v["Field"].([]uint8)) == "title" && string(v["Collation"].([]uint8)) != "utf8mb4_general_ci" {
			db.Db.Exec("ALTER TABLE `posts` CHANGE `title` `title` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
		}else if string(v["Field"].([]uint8)) == "content" && string(v["Collation"].([]uint8)) != "utf8mb4_general_ci" {
			db.Db.Exec("ALTER TABLE `posts` CHANGE `content` `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
		}else if string(v["Field"].([]uint8)) == "description" && string(v["Collation"].([]uint8)) != "utf8mb4_general_ci" {
			db.Db.Exec("ALTER TABLE `posts` CHANGE `description` `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
		}
	}
}

type Posts struct {
	gorm.Model
	AdminId	    uint
	Title		string `gorm:"type:text;character:utf8mb4"`
	Description string `grom:"type:text;character:utf8mb4"`
	Content     string `gorm:"type:text;character:utf8mb4"`
	Cover		string `gorm:"type:varchar(255);not null;default:\"\""`
	Likes		int	   `gorm:"not null;default:0"`
}

func (p *Posts) TableName() string {
	return "posts"
}

func (p *Posts) Total() int {
	var count int
	db.Db.Model(p).Count(&count)
	return count
}

func (p *Posts) List(fields []string, offset int, limit int) []Posts {
	var posts []Posts
	if fields != nil {
		db.Db.Model(Admin{}).Select(fields).Order("id desc").Offset(offset).Limit(limit).Find(&posts)
	} else {
		db.Db.Model(Admin{}).Offset(offset).Order("id desc").Limit(limit).Find(&posts)
	}

	return posts
}

func (p *Posts) Info(fields []string) error {
	if p.ID == 0 {
		return errors.New("invalid id")
	}
	db.Db.Model(p).Select(fields).First(p)
	return nil
}

/*
处理规则：
忽略图片、换行
1、第一段少于110个字符，截取第一段内容
2、第一段大于110个字符，截取到第110个字符
 */
func (p *Posts) CreateDescription() {
	if utf8.RuneCountInString(p.Content) == 0 {
		return
	}

	sr := []rune(p.Content)
	start := 0
	end := 0
	flag := false
	for i := 0; i < len(sr); i++ {
		if sr[i] == '<' && sr[i+1] == 'p' {
			flag = true
		}

		if flag {
			if start == 0 && sr[i] == '>' {
				start = i + 1
			} else if end <= start && sr[i] == '<' {
				end = i
			}

			if start != 0 && end != 0 && end > start {
				break
			}
		}
	}

	var result []rune
	if start <= end {
		result = sr[start:end]
	} else {
		result = sr[:]
	}

	if len(result) > 110 {
		result = result[:110]
	}
	p.Description = string(result)
}

func (p *Posts) InsertOrUpdate(ts []Tags) error {
	tx := db.Db.Begin()
	// 处理标签
	if len(ts) > 0 {
		for i, item := range ts {
			if item.TagId == 0 {
				item.CreatedAt = time.Now()
				item.UpdateAt = item.CreatedAt
				tx.Create(&item)
				ts[i] = item
			}
		}
	}

	if p.ID == 0 {
		db.Db.Model(p).Create(p)
		for _, tag := range ts {
			var item PostTags
			item.Pid = int(p.ID)
			item.Tid = tag.TagId
			item.CreatedAt = time.Now()
			tx.Create(&item)
		}
	} else {
		// 检查该文章标签关系
		pts := (&PostTags{}).GetPostTags(p.ID)
		for _, tag := range ts {
			exists := false
			var item PostTags
			for k, item := range pts {
				if item.Tid == tag.TagId {
					pts = append(pts[0:k], pts[k+1:]...)
					exists = true
					break
				}
			}
			if exists == false {
				item.Pid = int(p.ID)
				item.Tid = tag.TagId
				item.CreatedAt = time.Now()
				tx.Create(&item)
			}
		}

		if len(pts) > 0 {
			for _, item := range pts {
				db.Db.Where("tid = ? and pid = ?", item.Tid, item.Pid).Delete(PostTags{})
			}
		}

		db.Db.Model(p).Updates(map[string]interface{}{
			"updated_at": time.Now(),
			"title": p.Title,
			"description": p.Description,
			"admin_id": p.AdminId,
			"content": p.Content,
			"cover": p.Cover,
		})
	}

	tx.Commit()
	return nil
}

func (p *Posts) Delete() {
	if p.ID != 0 {
		db.Db.Delete(&p)
	}
}
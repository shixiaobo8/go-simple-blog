package api

import (
	"encoding/json"
	"fmt"
	"go-blog/controllers"
	"go-blog/db"
	"go-blog/models"
	"go-blog/utils"
	"strconv"
	"strings"
)

type PostsController struct {
	controllers.ApiBaseController
}

func (this *PostsController) List() {
	var p models.Posts
	var t models.Tags
	var pt models.PostTags

	var posts []models.Posts
	var count int
	var res = make([]map[string]interface{}, 0)
	var err error

	// 搜索条件
	type search struct {
		Title string
		Tag string
	}

	// 解析数据
	type ob struct {
		Page int
		PageSize int	`json:"page_size"`
		Search search
	}

	var data ob
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, utils.RC["failed"].Msg, nil)
		this.ServeJSON()
	}

	// 分页
	var page = 1
	if data.Page != 0 {
		page = data.Page
	}
	var pageSize = data.PageSize
	var offset = (page - 1) * pageSize

	data.Search.Title = strings.TrimSpace(data.Search.Title)

	// 需要查询字段
	var selectFields = []string{
		"id",
		"admin_id",
		"title",
		"description",
		"likes",
		"cover",
		"created_at",
	}

	if data.Search.Title != "" {
		handler := db.Db.Model(&p).Select(selectFields).Where("title like ?", fmt.Sprintf("%%%s%%", data.Search.Title)).Order("id desc").Offset(offset).Limit(pageSize)
		handler.Count(&count)
		handler.Find(&posts)

	} else if data.Search.Tag != "" {
		postsTable := p.TableName()
		postTagsTable := pt.TableName()

		sf := ""
		for _, item := range selectFields {
			sf += postsTable + "." + item + ","
		}
		sf = strings.TrimRight(sf, ",")

		db.Db.Model(t).Where("name = ?", data.Search.Tag).Find(&t)
		sql := fmt.Sprintf("SELECT %s FROM posts left join post_tags on %s.id = %s.pid WHERE %s.tid = %d and %s.deleted_at is null order by id desc", sf, postsTable, postTagsTable, postTagsTable, t.TagId, postsTable)

		type countRes struct{
			C int
		}
		var countResObj countRes
		db.Db.Raw(sql).Scan(&countResObj)
		count = countResObj.C
		rows, _ := db.Db.Raw(sql + fmt.Sprintf(" limit %d,%d", offset, pageSize)).Rows()
		defer rows.Close()
		for rows.Next() {
			var post models.Posts
			_ = db.Db.ScanRows(rows, &post)
			posts = append(posts, post)
		}
	} else {
		posts = p.List(selectFields, offset, pageSize)
		count = p.Total()
	}

	var temp = &models.Admin{}
	var records = map[int]*models.Admin{}
	for _, post := range posts {
		if _, ok := records[int(post.AdminId)]; ok {
			temp = records[int(post.AdminId)]
		} else {
			temp.ID = post.AdminId
			_ = temp.Info([]string{"nickname", "avatar"})
			records[int(temp.ID)] = temp
		}

		res = append(res, map[string]interface{}{
			"id": post.ID,
			"nickname": temp.Nickname,
			"title" : post.Title,
			"cover": post.Cover,
			"description": post.Description + "...",
			"likes" : post.Likes,
			"created_at" : post.CreatedAt,
		})
		temp = &models.Admin{}
	}
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, &map[string]interface{}{
		"rows" : res,
		"total" : count,
	})
	this.ServeJSON()
}

func (this *PostsController) Recommend() {
	var p models.Posts
	var t models.Tags
	var pt models.PostTags
	var posts []models.Posts
	var err error
	// 解析数据
	type ob struct {
		Size int
		Tags []string
	}

	var data ob
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, utils.RC["failed"].Msg, nil)
		this.ServeJSON()
	}

	if data.Size == 0 {
		data.Size = 5
	}

	// 需要查询字段
	var selectFields = []string{
		"id",
		"title",
		"created_at",
	}

	if len(data.Tags) > 0 {
		postsTable := p.TableName()
		postTagsTable := pt.TableName()

		sf := ""
		for _, item := range selectFields {
			sf += postsTable + "." + item + ","
		}
		sf = strings.TrimRight(sf, ",")

		list, _ := t.GetTagByName(data.Tags)
		var intags string
		for _, value := range list {
			intags += strconv.Itoa(value.TagId) + ","
		}
		intags = strings.TrimRight(intags, ",")

		sql := fmt.Sprintf("SELECT %s as count FROM posts left join post_tags on %s.id = %s.pid WHERE %s.tid IN (%s) and %s.deleted_at is null order by id desc limit %d", sf, postsTable, postTagsTable, postTagsTable, intags, postsTable, data.Size)
		rows, _ := db.Db.Raw(sql).Rows()
		defer rows.Close()
		for rows.Next() {
			var post models.Posts
			_ = db.Db.ScanRows(rows, &post)
			posts = append(posts, post)
		}
	}

	var res = make([]map[string]interface{}, 0)
	for _, post := range posts {
		res = append(res, map[string]interface{}{
			"id": post.ID,
			"title" : post.Title,
			"created_at" : post.CreatedAt,
		})
	}
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, &map[string]interface{}{
		"rows" : res,
	})
	this.ServeJSON()
}

func (this *PostsController) Info() {
	var p models.Posts
	var err error
	var id, _ = this.GetInt("id", 0)

	if id == 0 {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["invalid_args"].Code, utils.RC["invalid_args"].Msg, nil)
		this.ServeJSON()
	}
	p.ID = uint(id)
	_ = p.Info([]string{"id", "title", "content", "cover"})

	// 查询tag
	tagSlice, err := (&models.PostTags{}).GetTags(int(p.ID))

	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, utils.RC["failed"].Msg, nil)
		this.ServeJSON()
	}

	var ts []string
	for _, v := range tagSlice {
		ts = append(ts, string(v["name"].([]uint8)))
	}

	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, &map[string]interface{}{
		"id": p.ID,
		"title": p.Title,
		"content": p.Content,
		"cover": p.Cover,
		"tags": ts,
	})
	this.ServeJSON()
}
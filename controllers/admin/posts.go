package admin

import (
	"encoding/json"
	"fmt"
	"go-blog/controllers"
	"go-blog/db"
	"go-blog/lib"
	"go-blog/models"
	"go-blog/utils"
	"strconv"
	"strings"
)

type PostsController struct {
	controllers.AdminBaseController
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
		Labels []string
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

	// 前端渲染table列名
	columns := []map[string]interface{}{
		{
			"key" : "id",
			"title": "ID",
		},
		{
			"key" : "nickname",
			"title": "用户昵称",
		},
		{
			"key" : "avatar",
			"title": "头像",
		},
		{
			"key" : "title",
			"title" : "标题",
		},
		{
			"key": "likes",
			"title": "点赞数",
		},
		{
			"key" : "created_at",
			"title" : "发表时间",
		},
		{
			"key" : "update_at",
			"title" : "更新时间",
		},
	}

	// 需要查询字段
	var selectFields = []string{
		"id",
		"admin_id",
		"title",
		"likes",
		"created_at",
		"updated_at",
	}

	if data.Search.Title != "" {
		handler := db.Db.Model(&p).Select(selectFields).Where("title like ?", fmt.Sprintf("%%%s%%", data.Search.Title)).Order("id desc").Offset(offset).Limit(pageSize)
		handler.Count(&count)
		handler.Find(&posts)

	} else if len(data.Search.Labels) > 0 {
		postsTable := p.TableName()
		postTagsTable := pt.TableName()

		sf := ""
		for _, item := range selectFields {
			sf += postsTable + "." + item + ","
		}
		sf = strings.TrimRight(sf, ",")

		list, _ := t.GetTagByName(data.Search.Labels)
		var intags string
		for _, value := range list {
			intags += strconv.Itoa(value.TagId) + ","
		}
		intags = strings.TrimRight(intags, ",")

		sql := fmt.Sprintf(" (SELECT %s, count(*) as count FROM posts left join post_tags on %s.id = %s.pid WHERE %s.tid IN (%s) and %s.deleted_at is null GROUP BY %s.id) as tmp where count >= %d order by id desc", sf, postsTable, postTagsTable, postTagsTable, intags, postsTable, postsTable, len(data.Search.Labels))

		type countRes struct{
			C int
		}
		var countResObj countRes
		db.Db.Raw("select count(*) as c from" + sql).Scan(&countResObj)
		count = countResObj.C
		rows, _ := db.Db.Raw( "select * from" + sql + fmt.Sprintf(" limit %d,%d", offset, pageSize)).Rows()
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
	var avatar string
	var records = map[int]*models.Admin{}
	for _, post := range posts {
		if _, ok := records[int(post.AdminId)]; ok {
			temp = records[int(post.AdminId)]
		} else {
			temp.ID = post.AdminId
			_ = temp.Info([]string{"nickname", "avatar"})
			records[int(temp.ID)] = temp
		}

		avatar, _ = lib.Su.FormatPath(temp.Avatar)
		res = append(res, map[string]interface{}{
			"id": post.ID,
			"nickname": temp.Nickname,
			"avatar" : avatar,
			"title" : post.Title,
			"likes" : post.Likes,
			"created_at" : post.CreatedAt,
			"update_at" : post.UpdatedAt,
		})
		temp = &models.Admin{}
	}
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, &map[string]interface{}{
		"columns": columns,
		"rows" : res,
		"total" : count,
	})
	this.ServeJSON()
}

func (this *PostsController) Info() {
	var p models.Posts
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

func (this *PostsController) Update() {
	var err error
	type ob struct {
		models.Posts
		Tags string
	}
	var p ob

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &p)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["invalid_args"].Code, utils.RC["invalid_args"].Msg, nil)
		this.ServeJSON()
	}
	sess := this.GetSession(controllers.AdminSessionName)
	s := sess.(controllers.LoginSession)
	p.AdminId = s.Id

	// 获取tag
	ts := utils.SliceRemoveDupString(strings.Fields(strings.TrimSpace(p.Tags)))
	tss, err := (&models.Tags{}).GetTagByName(ts)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, utils.RC["failed"].Msg, nil)
		this.ServeJSON()
	}
	for _, item := range ts {
		exists := false
		for _, v := range tss {
			if v.Name == item {
				exists = true
				break
			}
		}
		if exists == false {
			tss = append(tss, models.Tags{
				Name:      item,
			})
		}
	}

	// 处理description
	p.CreateDescription()
	err = p.InsertOrUpdate(tss)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, utils.RC["failed"].Msg, nil)
		this.ServeJSON()
	}
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, nil)
	this.ServeJSON()
}

func (this *PostsController) Delete() {
	var m models.Posts
	var err error

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &m)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, utils.RC["failed"].Msg, nil)
		this.ServeJSON()
	}
	if m.ID != 0 {
		m.Delete()
	}
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, nil)
	this.ServeJSON()
}

type FilterType struct {
	Label string
	Value int
}
var Filters = []FilterType{
	{Label:"标题", Value:1},
	{Label:"标签", Value:2},
}

func (this *PostsController) FilterTag() {
	var f = struct {
		Query string
	}{}
	var t models.Tags
	var err error
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &f)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, utils.RC["failed"].Msg, nil)
		this.ServeJSON()
	}
	rows, err := db.Db.Model(&t).Select([]string{"tag_id", "name"}).Where("name like ?", fmt.Sprintf("%%%s%%", f.Query)).Rows()
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, utils.RC["failed"].Msg, nil)
		this.ServeJSON()
	}
	result := models.RowsToMap(rows)
	res := make([]map[string]interface{}, 0)
	for _, item := range result {
		var tmp = map[string]interface{}{
			"label": string(item["name"].([]byte)),
			"value": string(item["name"].([]byte)),
		}
		res = append(res, tmp)
	}
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, &map[string]interface{}{
		"tags": res,
	})
	this.ServeJSON()
}


package order

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
)

type GoodsCate struct {
	Id       int64
	ParentId int64
	Name     string
	Sort     int
	Child    []GoodsCate `orm:"-"`
}

func (this *GoodsCate) TableName() string {
	return models.TableName("goods_cate")
}
func init() {
	orm.RegisterModel(new(GoodsCate))
}
func GetTree() ([]GoodsCate, error) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("goods_cate"))
	var g []GoodsCate
	_, err := qs.OrderBy("-id").All(&g)
	return g, err
}

type NodeTree struct {
	Id       int64
	Name     string
	ParentId int64
	Sibling  map[int64]*NodeTree
}

func NodeToTree(data []GoodsCate) map[int64]*NodeTree {
	var handle map[int64]*NodeTree
	handle = make(map[int64]*NodeTree)
	head := new(NodeTree)
	head.Id = 0
	head.Name = "handle"
	head.ParentId = 0
	head.Sibling = make(map[int64]*NodeTree)
	handle[0] = head
	for _, k := range data {
		if _, ok := handle[k.Id]; ok {
			handle[k.Id].Id = k.Id
			handle[k.Id].ParentId = k.ParentId
			handle[k.Id].Name = k.Name
		} else {
			node := new(NodeTree)
			node.Id = k.Id
			node.ParentId = k.ParentId
			node.Name = k.Name
			node.Sibling = make(map[int64]*NodeTree)
			handle[k.Id] = node
		}
		if _, ok := handle[k.ParentId]; !ok {
			node := new(NodeTree)
			node.Sibling = make(map[int64]*NodeTree)
			node.Id = k.ParentId
			handle[k.ParentId] = node
		}
		handle[k.ParentId].Sibling[k.Id] = handle[k.Id]
	}
	return handle
}
func GetTreeArray(id int64, tree map[int64]*NodeTree) (res []int64) {
	res = append(res, id)
	if len(tree[id].Sibling) > 0 {
		for _, v := range tree[id].Sibling {
			res = append(res, GetTreeArray(v.Id, tree)...)
		}
	}
	return
}

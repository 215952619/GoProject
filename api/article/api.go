package article

import (
	"GoProject/database"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
)

func defaultHandler(c *gin.Context) {
	panic("panic msg !!!")
}

func overflow(_ *gin.Context, _ *database.User) (data interface{}, err error) {
	popularList, err := database.PopularList()
	if err != nil {
		return nil, err
	}
	archiveList, err := database.ArchiveList()
	if err != nil {
		return nil, err
	}
	typeOverflow, err := database.TypeOverflow()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"archive": archiveList,
		"popular": popularList,
		"type":    typeOverflow,
	}, err
}

func list(c *gin.Context, _ *database.User) (data interface{}, err error) {
	var articles []database.Article
	if err = database.DBM.Paginate(&articles, c); err != nil {
		return nil, err
	}
	return articles, nil
}

func detail(c *gin.Context, user *database.User) (data interface{}, err error) {
	id := c.Param("id")
	var article database.Article
	if err = database.DBM.First(&article, "id=?", id); err != nil {
		return nil, err
	}
	err = article.Hit(c, user)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func create(c *gin.Context, user *database.User) (data interface{}, err error) {
	var params CreateArticleParams
	err = c.ShouldBind(&params)
	if err != nil {
		return nil, err
	}

	var article database.Article
	paramJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(paramJson, &article); err != nil {
		return nil, err
	}
	if err = database.DBM.Create(&article); err != nil {
		return nil, err
	}
	return
}

func recommend(c *gin.Context, _ *database.User) (data interface{}, err error) {
	id := c.Param("id")
	var article database.Article
	if err = database.DBM.First(&article, "id=?", id); err != nil {
		return nil, err
	}
	if article.Private == true {
		return nil, errors.New("仅能推荐公开资源")
	}
	if err = database.DBM.Update(&article, "recommend", true); err != nil {
		return nil, err
	}
	return
}

func unrecommended(c *gin.Context, _ *database.User) (data interface{}, err error) {
	id := c.Param("id")
	var article database.Article
	if err = database.DBM.First(&article, "id=?", id); err != nil {
		return nil, err
	}
	if err = database.DBM.Update(&article, "recommend", false); err != nil {
		return nil, err
	}
	return
}

func top(c *gin.Context, _ *database.User) (data interface{}, err error) {
	id := c.Param("id")
	var article database.Article
	if err = database.DBM.First(&article, "id=?", id); err != nil {
		return nil, err
	}
	if article.Private == true {
		return nil, errors.New("仅能置顶公开资源")
	}
	if err = database.DBM.Update(&article, "top", true); err != nil {
		return nil, err
	}
	return
}

func untop(c *gin.Context, _ *database.User) (data interface{}, err error) {
	id := c.Param("id")
	var article database.Article
	if err = database.DBM.First(&article, "id=?", id); err != nil {
		return nil, err
	}
	if err = database.DBM.Update(&article, "top", false); err != nil {
		return nil, err
	}
	return
}

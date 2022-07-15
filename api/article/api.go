package article

import (
	"GoProject/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func defaultHandler(c *gin.Context) {
	panic("panic msg !!!")
}

func overflow(c *gin.Context) (data interface{}, err error) {
	return database.ArchiveList()
}

func list(c *gin.Context) (data interface{}, err error) {
	var articles []database.Article
	if err = database.DBM.Paginate(&articles, c); err != nil {
		return nil, err
	}
	return articles, nil
}

func detail(c *gin.Context) (data interface{}, err error) {
	id := c.Param("id")
	var article database.Article
	if err = database.DBM.First(&article, "id=?", id); err != nil {
		return nil, err
	}
	localLock.Lock()
	defer localLock.Unlock()
	database.DBM.Db.Model(&article).Update("hits", article.Hits+1)
	return article, nil
}

func create(c *gin.Context) (data interface{}, err error) {
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

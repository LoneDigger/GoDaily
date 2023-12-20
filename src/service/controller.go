package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"me.daily/src/bundle"
	"me.daily/src/fuzzy"
	"me.daily/src/log"
	"me.daily/src/token"
	"me.daily/src/util"
)

// @Summary 建立項目
// @Description 建立項目
// @Tags create
// @Accept json
// @Produce json
// @Param Body body bundle.CreateItemRequest true "建立項目"
// @Router /api/item [post]
func (s *Service) createItem(c *gin.Context) {
	var b bundle.CreateItemResponse
	var create bundle.CreateItemRequest

	err := c.BindJSON(&create)
	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		userId := c.GetInt("user_id")
		err := s.d.InsertItem(userId, create.Name, create.SubId, create.Price, create.Remark, create.Date)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "createItem",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 建立主類別
// @Description 建立主類別
// @Tags create
// @Accept json
// @Produce json
// @Param Body body bundle.CreateMainTypeResponse true "建立主類別"
// @Router /api/main [post]
func (s *Service) createMainType(c *gin.Context) {
	var b bundle.CreateMainTypeResponse
	var create bundle.CreateMainTypeRequest

	err := c.BindJSON(&create)
	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		userId := c.GetInt("user_id")
		mainId, err := s.d.InsertMainType(userId, create.Name)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
			b.Name = create.Name
			b.MainId = mainId
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "createItem",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 建立子類別
// @Description 建立子類別
// @Tags create
// @Accept json
// @Produce json
// @Param Body body bundle.CreateSubTypeResponse true "建立子類別"
// @Router /api/sub [post]
func (s *Service) createSubType(c *gin.Context) {
	var b bundle.CreateSubTypeResponse
	var create bundle.CreateSubTypeRequest

	err := c.BindJSON(&create)
	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		userId := c.GetInt("user_id")
		subId, err := s.d.InsertSubType(userId, create.MainId, create.Name, create.Increase)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
			b.Name = create.Name
			b.SubId = subId
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "createSubType",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 建立使用者
// @Description 建立使用者
// @Tags create
// @Accept json
// @Produce json
// @Param Body body bundle.CreateUserRequest true "建立使用者"
// @Router /api/user [post]
func (s *Service) createUser(c *gin.Context) {
	var b bundle.ErrorResponse
	var create bundle.CreateUserRequest

	err := c.BindJSON(&create)
	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		pw := util.HashPassword(create.Password)
		userId, err := s.d.CreateUser(create.Username, pw)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "createUser",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 刪除項目
// @Description 刪除項目
// @Tags delete
// @Param item_id path int true "項目編號"
// @Accept json
// @Produce json
// @Router /api/item/{item_id} [delete]
func (s *Service) deleteItem(c *gin.Context) {
	var b bundle.DeleteMainTypeResponse
	userId := c.GetInt("user_id")
	itemId, err := strconv.Atoi(c.Param("item_id"))

	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		err := s.d.DeleteItem(userId, itemId)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "deleteItem",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 刪除主類別名稱
// @Description 刪除主類別名稱
// @Tags delete
// @Param main_id path int true "主類別編號"
// @Accept json
// @Produce json
// @Router /api/main/{main_id} [delete]
func (s *Service) deleteMainType(c *gin.Context) {
	var b bundle.DeleteMainTypeResponse
	userId := c.GetInt("user_id")
	mainId, err := strconv.Atoi(c.Param("main_id"))

	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		err := s.d.DeleteMainType(userId, mainId)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "deleteMainType",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 刪除子類別名稱
// @Description 刪除子類別名稱
// @Tags delete
// @Param sub_id path int true "子類別編號"
// @Accept json
// @Produce json
// @Router /api/sub/{sub_id} [delete]
func (s *Service) deleteSubType(c *gin.Context) {
	var b bundle.DeleteMainTypeResponse
	userId := c.GetInt("user_id")
	subId, err := strconv.Atoi(c.Param("sub_id"))

	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		err := s.d.DeleteSubType(userId, subId)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "deleteSubType",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 登入
// @Description 登入
// @Tags login
// @Accept json
// @Produce json
// @Param Body body bundle.LoginRequest true "登入"
// @Router /api/login [post]
func (s *Service) login(c *gin.Context) {
	var b bundle.ErrorResponse
	var login bundle.LoginRequest

	err := c.BindJSON(&login)
	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		userId, pw, err := s.d.Login(login.Username)

		if err != nil {
			b.Code = err.Error()
		} else {
			// 比對密碼
			if !util.CheckPasswordHash(login.Password, pw) {
				b.Code = bundle.CodePassword
			} else {
				auth := token.NewToken(userId, login.Username, expiredTime*time.Second)
				c.SetCookie("Authorization", auth, expiredTime, "/", c.Request.Host, false, false)
				b.Code = bundle.CodeOk

				s.c.Add(strconv.Itoa(userId), nil, cache.DefaultExpiration)
			}
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "login",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 登出
// @Description 登出
// @Tags get
// @Router /api/logout [get]
func (s *Service) logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", c.Request.Host, false, false)

	c.Set("code", bundle.CodeOk)
	c.JSON(http.StatusOK, bundle.ErrorResponse{
		Code: bundle.CodeOk,
	})
}

// @Summary 取得全部類別
// @Description 取得全部類別
// @Tags get
// @Accept json
// @Produce json
// @Router /api/all [get]
func (s *Service) getAll(c *gin.Context) {
	var b bundle.GetAllTypeResponse
	userId := c.GetInt("user_id")

	all, err := s.d.GetAllType(userId)

	if err != nil {
		b.Code = err.Error()
	} else {
		b.Code = bundle.CodeOk
		b.List = all
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 取得單一項目
// @Description 取得單一項目
// @Tags get
// @Param item_id path int true "項目編號"
// @Accept json
// @Produce json
// @Router /api/item/{item_id} [get]
func (s *Service) getItem(c *gin.Context) {
	var b bundle.GetItemResponse

	userId := c.GetInt("user_id")
	itemId, err := strconv.Atoi(c.Param("item_id"))

	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		item, err := s.d.GetItem(userId, itemId)
		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
			b.Item = item
		}
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 由日期取得預覽項目
// @Description 由日期取得預覽項目
// @Tags get
// @Param start		query string true "起始日期"
// @Param end		query string true "結束日期"
// @Param content	query string false "關鍵字"
// @Accept json
// @Produce json
// @Router /api/items [get]
func (s *Service) getItems(c *gin.Context) {
	var b bundle.GetItemsResponse
	b.List = make([]bundle.PreviewItem, 0)

	userId := c.GetInt("user_id")
	startStr := c.Query("start")
	startDate, err := time.Parse(dateFormat, startStr)
	if err != nil {
		b.Code = bundle.CodeFormat
		c.JSON(http.StatusOK, b)
		return
	}

	endStr := c.Query("end")
	endDate, err := time.Parse(dateFormat, endStr)
	if err != nil {
		b.Code = bundle.CodeFormat
		c.Set("code", b.Code)
		c.JSON(http.StatusOK, b)
		return
	}

	// 日期錯誤
	if startDate.After(endDate) {
		b.Code = bundle.CodeDate
		c.Set("code", b.Code)
		c.JSON(http.StatusOK, b)
		return
	}

	// 查詢區間
	shiftTime := startDate.AddDate(dateRange, 0, 0)
	if endDate.After(shiftTime) {
		b.Code = bundle.CodeDate
		c.Set("code", b.Code)
		c.JSON(http.StatusOK, b)
		return
	}

	var items []bundle.PreviewItem
	content := c.Query("content")
	if len(content) == 0 {
		items, err = s.d.GetPerviewItemsByDate(userId, startStr, endStr)
	} else {
		items, err = s.d.LikeName(userId, content, startStr, endStr)
	}

	if err != nil {
		b.Code = err.Error()
	} else {
		b.Code = bundle.CodeOk
		b.List = items
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

func (s *Service) getLog(c *gin.Context) {
	c.String(http.StatusOK, log.LogHistory.String())
}

// @Summary 取得全部主類別
// @Description 取得全部主類別
// @Tags get
// @Accept json
// @Produce json
// @Router /api/main [get]
func (s *Service) getMainType(c *gin.Context) {
	userId := c.GetInt("user_id")
	main, err := s.d.GetMainType(userId)

	var b bundle.GetMainTypeResponse
	if err != nil {
		b.Code = err.Error()
	} else {
		b.Code = bundle.CodeOk
		b.List = main
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 取得這個月的主類別總和
// @Description 取得這個月的主類別總和
// @Tags get
// @Accept json
// @Produce json
// @Router /api/sum/main [get]
func (s *Service) getSumByMainType(c *gin.Context) {
	var b bundle.GetSumByMainTypeResponse
	userId := c.GetInt("user_id")

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)
	m, err := s.d.GetSumByMainType(userId, startDate.Format(dateFormat), endDate.Format(dateFormat))

	if err != nil {
		b.Code = err.Error()
	} else {
		b.Code = bundle.CodeOk
		b.List = m
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 取得前幾個月收支總和
// @Description 取得前幾個月收支總和
// @Tags get
// @Accept json
// @Produce json
// @Router /api/spend/month/{count} [get]
func (s *Service) getSpendByLastMonthly(c *gin.Context) {
	var b bundle.GetSpendByMonthlyResponse
	userId := c.GetInt("user_id")

	count, err := strconv.Atoi(c.Param("count"))
	if count > 12 {
		count = 12
	}

	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		now := time.Now()
		tmpDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		// 月底
		endDate := tmpDate.AddDate(0, 1, -1)
		// 月初
		start := tmpDate.AddDate(0, -count, 0)
		m, err := s.d.GetSumByMonth(userId, start.Format(dateFormat), endDate.Format(dateFormat))

		l := len(m)
		if l != count {
			list := make([]bundle.Monthly, count)
			index := 0
			for i := 0; i < count; i++ {
				list[i].Date = start.AddDate(0, i, 0)
				if index < l && list[i].Date.Equal(m[index].Date) {
					list[i].Sum = m[index].Sum
					index++
				}
			}

			m = list
		}

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
			b.List = m
		}
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 取得全部子類別
// @Description 取得全部子類別
// @Tags get
// @Param main_id path int true "主類別編號"
// @Accept json
// @Produce json
// @Router /api/sub/{main_id} [get]
func (s *Service) getSubType(c *gin.Context) {
	var b bundle.GetSubTypeResponse
	userId := c.GetInt("user_id")
	mainId, err := strconv.Atoi(c.Param("main_id"))

	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		sub, err := s.d.GetSubType(userId, mainId)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
			b.List = sub
		}
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 模糊搜尋名稱
// @Description 模糊搜尋名稱
// @Tags get
// @Param start		query string true "起始日期"
// @Param end		query string true "結束日期"
// @Param content	query string true "關鍵字"
// @Accept json
// @Produce json
// @Router /api/search/name [get]
func (s *Service) searchByName(c *gin.Context) {
	var b bundle.GetItemsResponse
	b.List = make([]bundle.PreviewItem, 0)

	userId := c.GetInt("user_id")
	startStr := c.Query("start")
	startDate, err := time.Parse(dateFormat, startStr)
	if err != nil {
		b.Code = bundle.CodeFormat
		c.JSON(http.StatusOK, b)
		return
	}

	endStr := c.Query("end")
	endDate, err := time.Parse(dateFormat, endStr)
	if err != nil {
		b.Code = bundle.CodeFormat
		c.Set("code", b.Code)
		c.JSON(http.StatusOK, b)
		return
	}

	// 日期錯誤
	if startDate.After(endDate) {
		b.Code = bundle.CodeDate
		c.Set("code", b.Code)
		c.JSON(http.StatusOK, b)
		return
	}

	// 查詢區間
	shiftTime := startDate.AddDate(dateRange, 0, 0)
	if endDate.After(shiftTime) {
		b.Code = bundle.CodeDate
		c.Set("code", b.Code)
		c.JSON(http.StatusOK, b)
		return
	}

	var items []bundle.PreviewItem
	content := c.Query("content")
	if len(content) == 0 {
		b.Code = bundle.CodeEmptyContent
		c.Set("code", b.Code)
		c.JSON(http.StatusOK, b)
	} else {
		items, err = s.d.GetPerviewItemsByDate(userId, startStr, endStr)
		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk

			// 過濾內文
			for _, item := range items {
				if fuzzy.FuzzySearch(item.Name, content, nil) {
					b.List = append(b.List, item)
				}
			}
		}
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 修改項目
// @Description 修改項目
// @Tags update
// @Accept json
// @Produce json
// @Param Body body bundle.UpdateItemRequest true "修改"
// @Router /api/item [put]
func (s *Service) updateItem(c *gin.Context) {
	var b bundle.UpdateItemResponse
	var update bundle.UpdateItemRequest

	err := c.BindJSON(&update)
	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		userId := c.GetInt("user_id")
		err := s.d.UpdateItem(userId, update.ItemId, update.Name, update.SubId, update.Price, update.Remark, update.Date)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "updateItem",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 修改主類別名稱
// @Description 修改主類別名稱
// @Tags update
// @Accept json
// @Produce json
// @Param Body body bundle.UpdateMainTypeRequest true "修改"
// @Router /api/main [put]
func (s *Service) updateMainType(c *gin.Context) {
	var b bundle.UpdateMainTypeResponse
	var update bundle.UpdateMainTypeRequest

	err := c.BindJSON(&update)
	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		userId := c.GetInt("user_id")
		err := s.d.UpdateMainType(userId, update.MainId, update.Name)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "updateMainType",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

// @Summary 修改子類別名稱
// @Description 修改子類別名稱
// @Tags update
// @Accept json
// @Produce json
// @Param Body body bundle.UpdateSubTypeRequest true "修改"
// @Router /api/sub [put]
func (s *Service) updateSubType(c *gin.Context) {
	var b bundle.UpdateSubTypeResponse
	var update bundle.UpdateSubTypeRequest

	err := c.BindJSON(&update)
	if err != nil {
		b.Code = bundle.CodeFormat
	} else {
		userId := c.GetInt("user_id")
		err := s.d.UpdateSubType(userId, update.SubId, update.Name, update.Increase)

		if err != nil {
			b.Code = err.Error()
		} else {
			b.Code = bundle.CodeOk
		}

		log.LogHistory.L.WithFields(logrus.Fields{
			"Method": "updateSubType",
			"UserId": userId,
			"Code":   b.Code,
		}).Info("Api")
	}

	c.Set("code", b.Code)
	c.JSON(http.StatusOK, b)
}

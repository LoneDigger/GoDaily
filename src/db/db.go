package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"me.daily/src/bundle"
)

var initMain []string
var initSub [][]string

func init() {
	initMain = []string{"收入", "餐費", "交通", "裝扮", "房租", "日用", "娛樂", "其他"}

	initSub = make([][]string, len(initMain))
	initSub[0] = append(initSub[0], "薪水", "租借", "股息", "基金", "接案", "販售")
	initSub[1] = append(initSub[1], "早餐", "午餐", "晚餐", "消夜", "飲料", "食物", "咖啡", "水")
	initSub[2] = append(initSub[2], "機票", "高鐵", "火車", "捷運", "輕軌", "客運", "公車", "船票", "計程車", "租車", "共享單車", "汽油", "停車費", "修車")
	initSub[3] = append(initSub[3], "衣著", "鞋子", "包包", "剪髮", "洗衣", "化妝品", "保養品")
	initSub[4] = append(initSub[4], "房租", "水費", "電費", "管理費")
	initSub[5] = append(initSub[5], "電話費", "醫療", "購物", "物品")
	initSub[6] = append(initSub[6], "網咖", "遊戲", "KTV", "電影", "書", "門票", "住宿")
	initSub[7] = append(initSub[7], "其他", "遺失記憶")
}

type Db struct {
	db *sqlx.DB
}

func NewDb(host, user, password, dbname string) *Db {
	connect := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, 5432, user, password, dbname)

	db, err := sqlx.Connect("postgres", connect)
	if err != nil {
		panic(err)
	}

	return &Db{
		db: db,
	}
}

// 確認子類別持有者
func (d *Db) checkSub(userId, subId int) error {
	s := `SELECT COUNT(1) FROM sub_types
			WHERE user_id=$1 AND id=$2 AND NOT deleted`
	count := 0
	err := d.db.QueryRow(s, userId, subId).Scan(&count)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	if count == 1 {
		return nil
	}

	return errors.New(bundle.CodeHold)
}

// 確認主類別名稱有無重複
func (d *Db) checkMainTypeName(userId int, name string) error {
	s := `SELECT COUNT(1) 
			FROM main_types
			WHERE user_id=$1 AND name=$2 AND NOT deleted`
	count := 0
	err := d.db.QueryRow(s, userId, name).Scan(&count)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	if count == 0 {
		return nil
	}

	return errors.New(bundle.CodeTypeRepeat)
}

// 確認主類別名持有者
func (d *Db) checkMainTypeOwner(userId, mainId int) error {
	s := `SELECT COUNT(1) 
			FROM main_types 
			WHERE user_id=$1 AND id=$2 AND NOT deleted`
	count := 0

	err := d.db.QueryRow(s, userId, mainId).Scan(&count)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	if count == 1 {
		return nil
	}

	return errors.New(bundle.CodeHold)
}

// 確認子類別名稱有無重複
func (d *Db) checkSubTypeName(userId, mainId int, name string) error {
	s := `SELECT COUNT(1)
			FROM sub_types
            WHERE user_id=$1 AND name=$2 AND main_id=$3 AND NOT deleted`
	count := 0
	err := d.db.QueryRow(s, userId, name, mainId).Scan(&count)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	if count == 0 {
		return nil
	}

	return errors.New(bundle.CodeTypeRepeat)
}

// 建立主類別
func createMain(tx *sqlx.Tx, userId int) ([]int, error) {
	var ids []int
	s := `INSERT INTO main_types (user_id, name) 
			VALUES ($1, $2)
			RETURNING id`

	id := 0
	for _, n := range initMain {
		err := tx.QueryRow(s, userId, n).Scan(&id)
		if err != nil {
			return nil, errors.New(bundle.CodeDb)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// 建立子類別
func createSub(tx *sqlx.Tx, mainIds []int) error {
	type tmp struct {
		UserId   int `db:"user_id"`
		MainId   int `db:"main_id"`
		Name     string
		Increase int
	}

	for i, sub := range initSub {
		increase := -1
		if i == 0 {
			increase = 1
		}

		m := make([]tmp, len(sub))

		for p, n := range sub {
			m[p].MainId = mainIds[i]
			m[p].Name = n
			m[p].Increase = increase
		}

		s := `INSERT INTO sub_types (user_id, main_id, name, increase) 
				VALUES (:user_id, :main_id, :name, :increase)`
		_, err := tx.NamedExec(s, m)
		if err != nil {
			return errors.New(bundle.CodeDb)
		}
	}

	return nil
}

// 新增使用者
func (d *Db) CreateUser(username, password string) (int, error) {
	var count int
	s := `SELECT COUNT(*) FROM users WHERE username=$1`
	err := d.db.QueryRow(s, username).Scan(&count)
	if err != nil {
		return -1, errors.New(bundle.CodeDb)
	}

	tx := d.db.MustBegin()
	defer tx.Rollback()

	if count == 0 {
		var userId int
		s = `INSERT INTO users (username, password) 
				VALUES ($1, $2)
				RETURNING id`
		err = tx.QueryRow(s, username, password).Scan(&userId)
		if err != nil {
			return -1, errors.New(bundle.CodeDb)
		}

		mainIds, err := createMain(tx, userId)
		if err != nil {
			return -1, err
		}

		if err = createSub(tx, mainIds); err != nil {
			return -1, err
		}

		tx.Commit()
		return userId, nil
	}

	return -1, errors.New(bundle.CodeUserRepeat)
}

// 刪除帳單項目
func (d *Db) DeleteItem(userId, id int) error {
	s := `DELETE FROM bills WHERE user_id=$1 AND id=$2`
	r, err := d.db.Exec(s, userId, id)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	row, _ := r.RowsAffected()

	if row == 0 {
		return errors.New(bundle.CodeNoData)
	}

	return nil
}

// 刪除主類型
func (d *Db) DeleteMainType(userId, id int) error {
	s := `UPDATE main_types SET deleted=true WHERE user_id=$1 AND id=$2 AND NOT deleted`
	r, err := d.db.Exec(s, userId, id)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	row, _ := r.RowsAffected()

	if row == 0 {
		return errors.New(bundle.CodeNoData)
	}

	return nil
}

// 刪除子類型
func (d *Db) DeleteSubType(userId, id int) (err error) {
	s := `UPDATE sub_types SET deleted=true WHERE user_id=$1 AND id=$2 AND NOT deleted`
	r, err := d.db.Exec(s, userId, id)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	row, _ := r.RowsAffected()

	if row == 0 {
		return errors.New(bundle.CodeNoData)
	}

	return nil
}

// 取得全部類別
func (d *Db) GetAllType(userId int) ([]bundle.AllType, error) {
	arr := make([]struct {
		MainId   int    `db:"main_id"`
		MainName string `db:"main_name"`
		SubId    int    `db:"sub_id"`
		SubName  string `db:"sub_name"`
		Increase bool   `db:"increase"`
	}, 0)

	s := `SELECT m.id AS main_id, m.name AS main_name, s.id AS sub_id, 
				s.name AS sub_name, s.increase > 0 AS increase
			FROM main_types AS m
			LEFT JOIN sub_types AS s
			ON m.id=s.main_id
			WHERE m.user_id=$1 AND s.user_id=$1 
				AND NOT m.deleted AND NOT s.deleted
			ORDER BY main_id, sub_id`
	err := d.db.Select(&arr, s, userId)
	if err != nil {
		err = errors.New(bundle.CodeDb)
	}

	ats := make([]bundle.AllType, 0)
	var at bundle.AllType
	tempId := -1
	for _, a := range arr {
		if tempId != a.MainId {

			if len(at.Subs) != 0 {
				ats = append(ats, at)
			}

			at.Id = a.MainId
			at.Name = a.MainName
			at.Subs = make([]bundle.Sub, 0)
			tempId = a.MainId
		}

		at.Subs = append(at.Subs, bundle.Sub{
			Id:       a.SubId,
			Name:     a.SubName,
			Increase: a.Increase,
		})
	}

	ats = append(ats, at)

	return ats, err
}

// 取得單項目
func (d *Db) GetItem(userId, itemId int) (bundle.Item, error) {
	var item bundle.Item

	s := `SELECT b.id, b.name, main_id, sub_id, price, remark, date
			FROM bills AS b
			LEFT JOIN sub_types AS s
			ON s.id=b.sub_id
			WHERE b.user_id=$1 AND b.id=$2`

	err := d.db.Get(&item, s, userId, itemId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New(bundle.CodeNoData)
		} else {
			err = errors.New(bundle.CodeDb)
		}
	}

	return item, err
}

// 用日期取得預覽項目
func (d *Db) GetPerviewItemsByDate(userId int, start, end string) ([]bundle.PreviewItem, error) {
	items := make([]bundle.PreviewItem, 0)

	s := `SELECT b.id, m.id AS "main_id", m.name AS "main_name", 
				s.id AS "sub_id", s.name AS "sub_name", b.name, 
				b.price, s.increase, TO_CHAR(b.date, 'yyyy-mm-dd') AS "date"
			FROM bills AS b
			LEFT JOIN sub_types AS s
			ON b.sub_id=s.id
			LEFT JOIN main_types AS m
			ON m.id=s.main_id
			WHERE b.user_id=$1 AND b.date BETWEEN $2 AND $3
			ORDER BY b.date, b.id`

	err := d.db.Select(&items, s, userId, start, end)
	if err != nil {
		err = errors.New(bundle.CodeDb)
	}

	return items, err
}

// 取得主類別
func (d *Db) GetMainType(userId int) ([]bundle.Main, error) {
	arr := make([]bundle.Main, 0)
	s := `SELECT id, name 
			FROM main_types 
			WHERE NOT deleted AND user_id=$1
			ORDER BY id`
	err := d.db.Select(&arr, s, userId)
	if err != nil {
		err = errors.New(bundle.CodeDb)
	}

	return arr, err
}

// 取得子類別
func (d *Db) GetSubType(userId, main_id int) ([]bundle.Sub, error) {
	arr := make([]bundle.Sub, 0)
	s := `SELECT id, name, increase>0 AS increase
			FROM sub_types
			WHERE NOT deleted AND user_id=$1 AND main_id=$2 
			ORDER BY id`
	err := d.db.Select(&arr, s, userId, main_id)
	if err != nil {
		err = errors.New(bundle.CodeDb)
	}

	return arr, err
}

// 取得主類別總和，排除收入
func (d *Db) GetSumByMainType(userId int, start, end string) ([]bundle.MainSumMonthly, error) {
	arr := make([]bundle.MainSumMonthly, 0)

	s := `SELECT SUM(b.price) AS sum, m.name
			FROM bills AS b
			LEFT JOIN sub_types AS s
			ON b.sub_id=s.id
			LEFT JOIN main_types AS m
			ON m.id=s.main_id
			WHERE b.user_id=$1 AND b.date BETWEEN $2 AND $3 AND s.increase<0
			GROUP BY m.id`

	err := d.db.Select(&arr, s, userId, start, end)
	if err != nil {
		err = errors.New(bundle.CodeDb)
	}

	return arr, err
}

// 取得月結總金額
func (d *Db) GetSumByMonth(userId int, start, end string) ([]bundle.Monthly, error) {
	items := make([]bundle.Monthly, 0)

	s := `SELECT SUM(s.increase * b.price) AS "sum", DATE_TRUNC('month', b.date) AS "date"
			FROM bills AS b
			LEFT JOIN sub_types AS s
			ON b.sub_id=s.id
			LEFT JOIN main_types AS m
			ON m.id=s.main_id
			WHERE b.user_id=$1 AND b.date BETWEEN $2 AND $3
			GROUP BY DATE_TRUNC('month', b.date)
			ORDER BY date`

	err := d.db.Select(&items, s, userId, start, end)
	if err != nil {
		err = errors.New(bundle.CodeDb)
	}

	return items, err
}

// 新增帳單項目
func (d *Db) InsertItem(userId int, name string, subId int, price int, remark, date string) error {
	err := d.checkSub(userId, subId)
	if err != nil {
		return err
	}

	s := `INSERT INTO bills (user_id, name, sub_id, price, remark, date) 
			VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = d.db.Exec(s, userId, name, subId, price, remark, date)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	return nil
}

// 新增主類型
func (d *Db) InsertMainType(userId int, name string) (int, error) {
	err := d.checkMainTypeName(userId, name)
	if err != nil {
		return 0, err
	}

	var id int
	s := `INSERT INTO main_types (user_id, name) 
			VALUES ($1, $2) RETURNING id`
	err = d.db.QueryRow(s, userId, name).Scan(&id)
	if err != nil {
		return 0, errors.New(bundle.CodeDb)
	}

	return id, nil
}

// 新增子類型
func (d *Db) InsertSubType(userId, mainId int, subName string, increase bool) (int, error) {
	err := d.checkMainTypeOwner(userId, mainId)
	if err != nil {
		return 0, err
	}

	err = d.checkSubTypeName(userId, mainId, subName)
	if err != nil {
		return 0, err
	}

	i := -1
	if increase {
		i = 1
	}

	var id int
	s := `INSERT INTO sub_types (name, user_id, main_id, increase) 
			VALUES ($1, $2, $3, $4)
			RETURNING id`
	err = d.db.QueryRow(s, subName, userId, mainId, i).Scan(&id)
	if err != nil {
		return 0, errors.New(bundle.CodeDb)
	}

	return id, nil
}

// 模糊搜尋名稱
func (d *Db) LikeName(userId int, keyword, start, end string) ([]bundle.PreviewItem, error) {
	items := make([]bundle.PreviewItem, 0)

	s := `SELECT b.id, m.id AS "main_id", m.name AS "main_name", 
				s.id AS "sub_id", s.name AS "sub_name", b.name, 
				b.price, s.increase, TO_CHAR(b.date, 'yyyy-mm-dd') AS "date"
			FROM bills AS b
			LEFT JOIN sub_types AS s
			ON b.sub_id=s.id
			LEFT JOIN main_types AS m
			ON m.id=s.main_id
			WHERE b.user_id=$1 AND b.date BETWEEN $2 AND $3 
				AND b.name LIKE $4 
			ORDER BY b.date, b.id`

	err := d.db.Select(&items, s, userId, start, end, keyword)
	if err != nil {
		err = errors.New(bundle.CodeDb)
	}

	return items, err
}

// 登入
func (d *Db) Login(username string) (int, string, error) {
	login := struct {
		Id       int
		Password string
	}{}

	s := `SELECT id, password FROM users WHERE username=$1`
	err := d.db.QueryRowx(s, username).StructScan(&login)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", errors.New(bundle.CodeUsername)
		}
		return 0, "", errors.New(bundle.CodeDb)
	}

	return login.Id, login.Password, err
}

// 更新帳單項目
func (d *Db) UpdateItem(userId, itemId int, name string, subId int, price int, remark, date string) error {
	err := d.checkSub(userId, subId)
	if err != nil {
		return err
	}

	s := `UPDATE bills SET name=$1, sub_id=$2, price=$3, remark=$4, date=$7
			WHERE user_id=$5 AND id=$6`
	r, err := d.db.Exec(s, name, subId, price, remark, userId, itemId, date)
	if err != nil {
		return errors.New(bundle.CodeHold)
	}

	row, _ := r.RowsAffected()
	if row == 0 {
		return errors.New(bundle.CodeNoData)
	}

	return nil
}

// 更新主類型
func (d *Db) UpdateMainType(userId, id int, name string) error {
	err := d.checkMainTypeName(userId, name)
	if err != nil {
		return err
	}

	s := `UPDATE main_types 
			SET name=$1
			WHERE id=$2 AND user_id=$3`
	r, err := d.db.Exec(s, name, id, userId)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	row, _ := r.RowsAffected()

	if row == 0 {
		return errors.New(bundle.CodeNoData)
	}

	return nil
}

// 更新子類型
func (d *Db) UpdateSubType(userId, subId int, name string, increase bool) error {
	s := `SELECT main_id 
			FROM sub_types 
			WHERE user_id=$1 AND sub_id=$2`

	var mainId int
	err := d.db.QueryRow(s,
		userId, subId).Scan(&mainId)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	err = d.checkSubTypeName(userId, mainId, name)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	i := -1
	if increase {
		i = 1
	}

	s = `UPDATE sub_types
			SET name=$1, increase=$2
			WHERE id=$3 AND user_id=$4`
	r, err := d.db.Exec(s, name, i, subId, userId)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	row, _ := r.RowsAffected()

	if row == 0 {
		return errors.New(bundle.CodeNoData)
	}

	return nil
}

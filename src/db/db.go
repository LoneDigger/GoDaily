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
	connect := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, 5432, user, password, dbname)

	db, err := sqlx.Connect("postgres", connect)
	if err != nil {
		panic(err)
	}

	return &Db{
		db: db,
	}
}

// 確認子類別持有者
func (d *Db) checkSub(userId, subId int) (bool, error) {
	s := `SELECT COUNT(1) FROM main_types AS m
			LEFT JOIN sub_types AS s
			ON s.main_id=m.id
			WHERE user_id=$1 AND s.id=$2`
	count := 0
	err := d.db.QueryRow(s, userId, subId).Scan(&count)
	if err != nil {
		return false, errors.New(bundle.CodeDb)
	}

	return count == 1, nil
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

		s := `INSERT INTO sub_types (main_id, name, increase) 
				VALUES (:main_id, :name, :increase)`
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
				VALUES ($1, $2) RETURNING id`
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
func (d *Db) DeleteItem(userId, id int) (err error) {
	s := `DELETE FROM bills where id=$1 AND user=$2`
	_, err = d.db.Exec(s, id, userId)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	return
}

// 刪除主類型
func (d *Db) DeleteMainType(userId, id int) (err error) {
	s := `UPDATE main_types SET deleted=true WHERE id=$1 AND user_id=$2`
	_, err = d.db.Exec(s, id, userId)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	return
}

// 刪除子類型
func (d *Db) DeleteSubType(userId, id int) (err error) {
	s := `UPDATE sub_types AS s SET deleted=true 
			FROM main_types AS m 
			WHERE s.main_id=m.id AND s.id=$1 AND m.user_id=$2`
	_, err = d.db.Exec(s, id, userId)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	return
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
			WHERE NOT m.deleted AND NOT s.deleted AND m.user_id=$1
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
			WHERE user_id=$1 AND b.id=$2`

	err := d.db.Get(&item, s, userId, itemId)
	if err != nil {
		err = errors.New(bundle.CodeDb)
	}

	return item, err
}

// 用日期取得預覽項目
func (d *Db) GetPerviewItemsByDate(userId, limit, offset int, start, end string) ([]bundle.PreviewItem, error) {
	var items []bundle.PreviewItem

	s := `SELECT b.id, m.name AS "main_name", s.name AS "sub_name", b.name, 
				b.price, s.increase, TO_CHAR(b.date, 'yyyy-mm-dd') AS "date"
			FROM bills AS b
			LEFT JOIN sub_types AS s
			ON b.sub_id=s.id
			LEFT JOIN main_types AS m
			ON m.id=s.main_id
			WHERE b.user_id=$1 AND b.date BETWEEN $2 AND $3
			ORDER BY b.date, b.id`

	//err := d.db.Select(&items, s, userId, start, end, limit, offset)
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
			WHERE NOT deleted AND user_id=$1 ORDER BY id`
	err := d.db.Select(&arr, s, userId)
	if err != nil {
		err = errors.New(bundle.CodeDb)
	}

	return arr, err
}

// 取得子類別
func (d *Db) GetSubType(userId, main_id int) ([]bundle.Sub, error) {
	arr := make([]bundle.Sub, 0)
	s := `SELECT s.id, s.name, s.increase > 0 AS increase
			FROM sub_types AS s LEFT JOIN main_types AS m ON s.main_id=m.id 
			WHERE NOT s.deleted AND m.user_id=$1 AND s.main_id=$2 ORDER BY s.id`
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
	b, err := d.checkSub(userId, subId)
	if err != nil {
		return err
	}

	if !b {
		return errors.New(bundle.CodeHold)
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
	var id int
	s := `INSERT INTO main_types (user_id, name) 
			VALUES ($1, $2) RETURNING id`
	err := d.db.QueryRow(s, userId, name).Scan(&id)
	if err != nil {
		return 0, errors.New(bundle.CodeDb)
	}

	return id, nil
}

// 新增子類型
func (d *Db) InsertSubType(userId, mainId int, subName string, increase bool) (int, error) {
	i := -1
	if increase {
		i = 1
	}

	s := `SELECT COUNT(*) FROM main_types WHERE user_id=$1 AND id=$2 AND NOT deleted`
	count := 0
	err := d.db.QueryRow(s, userId, mainId).Scan(&count)
	if err != nil {
		return 0, errors.New(bundle.CodeHold)
	}

	var id int
	s = `INSERT INTO sub_types (name, main_id, increase) 
			VALUES ($1, $2, $3) RETURNING id`
	err = d.db.QueryRow(s, subName, mainId, i).Scan(&id)
	if err != nil {
		return 0, errors.New(bundle.CodeDb)
	}

	return id, nil
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
	b, err := d.checkSub(userId, subId)
	if err != nil {
		return err
	}

	if !b {
		return errors.New(bundle.CodeHold)
	}

	s := `UPDATE bills SET name=$1, sub_id=$2, price=$3, remark=$4, date=$7, updated=CURRENT_TIMESTAMP
			WHERE user_id=$5 AND id=$6`
	_, err = d.db.Exec(s, name, subId, price, remark, userId, itemId, date)
	if err != nil {
		return errors.New(bundle.CodeHold)
	}

	return nil
}

// 更新主類型
func (d *Db) UpdateMainType(userId, id int, name string) (err error) {
	s := `UPDATE main_types SET name=$1, updated=CURRENT_TIMESTAMP
			WHERE id=$2 AND user_id=$3`
	_, err = d.db.Exec(s, name, id, userId)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	return
}

// 更新子類型
func (d *Db) UpdateSubType(userId, id int, name string, increase bool) (err error) {
	i := -1
	if increase {
		i = 1
	}

	s := `UPDATE sub_types s SET name=$1, increase=$2, updated=CURRENT_TIMESTAMP
			FROM main_types m
			WHERE m.id=s.main_id AND s.id=$3 AND m.user_id=$4`
	_, err = d.db.Exec(s, name, i, id, userId)
	if err != nil {
		return errors.New(bundle.CodeDb)
	}

	return
}

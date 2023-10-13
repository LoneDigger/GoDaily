package db

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// docker run -it --rm --name test_postgres -p 5434:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres postgres:14.4
func newDb() *Db {
	return NewDb("127.0.0.1", "postgres", "postgres", "postgres")
}

func TestDeleteItem(t *testing.T) {
	d := newDb()
	err := d.DeleteItem(1, 3)
	if err != nil {
		log.Fatal(err)
	}
}

func TestDeleteMainType(t *testing.T) {
	d := newDb()
	err := d.DeleteMainType(1, 3)
	if err != nil {
		log.Fatal(err)
	}
}

func TestDeleteSubType(t *testing.T) {
	d := newDb()
	err := d.DeleteSubType(1, 1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllType(t *testing.T) {
	d := newDb()
	all, err := d.GetAllType(1)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, a := range all {
		for _, b := range a.Subs {
			fmt.Println(a.Id, a.Name, b.Id, b.Name)
		}
	}
}

func TestGetItem(t *testing.T) {
	d := newDb()
	item, err := d.GetItem(1, -10)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(item)
}

func TestGetPerviewItemsByDate(t *testing.T) {
	d := newDb()
	now := time.Now()
	end := now.Format("2006-01-02")
	start := now.AddDate(0, 0, -5).Format("2006-01-02")
	items, err := d.GetPerviewItemsByDate(-1, start, end)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(len(items))
}

func TestGetMainType(t *testing.T) {
	d := newDb()
	all, err := d.GetMainType(1)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, a := range all {
		fmt.Println(a.Id, a.Name)
	}
}

func TestGetSubType(t *testing.T) {
	d := newDb()
	all, err := d.GetSubType(1, -1)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, a := range all {
		fmt.Println(a.Id, a.Name)
	}
}

func TestGetSumByMainType(t *testing.T) {
	d := newDb()
	now := time.Now()
	end := now.Format("2006-01-02")
	start := now.AddDate(0, 0, -5).Format("2006-01-02")
	all, err := d.GetSumByMainType(1, start, end)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, a := range all {
		fmt.Println(a.Name, a.Sum)
	}
}

func TestGetSumByMonth(t *testing.T) {
	d := newDb()
	now := time.Now()
	end := now.Format("2006-01-02")
	start := now.AddDate(0, -3, 0).Format("2006-01-02")
	all, err := d.GetSumByMonth(1, start, end)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, a := range all {
		fmt.Println(a.Date, a.Sum)
	}
}

func TestInsertItem(t *testing.T) {
	d := newDb()
	err := d.InsertItem(1, "test", 6, 10, "", "2022-10-10")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func TestInsertMainType(t *testing.T) {
	d := newDb()
	i, err := d.InsertMainType(1, "test")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(i)
}

func TestInsertSubType(t *testing.T) {
	d := newDb()
	i, err := d.InsertSubType(1, 1, "test", true)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(i)
}

func TestUpdateItem(t *testing.T) {
	d := newDb()
	err := d.UpdateItem(1, -10, "test", 1, 100, "remark", "2020-01-011")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func TestUpdateMainType(t *testing.T) {
	d := newDb()
	err := d.UpdateMainType(1, -1, "test")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func TestUpdateSubType(t *testing.T) {
	d := newDb()
	err := d.UpdateSubType(1, -1, "Test", false)
	if err != nil {
		log.Fatal(err)
		return
	}
}

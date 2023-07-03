package main

import (
	"context"
	"fmt"
	"github.com/jinquanbao/lx-learning-go/basic/models"
	"reflect"
)

type BaseEntity = models.BaseEntity
type ExampleM = *ExampleEntity

type DBBaseStore[M BaseEntity] interface {
	Create(ctx context.Context, value M) error
	Save(ctx context.Context, value M, omit ...string) error
	UpdateById(ctx context.Context, value M, selected ...string) error
	FirstById(ctx context.Context, id int) (M, error)
	GetById(ctx context.Context, id int) (M, error)
	GetByIds(ctx context.Context, ids interface{}) ([]M, error)
	DeleteByIds(ctx context.Context, ids interface{}) error
	GetMapByIds(ctx context.Context, ids []int) (map[int]M, error)
}

var _ DBBaseStore[BaseEntity] = &dbBase[BaseEntity]{}

type dbBase[M BaseEntity] struct {
	baseEntity M
}

func (ds dbBase[M]) Create(ctx context.Context, value M) error {
	//TODO implement me
	panic("implement me")
}

func (ds dbBase[M]) Save(ctx context.Context, value M, omit ...string) error {
	//TODO implement me
	panic("implement me")
}

func (ds dbBase[M]) UpdateById(ctx context.Context, value M, selected ...string) error {
	//TODO implement me
	panic("implement me")
}

func (ds dbBase[M]) FirstById(ctx context.Context, id int) (M, error) {
	//TODO implement me
	panic("implement me")
}

func (ds dbBase[M]) GetById(ctx context.Context, id int) (M, error) {
	var m M

	elemType := reflect.TypeOf(m)
	fmt.Printf(" elemType=%v", elemType)
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
		fmt.Printf(" elemType=%v", elemType)
		val := reflect.New(elemType)
		fmt.Printf(" kind=%v", val.Kind())
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				val.Set(reflect.New(elemType))
			}
		}
		val.Elem().FieldByName("Id").SetInt(1)
		val.Elem().FieldByName("Name").SetString("exam")

		a := val.Interface()
		if v, ok := a.(M); ok {
			m = v
		}
	}

	return m, nil
}

func (ds dbBase[M]) GetByIds(ctx context.Context, ids interface{}) ([]M, error) {
	//TODO implement me
	panic("implement me")
}

func (ds dbBase[M]) DeleteByIds(ctx context.Context, ids interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (ds dbBase[M]) GetMapByIds(ctx context.Context, ids []int) (map[int]M, error) {
	//TODO implement me
	panic("implement me")
}

type ExampleStore interface {
	DBBaseStore[ExampleM]
}

var _ ExampleStore = example{}

type example struct {
	dbBase[ExampleM]
}

type ExampleEntity struct {
	Id   int
	Name string
}

func (e ExampleEntity) TableName() string {
	//TODO implement me
	panic("implement me")
}

func main() {
	var m ExampleEntity
	example := example{dbBase: dbBase[*ExampleEntity]{baseEntity: &m}}

	r, _ := example.GetById(context.Background(), 1)
	fmt.Println()
	fmt.Println(r)
}

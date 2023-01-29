package repository

import (
	"context"
	"encoding/json"
	"github.com/whj1990/go-core/handler"
	"github.com/whj1990/go-core/util"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type BaseCommonDBData struct {
	Id            int64 `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CreatedUserId int64
	UpdatedUserId int64
	Deleted       int
}

type BaseDBData struct {
	BaseCommonDBData
	OrganizationId int64
}

type basePrivateRepo interface {
	Db(ctx context.Context) *gorm.DB

	BaseGet(ctx context.Context, qb *QueryBuilder) (interface{}, error)
	BaseList(ctx context.Context, qb *QueryBuilder) (interface{}, int64, error)
	BaseCount(ctx context.Context, qb *QueryBuilder) (int64, error)

	BaseUpdate(ctx context.Context, qb *QueryBuilder, value interface{}) error
	BaseUpdateWithTX(ctx context.Context, tx *gorm.DB, qb *QueryBuilder, value interface{}) error
}

type BaseCommonRepo interface {
	AddWithTX(ctx context.Context, tx *gorm.DB, value interface{}) (int64, error)
	BatchAddWithTX(ctx context.Context, tx *gorm.DB, value interface{}) ([]int64, error)
	UpdateByIdWithTX(ctx context.Context, tx *gorm.DB, id int64, value interface{}) error
	DeleteByIdWithTX(ctx context.Context, tx *gorm.DB, id int64) error

	Add(ctx context.Context, value interface{}) (int64, error)
	BatchAdd(ctx context.Context, value interface{}) ([]int64, error)
	UpdateById(ctx context.Context, id int64, value interface{}) error
	DeleteById(ctx context.Context, id int64) error

	Transaction(f func(tx *gorm.DB) error) error
}

type BaseRepo struct {
	Db    *gorm.DB
	Model interface{}
}

func setDBDataFieldValue(f string, fieldValue reflect.Value, v interface{}) {
	if value, ok := v.([]uint8); ok {
		v = string(value)
	}
	switch fieldValue.FieldByName(util.CamelString(f)).Type().String() {
	case "datatypes.JSON":
		var jsonData datatypes.JSON
		json.Unmarshal([]byte(v.(string)), &jsonData)
		v = jsonData
	case "*time.Time":
		if value, ok := v.(time.Time); ok {
			v = &value
		}
	case "int":
		switch v.(type) {
		case int64:
			v = int(v.(int64))
		case string:
			v, _ = strconv.Atoi(v.(string))
		}
	case "int32":
		switch v.(type) {
		case int64:
			v = int32(v.(int64))
		case string:
			v2, _ := strconv.Atoi(v.(string))
			v = int32(v2)
		}
	case "int64":
		switch v.(type) {
		case string:
			v, _ = strconv.ParseInt(v.(string), 10, 64)
		}
	case "float64":
		switch v.(type) {
		case float64:
			v = v.(float64)
		case string:
			v, _ = strconv.ParseFloat(v.(string), 64)
		}
	}
	fieldValue.FieldByName(util.CamelString(f)).Set(reflect.ValueOf(v))
}

func (r *BaseRepo) convertDBData(dataMap map[string]interface{}) reflect.Value {
	dbData := reflect.New(reflect.TypeOf(r.Model)).Elem()
	for k, v := range dataMap {
		if v != nil {
			keyArr := strings.Split(k, "__")
			if len(keyArr) == 1 {
				setDBDataFieldValue(k, dbData, v)
			} else {
				joinTableAliasArr := strings.Split(keyArr[0], ".")
				joinFieldValue := dbData
				for _, f := range joinTableAliasArr {
					joinFieldValue = joinFieldValue.FieldByName(f)
				}
				setDBDataFieldValue(keyArr[1], joinFieldValue, v)
			}
		}
	}
	return dbData
}

func (r *BaseRepo) BaseGet(ctx context.Context, qb *QueryBuilder) (interface{}, error) {
	organizationId, err := util.GetMetaInfoCurrentOrganizationId(ctx)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	var dataMap map[string]interface{}
	err = qb.build(r.Db, r.Model, ctx, organizationId).Find(&dataMap).Error
	return r.convertDBData(dataMap).Addr().Interface(), handler.HandleError(err)
}

func (r *BaseRepo) BaseList(ctx context.Context, qb *QueryBuilder) (interface{}, int64, error) {
	organizationId, err := util.GetMetaInfoCurrentOrganizationId(ctx)
	if err != nil {
		return nil, 0, handler.HandleError(err)
	}
	count, err := r.BaseCount(ctx, qb)
	if err != nil || count == 0 {
		return nil, count, handler.HandleError(err)
	}
	var dataMapList []map[string]interface{}
	err = qb.build(r.Db, r.Model, ctx, organizationId).Find(&dataMapList).Error
	if err != nil {
		return nil, count, handler.HandleError(err)
	}
	result := reflect.New(reflect.SliceOf(reflect.TypeOf(r.Model)))
	resultValue := make([]reflect.Value, len(dataMapList))
	for i := 0; i < len(dataMapList); i++ {
		resultValue[i] = r.convertDBData(dataMapList[i])
	}
	resultElem := result.Elem()
	resultElem.Set(reflect.Append(resultElem, resultValue...))
	return result.Interface(), count, nil
}

func (r *BaseRepo) BaseCount(ctx context.Context, qb *QueryBuilder) (int64, error) {
	organizationId, err := util.GetMetaInfoCurrentOrganizationId(ctx)
	if err != nil {
		return 0, handler.HandleError(err)
	}
	var count int64
	err = qb.cleanListParam().build(r.Db, r.Model, ctx, organizationId).Count(&count).Error
	return count, handler.HandleError(err)
}

func (r *BaseRepo) AddWithTX(ctx context.Context, tx *gorm.DB, value interface{}) (int64, error) {
	currentUserId, err := util.GetMetaInfoCurrentUserId(ctx)
	if err != nil {
		return 0, handler.HandleError(err)
	}
	currentOrganizationId, err := util.GetMetaInfoCurrentOrganizationId(ctx)
	if err != nil {
		return 0, handler.HandleError(err)
	}
	data := reflect.ValueOf(value).Elem()
	data.FieldByName("CreatedUserId").Set(reflect.ValueOf(currentUserId))
	data.FieldByName("UpdatedUserId").Set(reflect.ValueOf(currentUserId))
	_, ok := data.Type().FieldByName("OrganizationId")
	if ok && currentOrganizationId != 0 {
		data.FieldByName("OrganizationId").Set(reflect.ValueOf(currentOrganizationId))
	}
	err = tx.WithContext(ctx).Model(r.Model).Create(value).Error
	if err != nil {
		return 0, handler.HandleError(err)
	}
	return data.FieldByName("Id").Int(), nil
}

func (r *BaseRepo) BatchAddWithTX(ctx context.Context, tx *gorm.DB, value interface{}) ([]int64, error) {
	currentUserId, err := util.GetMetaInfoCurrentUserId(ctx)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	currentOrganizationId, err := util.GetMetaInfoCurrentOrganizationId(ctx)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	reflectValue := reflect.Indirect(reflect.ValueOf(value))
	for i := 0; i < reflectValue.Len(); i++ {
		data := reflectValue.Index(i)
		data.FieldByName("CreatedUserId").Set(reflect.ValueOf(currentUserId))
		data.FieldByName("UpdatedUserId").Set(reflect.ValueOf(currentUserId))
		_, ok := data.Type().FieldByName("OrganizationId")
		if ok && currentOrganizationId != 0 {
			data.FieldByName("OrganizationId").Set(reflect.ValueOf(currentOrganizationId))
		}
	}
	err = tx.WithContext(ctx).Model(r.Model).Create(value).Error
	if err != nil {
		return nil, handler.HandleError(err)
	}
	ids := make([]int64, reflectValue.Len())
	for i := 0; i < reflectValue.Len(); i++ {
		data := reflectValue.Index(i)
		ids[i] = data.FieldByName("Id").Int()
	}
	return ids, nil
}

func (r *BaseRepo) BaseUpdateWithTX(ctx context.Context, tx *gorm.DB, qb *QueryBuilder, value interface{}) error {
	currentUserId, err := util.GetMetaInfoCurrentUserId(ctx)
	if err != nil {
		return handler.HandleError(err)
	}
	data := reflect.ValueOf(value).Elem()
	data.FieldByName("UpdatedUserId").Set(reflect.ValueOf(currentUserId))
	data.FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now()))
	return handler.HandleError(qb.cleanListParam().buildUpdate(tx, r.Model, ctx).Updates(value).Error)
}

func (r *BaseRepo) UpdateByIdWithTX(ctx context.Context, tx *gorm.DB, id int64, value interface{}) error {
	return r.BaseUpdateWithTX(ctx, tx, NewQueryBuilder().Where([]*Condition{
		NewAndCondition("id = ?", id),
	}), value)
}

func (r *BaseRepo) DeleteByIdWithTX(ctx context.Context, tx *gorm.DB, id int64) error {
	value := reflect.New(reflect.TypeOf(r.Model)).Interface()
	data := reflect.ValueOf(value).Elem()
	data.FieldByName("Deleted").Set(reflect.ValueOf(1))
	return r.UpdateByIdWithTX(ctx, tx, id, value)
}

func (r *BaseRepo) Add(ctx context.Context, value interface{}) (int64, error) {
	return r.AddWithTX(ctx, r.Db, value)
}

func (r *BaseRepo) BatchAdd(ctx context.Context, value interface{}) ([]int64, error) {
	return r.BatchAddWithTX(ctx, r.Db, value)
}

func (r *BaseRepo) UpdateById(ctx context.Context, id int64, value interface{}) error {
	return r.UpdateByIdWithTX(ctx, r.Db, id, value)
}

func (r *BaseRepo) BaseUpdate(ctx context.Context, qb *QueryBuilder, value interface{}) error {
	return r.BaseUpdateWithTX(ctx, r.Db, qb, value)
}

func (r *BaseRepo) DeleteById(ctx context.Context, id int64) error {
	return r.DeleteByIdWithTX(ctx, r.Db, id)
}

func (r *BaseRepo) Transaction(f func(tx *gorm.DB) error) error {
	return r.Db.Transaction(f)
}

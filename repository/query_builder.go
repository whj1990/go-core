package repository

import (
	"context"
	"fmt"
	"github.com/whj1990/go-core/util"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
)

const (
	AND = 0
	OR  = 1
)

type QueryBuilder struct {
	selectFields []string
	where        []*Condition
	order        string
	pageNum      int32
	pageSize     int32
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func NewListQueryBuilder(order string, pageNum, pageSize int32) *QueryBuilder {
	return &QueryBuilder{order: order, pageNum: pageNum, pageSize: pageSize}
}

func (qb *QueryBuilder) Select(fields []string) *QueryBuilder {
	qb.selectFields = fields
	return qb
}

func (qb *QueryBuilder) Where(condition []*Condition) *QueryBuilder {
	qb.where = condition
	return qb
}

func (qb *QueryBuilder) cleanListParam() *QueryBuilder {
	newQb := *qb
	newQb.order = ""
	newQb.pageNum = 0
	newQb.pageSize = 0
	return &newQb
}

func (qb *QueryBuilder) buildUpdate(db *gorm.DB, model interface{}, ctx context.Context) *gorm.DB {
	tx := (*db).WithContext(ctx).Model(model)
	if len(qb.selectFields) > 0 {
		tx = tx.Select(qb.selectFields)
	}
	if len(qb.where) > 0 {
		tx = tx.Where(buildCondition(qb.where, db))
	}
	return tx
}

func (qb *QueryBuilder) build(db *gorm.DB, model interface{}, ctx context.Context, organizationId int64) *gorm.DB {
	tx := (*db).WithContext(ctx).Model(model)
	t := reflect.TypeOf(model)
	mainTableName := getTableName(t)
	var totalSelectFields []string
	addMainSelectFields(t, mainTableName, &totalSelectFields)
	buildJoinPreloading(tx, db, t, mainTableName, &totalSelectFields, []string{}, organizationId)
	if len(qb.selectFields) > 0 {
		tx = tx.Select(qb.selectFields)
	} else {
		tx = tx.Select(totalSelectFields)
	}
	value := reflect.New(reflect.TypeOf(model)).Interface()
	tx = tx.Where(buildDefaultQuery(value.(schema.Tabler).TableName(), organizationId, t))

	if len(qb.where) > 0 {
		tx = tx.Where(buildCondition(qb.where, db))
	}
	if qb.order != "" {
		tx = tx.Order(qb.order)
	}
	if qb.pageNum > 0 && qb.pageSize > 0 {
		tx = tx.Offset(int((qb.pageNum - 1) * qb.pageSize))
		tx = tx.Limit(int(qb.pageSize))
	}
	return tx
}

func getTagForeignKeyReferences(field reflect.StructField) (string, string) {
	foreignKey := ""
	references := "Id"
	tag := field.Tag.Get("gorm")
	if tag != "" {
		for _, v := range strings.Split(strings.ReplaceAll(tag, " ", ""), ";") {
			if strings.HasPrefix(v, "foreignKey:") {
				foreignKey = v[11:]
				continue
			}
			if strings.HasPrefix(v, "references:") {
				references = v[11:]
				continue
			}
		}
	}
	return foreignKey, references
}

func buildJoinPreloading(tx *gorm.DB, db *gorm.DB, t reflect.Type, mainTableName string, totalSelectFields *[]string, parent []string, organizationId int64) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		foreignKey, references := getTagForeignKeyReferences(field)
		if foreignKey != "" {
			current := append(parent, field.Name)
			joinTableAlias := getTableAlias(current)
			addJoinSelectFields(field.Type, joinTableAlias, totalSelectFields)
			joinTableName := getTableName(field.Type)
			mainTableAlias := mainTableName
			if len(parent) > 0 {
				mainTableAlias = getTableAlias(parent)
			}
			tx = tx.Joins(fmt.Sprintf("LEFT JOIN `%s` `%s` ON `%s`.%s = `%s`.`%s` AND %s",
				joinTableName, joinTableAlias, mainTableAlias, util.SnakeString(foreignKey), joinTableAlias,
				util.SnakeString(references), buildDefaultQuery(joinTableAlias, organizationId, field.Type)))
			buildJoinPreloading(tx, db, field.Type, mainTableName, totalSelectFields, current, organizationId)
		} else if field.Anonymous {
			buildJoinPreloading(tx, db, field.Type, mainTableName, totalSelectFields, parent, organizationId)
		}
	}
}

func addMainSelectFields(t reflect.Type, tableName string, totalSelectFields *[]string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if foreignKey, _ := getTagForeignKeyReferences(field); foreignKey == "" {
			if !field.Anonymous {
				*totalSelectFields = append(*totalSelectFields, fmt.Sprintf("`%s`.`%s`", tableName, util.SnakeString(field.Name)))
			} else {
				subField, _ := t.FieldByName(field.Name)
				addMainSelectFields(subField.Type, tableName, totalSelectFields)
			}
		}
	}
}

func addJoinSelectFields(t reflect.Type, tableAlias string, totalSelectFields *[]string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if foreignKey, _ := getTagForeignKeyReferences(field); foreignKey == "" {
			if !field.Anonymous {
				fieldSnakeStr := util.SnakeString(field.Name)
				*totalSelectFields = append(*totalSelectFields, fmt.Sprintf("`%s`.`%s` AS `%s__%s`", tableAlias, fieldSnakeStr, tableAlias, fieldSnakeStr))
			} else {
				subField, _ := t.FieldByName(field.Name)
				addJoinSelectFields(subField.Type, tableAlias, totalSelectFields)
			}
		}
	}
}

func getTableAlias(arr []string) string {
	return strings.Join(arr, ".")
}

func getTableName(t reflect.Type) string {
	return reflect.New(t).MethodByName("TableName").Call([]reflect.Value{})[0].Interface().(string)
}

func buildDefaultQuery(alias string, organizationId int64, t reflect.Type) string {
	_, ok := t.FieldByName("OrganizationId")
	if organizationId != 0 && ok {
		return fmt.Sprintf("`%s`.`deleted` = 0 and `%s`.`organization_id` = %d", alias, alias, organizationId)
	} else {
		return fmt.Sprintf("`%s`.`deleted` = 0", alias)
	}
}

func buildCondition(conditions []*Condition, db *gorm.DB) *gorm.DB {
	var result = db
	for _, condition := range conditions {
		subConditions, ok := condition.query.([]*Condition)
		if ok {
			switch condition.andOr {
			case AND:
				result = result.Where(buildCondition(subConditions, db))
			case OR:
				result = result.Or(buildCondition(subConditions, db))
			}
		} else {
			switch condition.andOr {
			case AND:
				result = result.Where(condition.query, condition.args...)
			case OR:
				result = result.Or(condition.query, condition.args...)
			}
		}
	}
	return result
}

type Condition struct {
	andOr int
	query interface{}
	args  []interface{}
}

func NewAndCondition(query interface{}, args ...interface{}) *Condition {
	return &Condition{AND, query, args}
}

func NewOrCondition(query interface{}, args ...interface{}) *Condition {
	return &Condition{OR, query, args}
}

func NewAndSubCondition(query []*Condition) *Condition {
	return &Condition{andOr: AND, query: query}
}

func NewOrSubCondition(query []*Condition) *Condition {
	return &Condition{andOr: OR, query: query}
}

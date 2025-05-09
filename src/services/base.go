package services

import (
	"context"
	"fmt"
	"github.com/wahid-wex/e-commerce-backend/api/dto"
	"github.com/wahid-wex/e-commerce-backend/api/error_handler"
	"github.com/wahid-wex/e-commerce-backend/common"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/data/db"
	logging "github.com/wahid-wex/e-commerce-backend/logs"
	"gorm.io/gorm"
	"math"
	"reflect"
	"strings"
)

type preload struct {
	string
}

type BaseService[T any, Tc any, Tu any, Tr any] struct {
	Database *gorm.DB
	Logger   logging.Logger
	Preloads []preload
}

func NewBaseService[T any, Tc any, Tu any, Tr any](cfg *config.Config) *BaseService[T, Tc, Tu, Tr] {
	return &BaseService[T, Tc, Tu, Tr]{
		Database: db.GetDb(),
		Logger:   logging.NewLogger(cfg),
	}
}

func (s *BaseService[T, Tc, Tu, Tr]) Create(ctx context.Context, req *Tc) (*Tr, error) {
	model, _ := common.TypeConverter[T](req)
	tx := s.Database.WithContext(ctx).Begin()
	err := tx.
		Create(model).
		Error
	if err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	bm, _ := common.TypeConverter[gorm.Model](model)
	id := int(bm.ID)
	return s.GetById(ctx, id)
}
func (s *BaseService[T, Tc, Tu, Tr]) Update(ctx context.Context, id int, req *Tu) (*Tr, error) {
	updateMap, _ := common.TypeConverter[map[string]interface{}](req)
	snakeMap := map[string]interface{}{}
	for k, v := range *updateMap {
		snakeMap[common.ToSnakeCase(k)] = v
	}
	model := new(T)
	tx := s.Database.WithContext(ctx).Begin()
	if err := tx.Model(model).
		Where("id = ?", id).
		Updates(snakeMap).
		Error; err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	return s.GetById(ctx, id)
}

func (s *BaseService[T, Tc, Tu, Tr]) Delete(ctx context.Context, id int) error {
	err := s.Database.Delete("id = ?", id).Error

	if err != nil {
		s.Logger.Error(logging.Postgres, logging.Delete, error_handler.RecordNotFound, nil)
		return &error_handler.ServiceError{EndUserMessage: error_handler.RecordNotFound}
	}
	return nil
}

func (s *BaseService[T, Tc, Tu, Tr]) GetById(ctx context.Context, id int) (*Tr, error) {
	ID := uint(id)
	model := new(T)
	db := Preload(s.Database, s.Preloads)
	err := db.
		Where("id = ?", ID).
		First(model).
		Error
	if err != nil {
		return nil, err
	}
	return common.TypeConverter[Tr](model)
}

func (s *BaseService[T, Tc, Tu, Tr]) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter) (*dto.PagedList[Tr], error) {
	res, err := Paginate[T, Tr](req, s.Preloads, s.Database)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewPagedList[T any](items *[]T, count int64, pageNumber int, pageSize int64) *dto.PagedList[T] {
	pl := &dto.PagedList[T]{
		PageNumber: pageNumber,
		TotalRows:  count,
		Items:      items,
	}
	pl.TotalPages = int(math.Ceil(float64(count) / float64(pageSize)))
	pl.HasNextPage = pl.PageNumber < pl.TotalPages
	pl.HasPreviousPage = pl.PageNumber > 1

	return pl
}

func Paginate[T any, Tr any](pagination *dto.PaginationInputWithFilter, preloads []preload, db *gorm.DB) (*dto.PagedList[Tr], error) {
	model := new(T)
	var items *[]T
	var rItems *[]Tr
	db = Preload(db, preloads)
	query := getQuery[T](&pagination.DynamicFilter)
	sort := getSort[T](&pagination.DynamicFilter)

	var totalRows int64 = 0

	db.
		Model(model).
		Where(query).
		Count(&totalRows)

	err := db.
		Where(query).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetPageSize()).
		Order(sort).
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}
	rItems, err = common.TypeConverter[[]Tr](items)
	if err != nil {
		return nil, err
	}
	return NewPagedList(rItems, totalRows, pagination.PageNumber, int64(pagination.PageSize)), err

}

// a short sample query
// {
//
//	"filter": {
//		"Name": { // Capital first letter
//			"filterType": "text",
//			"from": "x",
//			"to": "",
//			"type": "equals"
//		}
//	},
//
//		"pageNumber":1,
//		"pageSize": 10,
//		"sort": []
//	}
//
// getQuery
func getQuery[T any](filter *dto.DynamicFilter) string {
	t := new(T)
	typeT := reflect.TypeOf(*t)
	query := make([]string, 0)
	if filter.Filter != nil {
		for name, filter := range filter.Filter {
			if fld, ok := typeT.FieldByName(name); ok {
				query = append(query, generateDynamicFilter(fld, filter))
			}
		}
	}
	return strings.Join(query, " AND ")
}

func generateDynamicFilter(fld reflect.StructField, filter dto.Filter) string {
	conditionQuery := ""
	fld.Name = common.ToSnakeCase(fld.Name)
	switch filter.Type {
	case "contains":
		conditionQuery = fmt.Sprintf("%s ILike '%%%s%%'", fld.Name, filter.From)
	case "notContains":
		conditionQuery = fmt.Sprintf("%s not ILike '%%%s%%'", fld.Name, filter.From)
	case "startsWith":
		conditionQuery = fmt.Sprintf("%s ILike '%s%%'", fld.Name, filter.From)
	case "endsWith":
		conditionQuery = fmt.Sprintf("%s ILike '%%%s'", fld.Name, filter.From)
	case "equals":
		conditionQuery = fmt.Sprintf("%s = '%s'", fld.Name, filter.From)
	case "notEqual":
		conditionQuery = fmt.Sprintf("%s != '%s'", fld.Name, filter.From)
	case "lessThan":
		conditionQuery = fmt.Sprintf("%s < %s", fld.Name, filter.From)
	case "lessThanOrEqual":
		conditionQuery = fmt.Sprintf("%s <= %s", fld.Name, filter.From)
	case "greaterThan":
		conditionQuery = fmt.Sprintf("%s > %s", fld.Name, filter.From)
	case "greaterThanOrEqual":
		conditionQuery = fmt.Sprintf("%s >= %s", fld.Name, filter.From)
	case "inRange":
		if fld.Type.Kind() == reflect.String {
			conditionQuery = fmt.Sprintf("%s >= '%s%%' AND ", fld.Name, filter.From)
			conditionQuery += fmt.Sprintf("%s <= '%%%s'", fld.Name, filter.To)
		} else {
			conditionQuery = fmt.Sprintf("%s >= %s AND ", fld.Name, filter.From)
			conditionQuery += fmt.Sprintf("%s <= %s", fld.Name, filter.To)
		}
	}
	return conditionQuery
}

// getSort
func getSort[T any](filter *dto.DynamicFilter) string {
	t := new(T)
	typeT := reflect.TypeOf(*t)
	sort := make([]string, 0)
	if filter.Sort != nil {
		for _, tp := range *filter.Sort {
			fld, ok := typeT.FieldByName(tp.ColId)
			if ok && (tp.Sort == "asc" || tp.Sort == "desc") {
				fld.Name = common.ToSnakeCase(fld.Name)
				sort = append(sort, fmt.Sprintf("%s %s", fld.Name, tp.Sort))
			}
		}
	}
	return strings.Join(sort, ", ")
}

// Preload
func Preload(db *gorm.DB, preloads []preload) *gorm.DB {
	for _, item := range preloads {
		db = db.Preload(item.string)
	}
	return db
}

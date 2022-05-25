package repositories

import (
	"context"
	"fmt"
	"errors"
	"sync"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"github.com/duolacloud/crud-core/types"
	"github.com/duolacloud/crud-core-gorm/query"
	"github.com/mitchellh/mapstructure"
)

type GormCrudRepositoryOptions struct {
}

type GormCrudRepositoryOption func (*GormCrudRepositoryOptions)

type GormCrudRepository[DTO any, CreateDTO any, UpdateDTO any] struct {
	DB *gorm.DB
	Schema *schema.Schema
	Options *GormCrudRepositoryOptions
}

func NewGormCrudRepository[DTO any, CreateDTO any, UpdateDTO any](
	DB *gorm.DB,
	opts ...GormCrudRepositoryOption,
) *GormCrudRepository[DTO, CreateDTO, UpdateDTO] {
	r := &GormCrudRepository[DTO, CreateDTO, UpdateDTO]{
		DB: DB,
	}

	var dto DTO
	r.Schema, _ = schema.Parse(&dto, &sync.Map{}, schema.NamingStrategy{})

	r.Options = &GormCrudRepositoryOptions{}
	for _, o := range opts {
		o(r.Options)
	}

	return r
}

func (r *GormCrudRepository[DTO, CreateDTO, UpdateDTO]) Create(c context.Context, createDTO *CreateDTO, opts ...types.CreateOption) (*DTO, error) {
	res := r.DB.Create(createDTO)
	if res.Error != nil {
		return nil, res.Error
	}

	var dto DTO
	_ = mapstructure.Decode(createDTO, &dto)

	return &dto, nil
}

func (r *GormCrudRepository[DTO, CreateDTO, UpdateDTO]) Delete(c context.Context, id types.ID) error {
	/*
	model, err := r.Get(c, id)
	if err != nil {
		return err
	}*/

	filter := make(map[string]interface{})

	if len(r.Schema.PrimaryFields) == 1 {
		fName := r.Schema.PrimaryFields[0].DBName
		filter[fName] = id	
	} else if len(r.Schema.PrimaryFields) > 1 {
		ids, ok := id.(map[string]interface{})
		if !ok {
			return errors.New("invalid id, not match")
		}

		if len(ids) != len(r.Schema.PrimaryFields) {
			return errors.New("invalid id, size not match")	
		}

		for _, primaryField := range r.Schema.PrimaryFields {
			filter[primaryField.DBName] = ids[primaryField.Name]
		}
	}

	fmt.Printf("PrimaryFields: table: %s, %v\n", r.Schema.Table, filter)

	var dto DTO
	res := r.DB.Delete(&dto, filter)
	return res.Error
}

func (r *GormCrudRepository[DTO, CreateDTO, UpdateDTO]) Update(c context.Context, id types.ID, updateDTO *UpdateDTO, opts ...types.UpdateOption) (*DTO, error) {
	res := r.DB.Save(updateDTO)
	if res.Error != nil {
		return nil, res.Error
	}

	// TODO
	return nil, nil
}

func (r *GormCrudRepository[DTO, CreateDTO, UpdateDTO]) Get(c context.Context, id types.ID) (*DTO, error) {
	var dto DTO
	err := r.DB.First(&dto, id).Error
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func (r *GormCrudRepository[DTO, CreateDTO, UpdateDTO]) Query(c context.Context, q *types.PageQuery) ([]*DTO, error) {
	filterQueryBuilder := query.NewFilterQueryBuilder(r.Schema)

	db, err := filterQueryBuilder.BuildQuery(q, r.DB)
	if err != nil {
		return nil, err
	}

	var dtos []*DTO
	res := db.Find(&dtos)
	if res.Error != nil {
		return nil, err
	}

	return dtos, nil
}

func (r *GormCrudRepository[DTO, CreateDTO, UpdateDTO]) Count(c context.Context, q *types.PageQuery) (int64, error) {
	/*
	filterQueryBuilder := query.NewFilterQueryBuilder[DTO](r.Schema, r.Options.StrictValidation)

	mq, err := filterQueryBuilder.BuildQuery(q);
	if err != nil {
		return 0, err
	}
	*/

	// TODO
	count := int64(0)
	return count, nil
}

func (r *GormCrudRepository[DTO, CreateDTO, UpdateDTO]) QueryOne(c context.Context, filter map[string]interface{}) (*DTO, error) {
	return nil, nil
}

func (r *GormCrudRepository[DTO, CreateDTO, UpdateDTO]) Aggregate(
	c context.Context,
	filter map[string]interface{},
	aggregateQuery *types.AggregateQuery,
) ([]*types.AggregateResponse, error) {
	return nil, nil
}

func (r *GormCrudRepository[DTO, CreateDTO, UpdateDTO]) CursorQuery(c context.Context, query *types.CursorQuery) ([]*DTO, *types.CursorExtra, error) {
	return nil, nil, nil
}
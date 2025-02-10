package infrastructure

import (
	"context"
	"errors"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"gorm.io/gorm"
)

type DBDriver interface {
	Begin() Transaction

	First(ctx context.Context, model any, filter Filter) error
	Find(ctx context.Context, model any, filter Filter) error
	Create(ctx context.Context, model any) error
	Update(ctx context.Context, model any, filter Filter) error
	Delete(ctx context.Context, model any, filter Filter) error
	Save(ctx context.Context, model any) error
	Count(ctx context.Context, model any, filter Filter) (int64, error)
	Raw(ctx context.Context, model any, query string, args []string) error

	FirstTx(ctx context.Context, model any, filter Filter, tx Transaction) error
	FindTx(ctx context.Context, model any, filter Filter, tx Transaction) error
	CreateTx(ctx context.Context, model any, tx Transaction) error
	UpdateTx(ctx context.Context, model any, filter Filter, tx Transaction) error
	DeleteTx(ctx context.Context, model any, filter Filter, tx Transaction) error
	SaveTx(ctx context.Context, model any, tx Transaction) error
	CountTx(ctx context.Context, model any, filter Filter, tx Transaction) (int64, error)
	RawTx(ctx context.Context, model any, query string, args []string, tx Transaction) error
}

type dbDriver struct {
	db *gorm.DB
}

func NewDBDriver(db *gorm.DB) DBDriver {
	return &dbDriver{db: db}
}

func (d *dbDriver) Begin() Transaction {
	tx := d.db.Begin()
	return &Tx{tx: tx}
}

func (d *dbDriver) First(ctx context.Context, model any, filter Filter) error {
	var result *gorm.DB
	if filter == nil {
		result = d.db.WithContext(ctx).First(model)
	} else {
		result = filter.GenerateMods(d.db.WithContext(ctx)).First(model)
	}
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return objects.ErrNotFound
		}
		return result.Error
	}
	return nil
}

func (d *dbDriver) Create(ctx context.Context, model any) error {
	return d.db.WithContext(ctx).Create(model).Error
}

func (d *dbDriver) Find(ctx context.Context, model any, filter Filter) error {
	if filter == nil {
		return d.db.WithContext(ctx).Find(model).Error
	}
	return filter.GenerateMods(d.db.WithContext(ctx)).Find(model, filter).Error
}

func (d *dbDriver) Update(ctx context.Context, model any, filter Filter) error {
	if filter == nil {
		return d.db.WithContext(ctx).Updates(model).Error
	}
	table := d.db.WithContext(ctx).Model(model).Select("*").Omit("created_at")
	return filter.GenerateMods(table).Updates(model).Error
}

func (d *dbDriver) Delete(ctx context.Context, model any, filter Filter) error {
	if filter == nil {
		return d.db.WithContext(ctx).Delete(model).Error
	}
	return filter.GenerateMods(d.db.WithContext(ctx)).Delete(model).Error
}

func (d *dbDriver) Save(ctx context.Context, model any) error {
	return d.db.WithContext(ctx).Save(model).Error
}

func (d *dbDriver) Count(ctx context.Context, model any, filter Filter) (int64, error) {
	var count int64
	var result *gorm.DB
	if filter == nil {
		result = d.db.WithContext(ctx).Model(model).Count(&count)
	} else {
		result = filter.GenerateMods(d.db.WithContext(ctx)).Model(model).Count(&count)
	}
	return count, result.Error
}

func (d *dbDriver) Raw(ctx context.Context, model any, query string, args []string) error {
	return d.db.WithContext(ctx).Raw(query, args).Scan(model).Error
}

func (d *dbDriver) FirstTx(ctx context.Context, model any, filter Filter, tx Transaction) error {
	var result *gorm.DB
	if filter == nil {
		result = tx.GetTx().WithContext(ctx).First(model)
	} else {
		result = filter.GenerateMods(tx.GetTx().WithContext(ctx)).First(model)
	}
	return result.Error
}

func (d *dbDriver) CreateTx(ctx context.Context, model any, tx Transaction) error {
	return tx.GetTx().WithContext(ctx).Create(model).Error
}

func (d *dbDriver) FindTx(ctx context.Context, model any, filter Filter, tx Transaction) error {
	if filter == nil {
		return tx.GetTx().WithContext(ctx).Find(model).Error
	}
	return filter.GenerateMods(tx.GetTx().WithContext(ctx)).Find(model).Error
}

func (d *dbDriver) UpdateTx(ctx context.Context, model any, filter Filter, tx Transaction) error {
	if filter == nil {
		return tx.GetTx().WithContext(ctx).Model(model).Select("*").Omit("created_at").Updates(model).Error
	}
	table := tx.GetTx().WithContext(ctx).Model(model).Select("*").Omit("created_at")
	return filter.GenerateMods(table).Updates(model).Error
}

func (d *dbDriver) DeleteTx(ctx context.Context, model any, filter Filter, tx Transaction) error {
	if filter == nil {
		return tx.GetTx().WithContext(ctx).Delete(model).Error
	}
	return filter.GenerateMods(tx.GetTx().WithContext(ctx)).Delete(model).Error
}

func (d *dbDriver) SaveTx(ctx context.Context, model any, tx Transaction) error {
	return tx.GetTx().WithContext(ctx).Save(model).Error
}

func (d *dbDriver) CountTx(ctx context.Context, model any, filter Filter, tx Transaction) (int64, error) {
	var count int64
	var result *gorm.DB
	if filter == nil {
		result = tx.GetTx().WithContext(ctx).Model(model).Count(&count)
	} else {
		result = filter.GenerateMods(tx.GetTx().WithContext(ctx)).Model(model).Count(&count)
	}
	return count, result.Error
}

func (d *dbDriver) RawTx(ctx context.Context, model any, query string, args []string, tx Transaction) error {
	return tx.GetTx().WithContext(ctx).Raw(query, args).Scan(model).Error
}

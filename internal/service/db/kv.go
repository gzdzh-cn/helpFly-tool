package db

import (
	"context"
	"fmt"

	"helpfly/internal/dao"
	"helpfly/internal/model/do"
)

func GetKV(ctx context.Context, key string) ([]byte, error) {
	if err := ready(); err != nil {
		return nil, err
	}
	var row struct {
		Value []byte `orm:"value"`
	}
	if err := dao.KvStore.Ctx(normalizeContext(ctx)).Where(dao.KvStore.Columns().ItemKey, key).Scan(&row); err != nil {
		if isNoRows(err) {
			return nil, notFoundError("KV 数据不存在")
		}
		return nil, fmt.Errorf("查询 KV 数据失败: %w", err)
	}
	return row.Value, nil
}

func SetKV(ctx context.Context, key string, value []byte) error {
	if err := ready(); err != nil {
		return err
	}
	_, err := dao.KvStore.Ctx(normalizeContext(ctx)).Data(do.KvStore{
		ItemKey: key,
		Value:   value,
	}).OnConflict(dao.KvStore.Columns().ItemKey).Save()
	if err != nil {
		return fmt.Errorf("保存 KV 数据失败: %w", err)
	}
	return nil
}

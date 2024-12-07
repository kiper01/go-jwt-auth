package migration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(pool *pgxpool.Pool, migrationsPath string) error {

	ctx := context.TODO()

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return err
	}

	// Проверяем наличие каталога с миграциями
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("migrations directory '%s' does not exist", absPath)
	}

	// Получаем список SQL файлов в каталоге миграций
	files, err := os.ReadDir(absPath)
	if err != nil {
		return err
	}

	// Открываем соединение с базой данных
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// Запускаем транзакцию
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Применяем каждый SQL файл миграции
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlBytes, err := os.ReadFile(filepath.Join(absPath, file.Name()))
			if err != nil {
				return err
			}

			// Выполняем SQL запрос
			if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
				return fmt.Errorf("failed to execute migration '%s': %v", file.Name(), err)
			}
		}
	}

	// Фиксируем транзакцию
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	fmt.Println("Migration successful!")
	return nil
}

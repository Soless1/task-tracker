package task

import "context"

func Delete(ctx context.Context, repo TaskRepository, ID int64) error {
	err := repo.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}

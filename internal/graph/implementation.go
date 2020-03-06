package graph

import (
	"context"
	"fmt"

	dbmodel "github.com/moutend/gqlgen-todoapp/internal/db/model"
	"github.com/moutend/gqlgen-todoapp/internal/graph/model"
	"github.com/moutend/gqlgen-todoapp/internal/jwt"
	"github.com/moutend/gqlgen-todoapp/internal/middleware/auth"
	"github.com/moutend/gqlgen-todoapp/internal/middleware/db"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func (r *mutationResolver) createTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	user := auth.ForContext(ctx)

	if user == nil {
		return nil, fmt.Errorf("access denied")
	}

	tx := db.ForContext(ctx)

	if tx == nil {
		return nil, fmt.Errorf("internal error")
	}

	task := &dbmodel.Task{
		Title:   null.StringFrom(input.Title),
		Content: null.StringFrom(input.Content),
		UserID:  user.ID,
	}

	if err := task.Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Error(err)
		return nil, err
	}

	result := &model.Task{
		ID:      fmt.Sprint(task.ID),
		Title:   task.Title.String,
		Content: task.Content.String,
		User: &model.User{
			ID:   fmt.Sprint(user.ID),
			Name: user.Name,
		},
	}

	return result, nil
}

func (r *mutationResolver) createUser(ctx context.Context, input model.NewUser) (string, error) {
	tx := db.ForContext(ctx)

	if tx == nil {
		return "", fmt.Errorf("internal error")
	}

	hash, err := auth.HashPassword(input.Password)

	if err != nil {
		return "", fmt.Errorf("failed to create user")
	}

	user := &dbmodel.User{
		Name:         input.Name,
		PasswordHash: hash,
	}

	if err := user.Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Error(err)
		return "", fmt.Errorf("failed to create user")
	}

	token, err := jwt.GenerateToken(user.Name)

	if err != nil {
		tx.Error(err)
		return "", fmt.Errorf("failed to create user")
	}

	return token, nil
}

func (r *mutationResolver) login(ctx context.Context, input model.Login) (string, error) {
	tx := db.ForContext(ctx)

	if tx == nil {
		return "", fmt.Errorf("internal error")
	}

	user, err := dbmodel.Users(dbmodel.UserWhere.Name.EQ(input.Name)).One(ctx, tx)

	if err != nil {
		tx.Error(err)
		return "", fmt.Errorf("failed to login")
	}
	if !auth.IsValidPassword(input.Password, user.PasswordHash) {
		return "", fmt.Errorf("failed to login")
	}

	token, err := jwt.GenerateToken(user.Name)

	if err != nil {
		return "", fmt.Errorf("failed to login")
	}

	return token, nil
}

func (r *mutationResolver) refreshToken(ctx context.Context, input model.RefreshTokenInput) (token string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("access denied")
		}
	}()

	var name string

	name, err = jwt.ParseToken(input.Token)

	if err != nil {
		return "", fmt.Errorf("access denied")
	}

	token, err = jwt.GenerateToken(name)

	if err != nil {
		return "", fmt.Errorf("access denied")
	}

	return token, nil
}

func (r *queryResolver) tasks(ctx context.Context) ([]*model.Task, error) {
	user := auth.ForContext(ctx)

	if user == nil {
		return nil, fmt.Errorf("access denied")
	}

	tx := db.ForContext(ctx)

	if tx == nil {
		return nil, fmt.Errorf("internal error")
	}

	tasks, err := dbmodel.Tasks(dbmodel.TaskWhere.UserID.EQ(user.ID)).All(ctx, tx)

	if err != nil {
		panic(err)
	}

	results := make([]*model.Task, len(tasks))

	for i, task := range tasks {
		results[i] = &model.Task{
			ID:      fmt.Sprint(task.ID),
			Title:   task.Title.String,
			Content: task.Content.String,
			User: &model.User{
				ID:   fmt.Sprint(user.ID),
				Name: user.Name,
			},
		}
	}

	return results, nil
}

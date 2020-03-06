// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"fmt"

	"github.com/moutend/gqlgen-todoapp/internal/graph/generated"
	"github.com/moutend/gqlgen-todoapp/internal/graph/model"
)

func (r *mutationResolver) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	return r.createTask(ctx, input)
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	return r.createUser(ctx, input)
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	return r.login(ctx, input)
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	return r.refreshToken(ctx, input)
}

func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {
	return r.tasks(ctx)
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

package graphql

import "context"

type Resolver struct{}

func (r *queryResolver) Posts(ctx context.Context) ([]*Post, error) {
	// TODO: Implement
	return nil, nil
}

func (r *queryResolver) Post(ctx context.Context) (*Post, error) {
	// TODO: Implement
	return nil, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, title string, body string) (*Post, error) {
	// TODO: Implement
	return nil, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	// TODO: Implement
	return false, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, postID string, body string, parentID string) (*Comment, error) {
	// TODO: Implement
	return nil, nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (bool, error) {
	// TODO: Implement
	return false, nil
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

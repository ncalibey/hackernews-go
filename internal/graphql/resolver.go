package graphql

import (
	"context"
	"fmt"

	"github.com/ncalibey/hackernews-go/internal/models"
)

var (
	links = []*models.Link{
		{
			ID:          "link-0",
			URL:         "www.howtographql.com",
			Description: "Fullstack tutorial for GraphQL",
		},
	}
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Info(ctx context.Context) (string, error) {
	return "This is the API of a Hackernews Clone", nil
}

func (r *queryResolver) Feed(ctx context.Context) ([]*models.Link, error) {
	return links, nil
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Post(ctx context.Context, url, description string) (*models.Link, error) {
	idCount := len(links)
	link := &models.Link{
		ID:          fmt.Sprintf("link-%d", idCount),
		URL:         url,
		Description: description,
	}
	links = append(links, link)

	return link, nil
}

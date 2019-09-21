package graphql

import (
	"context"

	"github.com/ncalibey/hackernews-go/internal/prisma"
)

type Resolver struct {
	Prisma *prisma.Client
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Info(ctx context.Context) (string, error) {
	return "This is the API of a Hackernews Clone", nil
}

func (r *queryResolver) Feed(ctx context.Context) ([]*prisma.Link, error) {
	links, err := r.Prisma.Links(&prisma.LinksParams{}).Exec(ctx)
	var pLinks []*prisma.Link

	for _, link := range links {
		nlink := &prisma.Link{
			ID:          link.ID,
			CreatedAt:   link.CreatedAt,
			Description: link.Description,
			Url:         link.Url,
		}
		pLinks = append(pLinks, nlink)
	}

	return pLinks, err
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Post(ctx context.Context, url, description string) (*prisma.Link, error) {
	link, err := r.Prisma.CreateLink(prisma.LinkCreateInput{
		Url:         url,
		Description: description,
	}).Exec(ctx)
	return link, err
}

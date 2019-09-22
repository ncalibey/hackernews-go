package graphql

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/ncalibey/hackernews-go/internal/prisma"
	"github.com/ncalibey/hackernews-go/internal/token"
)

type Resolver struct {
	Prisma *prisma.Client
}

func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) User() UserResolver         { return &userResolver{r} }
func (r *Resolver) Link() LinkResolver         { return &linkResolver{r} }

//////////////////////////////////////////////////////////////////////////////////////////
//// Query Resolver //////////////////////////////////////////////////////////////////////

type queryResolver struct{ *Resolver }

func (r *queryResolver) Info(ctx context.Context) (string, error) {
	return "This is the API of a Hackernews Clone", nil
}

func (r *queryResolver) Feed(ctx context.Context) ([]*prisma.Link, error) {
	links, err := r.Prisma.Links(&prisma.LinksParams{}).Exec(ctx)
	return convertLinks(links), err
}

//////////////////////////////////////////////////////////////////////////////////////////
//// Mutation Resolver ///////////////////////////////////////////////////////////////////

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Post(ctx context.Context, url, description string) (*prisma.Link, error) {
	auth, ok := ctx.Value("authToken").(string)
	if !ok {
		return nil, fmt.Errorf("error getting token")
	}
	userID, err := token.GetUserID(auth)
	if err != nil {
		return nil, err
	}
	link, err := r.Prisma.CreateLink(prisma.LinkCreateInput{
		Url:         url,
		Description: description,
		PostedBy: &prisma.UserCreateOneWithoutLinksInput{
			Connect: &prisma.UserWhereUniqueInput{
				ID: &userID,
			},
		},
	}).Exec(ctx)
	return link, err
}

func (r *mutationResolver) Signup(ctx context.Context, email, password, name string) (*AuthPayload, error) {
	pw, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	user, err := r.Prisma.CreateUser(prisma.UserCreateInput{
		Name:     name,
		Email:    email,
		Password: pw,
	}).Exec(ctx)
	if err != nil {
		return nil, err
	}
	tokenString, err := token.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		User:  user,
		Token: &tokenString,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*AuthPayload, error) {
	user, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		Email: &email,
	}).Exec(ctx)
	if err != nil {
		return nil, err
	}
	if valid := CheckPasswordHash(password, user.Password); !valid {
		return nil, errors.New("invalid password")
	}
	tokenString, err := token.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		User:  user,
		Token: &tokenString,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////////////////
//// User Resolver ///////////////////////////////////////////////////////////////////////

type userResolver struct{ *Resolver }

func (r *userResolver) Links(ctx context.Context, obj *prisma.User) ([]*prisma.Link, error) {
	var id *string
	*id = obj.ID
	links, err := r.Prisma.User(prisma.UserWhereUniqueInput{ID: id}).Links(nil).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return convertLinks(links), nil
}

//////////////////////////////////////////////////////////////////////////////////////////
//// Link Resolver ///////////////////////////////////////////////////////////////////////

type linkResolver struct{ *Resolver }

func (r *linkResolver) PostedBy(ctx context.Context, obj *prisma.Link) (*prisma.User, error) {
	var id *string
	*id = obj.ID
	user, err := r.Prisma.Link(prisma.LinkWhereUniqueInput{ID: id}).PostedBy().Exec(nil)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//////////////////////////////////////////////////////////////////////////////////////////
//// Helper Functions ////////////////////////////////////////////////////////////////////

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func convertLinks(links []prisma.Link) []*prisma.Link {
	var pointerLinks []*prisma.Link
	for _, link := range links {
		l := &prisma.Link{
			ID:          link.ID,
			CreatedAt:   link.CreatedAt,
			Description: link.Description,
			Url:         link.Url,
		}
		pointerLinks = append(pointerLinks, l)
	}
	return pointerLinks
}

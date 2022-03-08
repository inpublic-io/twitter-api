package grpc

import (
	"context"
	"time"

	twitterv1 "github.com/inpublic-io/inpublicapis/twitter/v1"
	"github.com/inpublic-io/twitter-api/clients"
	"github.com/vniche/twitter-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

const inpublicTwitterID = "1494121132508733447"

// usersServer is used to implement io.inpublic.twitter.v1.Users
type usersServer struct {
	twitterv1.UnimplementedUsersServer
}

func (s *usersServer) LookupInpublic(ctx context.Context, _ *emptypb.Empty) (*twitterv1.User, error) {
	twitterClient := clients.TwitterClient()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	user, err := twitterClient.LookupUserByID(ctx, inpublicTwitterID, map[string][]string{
		twitter.UserGenericQueryParameters.UserFields: {
			twitter.UserFields.PublicMetrics,
		},
	})
	if err != nil {
		return nil, err
	}

	return &twitterv1.User{
		Id:             user.ID,
		Name:           user.Name,
		Username:       user.Username,
		FollowersCount: int32(user.PublicMetrics.FollowersCount),
	}, nil
}

package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
)

// Stream is the resolver for the stream field.
func (r *queryResolver) Stream(ctx context.Context) (*gqlmodel.Stream, error) {
	panic(fmt.Errorf("not implemented: Stream - stream"))
}

// StreamInfo is the resolver for the streamInfo field.
func (r *subscriptionResolver) StreamInfo(ctx context.Context) (<-chan *gqlmodel.Stream, error) {
	r.streamViewers.Inc()

	userID, userIdErr := r.sessionStorage.GetUserID(ctx)
	if userIdErr != nil {
		return nil, userIdErr
	}

	if userIdErr == nil {
		user, err := r.userRepo.FindByID(ctx, uuid.MustParse(userID))
		if err != nil {
			return nil, err
		}

		chattersLock.Lock()
		r.streamChatters[user.ID.String()] = gqlmodel.Chatter{
			User: &gqlmodel.User{
				ID:          user.ID.String(),
				Name:        user.Name,
				DisplayName: user.DisplayName,
				Color:       user.Color,
				Roles:       nil,
				IsBanned:    user.Banned,
				CreatedAt:   user.CreatedAt,
				AvatarURL:   user.AvatarUrl,
			},
		}

		chattersLock.Unlock()
	}

	chann := make(chan *gqlmodel.Stream)

	go func() {
		for {
			select {
			case <-ctx.Done():
				chattersLock.Lock()
				defer chattersLock.Unlock()

				r.streamViewers.Dec()

				if userIdErr == nil {
					delete(r.streamChatters, userID)
				}

				close(chann)
				return
			default:
				chatters := lo.Values(r.streamChatters)
				slices.SortFunc(
					chatters,
					func(a, b gqlmodel.Chatter) int {
						return strings.Compare(a.User.Name, b.User.Name)
					},
				)

				chann <- &gqlmodel.Stream{
					Viewers:  int(r.streamViewers.Load()),
					Chatters: chatters,
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return chann, nil
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
var chattersLock = sync.Mutex{}

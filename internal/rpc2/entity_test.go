package rpc2

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/token/null"

	types "github.com/NetAuth/Protocol"
	pb "github.com/NetAuth/Protocol/v2"
)

func TestEntityCreate(t *testing.T) {
	cases := []struct {
		req      pb.EntityRequest
		wantErr  error
		readonly bool
	}{
		{
			// Works, entity is created.
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("test1"),
				},
			},
			wantErr:  nil,
			readonly: false,
		},
		{
			// Fails, server is in read-only mode
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("test1"),
				},
			},
			wantErr:  ErrReadOnly,
			readonly: true,
		},
		{
			// Fails, token is invalid
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.InvalidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("test1"),
				},
			},
			wantErr:  ErrUnauthenticated,
			readonly: false,
		},
		{
			// Fails, token lacks capabilities
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidEmptyToken,
				},
				Entity: &types.Entity{
					ID: proto.String("test1"),
				},
			},
			wantErr:  ErrRequestorUnqualified,
			readonly: false,
		},
		{
			// Fails, duplicate resource
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					// This gets created by
					// initTree which fills in the
					// tree for testing purposes.
					ID: proto.String("entity1"),
				},
			},
			wantErr:  ErrExists,
			readonly: false,
		},
		{
			// Fails, internal write error
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Entity: &types.Entity{
					ID: proto.String("save-error"),
				},
			},
			wantErr:  ErrInternal,
			readonly: false,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		initTree(t, s)
		s.readonly = c.readonly
		if _, err := s.EntityCreate(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestEntityUpdate(t *testing.T) {
	cases := []struct {
		req      pb.EntityRequest
		readonly bool
		wantErr  error
	}{
		{
			// Works, will change the metadata DisplayName
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Data: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: false,
			wantErr:  nil,
		},
		{
			// Fails, server is in read-only mode
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Data: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: true,
			wantErr:  ErrReadOnly,
		},
		{
			// Fails, token is invalid
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.InvalidToken,
				},
				Data: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: false,
			wantErr:  ErrUnauthenticated,
		},
		{
			// Fails, token has no capabilities
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidEmptyToken,
				},
				Data: &types.Entity{
					ID: proto.String("entity1"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: false,
			wantErr:  ErrRequestorUnqualified,
		},
		{
			// Fails, entity does not exist
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Data: &types.Entity{
					ID: proto.String("does-not-exist"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: false,
			wantErr:  ErrDoesNotExist,
		},
		{
			// Fails, db write failure
			req: pb.EntityRequest{
				Auth: &pb.AuthData{
					Token: &null.ValidToken,
				},
				Data: &types.Entity{
					ID: proto.String("load-error"),
					Meta: &types.EntityMeta{
						DisplayName: proto.String("First Entity"),
					},
				},
			},
			readonly: false,
			wantErr:  ErrInternal,
		},
	}

	for i, c := range cases {
		s := newServer(t)
		s.readonly = c.readonly
		initTree(t, s)
		if _, err := s.EntityUpdate(context.Background(), &c.req); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
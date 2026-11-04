package command

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/product"
	"webtranshome/internal/modules/catalog/domain/room"
)

type assignmentProductRepoStub struct {
	product.Product
	findErr error
}

func (r *assignmentProductRepoStub) ListProducts(context.Context) ([]product.Product, error) { return nil, nil }
func (r *assignmentProductRepoStub) FindProduct(context.Context, int64) (product.Product, error) {
	return product.Product{ID: 7, Name: "Lamp", Slug: "lamp", CategoryID: 1, PriceGeneral: productFloat64Pointer(1000), UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA"}, r.findErr
}
func (r *assignmentProductRepoStub) CreateProduct(context.Context, product.Product) (product.Product, error) { return product.Product{}, nil }
func (r *assignmentProductRepoStub) UpdateProduct(context.Context, product.Product) error { return nil }
func (r *assignmentProductRepoStub) ReplaceProductRooms(context.Context, int64, []int64) error { return nil }

type assignmentRoomRepoStub struct {
	findErr error
}

func (r *assignmentRoomRepoStub) ListRooms(context.Context) ([]room.Room, error) { return nil, nil }
func (r *assignmentRoomRepoStub) FindRoom(context.Context, int64) (room.Room, error) {
	return room.Room{ID: 1, Name: "Living", Slug: "living", ImageCover: "/living.png", SortOrder: 1, IsActive: true}, r.findErr
}
func (r *assignmentRoomRepoStub) CreateRoom(context.Context, room.Room) (room.Room, error) { return room.Room{}, nil }
func (r *assignmentRoomRepoStub) UpdateRoom(context.Context, room.Room) error { return nil }

func TestReplaceProductRooms(t *testing.T) {
	productRepo := &assignmentProductRepoStub{}
	roomRepo := &assignmentRoomRepoStub{}
	if err := NewReplaceProductRooms(productRepo, roomRepo, productRepo).Execute(context.Background(), ReplaceProductRoomsInput{ProductID: 7, RoomIDs: []int64{1}}); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestReplaceProductRoomsRejectsInvalidAssignment(t *testing.T) {
	err := NewReplaceProductRooms(&assignmentProductRepoStub{}, &assignmentRoomRepoStub{}, &assignmentProductRepoStub{}).Execute(context.Background(), ReplaceProductRoomsInput{ProductID: 0, RoomIDs: []int64{1}})
	if !errors.Is(err, product.ErrInvalidProductID) {
		t.Fatalf("Execute() error = %v, want %v", err, product.ErrInvalidProductID)
	}
}

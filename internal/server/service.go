package server

import (
	"context"
	"database/sql"
	"io"
	"log"
	"path/filepath"

	"grpc-go/internal/account"
	"grpc-go/internal/orders"
	"grpc-go/internal/preferences"
	"grpc-go/internal/roles"
	"grpc-go/pkg/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	db *sql.DB
	proto.UnimplementedMyServiceServer
	accountService    *account.Service
	preferenceService *preferences.Service
	roleService       *roles.Service
	orderService      *orders.Service
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db:                db,
		accountService:    account.NewService(db),
		preferenceService: preferences.NewService(db),
		roleService:       roles.NewService(db),
		orderService:      orders.NewService(db),
	}
}

// --- Account Methods ---
func (s *Service) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.AccountResponse, error) {
	acc, err := s.accountService.CreateAccount(req.Username, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.AccountResponse{
		Id:       acc.ID,
		Username: acc.Username,
		Email:    acc.Email,
		Status:   acc.Status,
	}, nil
}

func (s *Service) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.AccountResponse, error) {
	acc, err := s.accountService.GetAccount(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &proto.AccountResponse{
		Id:       acc.ID,
		Username: acc.Username,
		Email:    acc.Email,
		Status:   acc.Status,
	}, nil
}

func (s *Service) UpdateAccount(ctx context.Context, req *proto.UpdateAccountRequest) (*proto.AccountResponse, error) {
	acc, err := s.accountService.UpdateAccount(req.Id, req.Username, req.Email, req.Password, req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.AccountResponse{
		Id:       acc.ID,
		Username: acc.Username,
		Email:    acc.Email,
		Status:   acc.Status,
	}, nil
}

func (s *Service) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	if err := s.accountService.DeleteAccount(req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.DeleteAccountResponse{Success: true}, nil
}

// --- Preference Methods ---
func (s *Service) CreatePreference(ctx context.Context, req *proto.CreatePreferenceRequest) (*proto.PreferenceResponse, error) {
	pref, err := s.preferenceService.CreatePreference(req.AccountId, req.Theme, req.Notifications, req.Locale)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.PreferenceResponse{
		Id:            pref.ID,
		AccountId:     pref.AccountID,
		Theme:         pref.Theme,
		Notifications: pref.Notifications,
		Locale:        pref.Locale,
	}, nil
}

func (s *Service) GetPreference(ctx context.Context, req *proto.GetPreferenceRequest) (*proto.PreferenceResponse, error) {
	pref, err := s.preferenceService.GetPreference(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &proto.PreferenceResponse{
		Id:            pref.ID,
		AccountId:     pref.AccountID,
		Theme:         pref.Theme,
		Notifications: pref.Notifications,
		Locale:        pref.Locale,
	}, nil
}

func (s *Service) UpdatePreference(ctx context.Context, req *proto.UpdatePreferenceRequest) (*proto.PreferenceResponse, error) {
	pref, err := s.preferenceService.UpdatePreference(req.Id, req.Theme, req.Notifications, req.Locale)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.PreferenceResponse{
		Id:            pref.ID,
		AccountId:     pref.AccountID,
		Theme:         pref.Theme,
		Notifications: pref.Notifications,
		Locale:        pref.Locale,
	}, nil
}

func (s *Service) DeletePreference(ctx context.Context, req *proto.DeletePreferenceRequest) (*proto.DeletePreferenceResponse, error) {
	if err := s.preferenceService.DeletePreference(req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.DeletePreferenceResponse{Success: true}, nil
}

// --- Role Methods ---
func (s *Service) CreateRole(ctx context.Context, req *proto.CreateRoleRequest) (*proto.RoleResponse, error) {
	role, err := s.roleService.CreateRole(req.Name, req.Permissions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.RoleResponse{
		Id:          role.ID,
		Name:        role.Name,
		Permissions: role.Permissions,
	}, nil
}

func (s *Service) GetRole(ctx context.Context, req *proto.GetRoleRequest) (*proto.RoleResponse, error) {
	role, err := s.roleService.GetRole(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &proto.RoleResponse{
		Id:          role.ID,
		Name:        role.Name,
		Permissions: role.Permissions,
	}, nil
}

func (s *Service) UpdateRole(ctx context.Context, req *proto.UpdateRoleRequest) (*proto.RoleResponse, error) {
	role, err := s.roleService.UpdateRole(req.Id, req.Name, req.Permissions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.RoleResponse{
		Id:          role.ID,
		Name:        role.Name,
		Permissions: role.Permissions,
	}, nil
}

func (s *Service) DeleteRole(ctx context.Context, req *proto.DeleteRoleRequest) (*proto.DeleteRoleResponse, error) {
	if err := s.roleService.DeleteRole(req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.DeleteRoleResponse{Success: true}, nil
}

// --- Order Methods ---
func (s *Service) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.OrderResponse, error) {
	order, err := s.orderService.CreateOrder(req.AccountId, convertOrderItems(req.Items))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.OrderResponse{
		Id:        order.ID,
		AccountId: order.AccountID,
		Items:     convertToProtoOrderItems(order.Items),
		Total:     order.Total,
		Status:    order.Status,
	}, nil
}

func (s *Service) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.OrderResponse, error) {
	order, err := s.orderService.GetOrder(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &proto.OrderResponse{
		Id:        order.ID,
		AccountId: order.AccountID,
		Items:     convertToProtoOrderItems(order.Items),
		Total:     order.Total,
		Status:    order.Status,
	}, nil
}

func (s *Service) UpdateOrder(ctx context.Context, req *proto.UpdateOrderRequest) (*proto.OrderResponse, error) {
	order, err := s.orderService.UpdateOrder(req.Id, convertOrderItems(req.Items), req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.OrderResponse{
		Id:        order.ID,
		AccountId: order.AccountID,
		Items:     convertToProtoOrderItems(order.Items),
		Total:     order.Total,
		Status:    order.Status,
	}, nil
}

func (s *Service) DeleteOrder(ctx context.Context, req *proto.DeleteOrderRequest) (*proto.DeleteOrderResponse, error) {
	if err := s.orderService.DeleteOrder(req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.DeleteOrderResponse{Success: true}, nil
}

func (s *Service) UploadFile(stream proto.MyService_UploadFileServer) error {
	var filename string
	var fileData []byte
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		if filename == "" && req.Filename != "" {
			filename = req.Filename
			ext := filepath.Ext(filename)

			if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".pdf" {
				return status.Error(codes.InvalidArgument, "unsupported file type")
			}
		}
		fileData = append(fileData, req.ChunkData...)
		if len(req.ChunkData) > 5*1024*1024 {
			return status.Error(codes.InvalidArgument, "chunk size exceeds limit")
		}
	}
	log.Printf("Received file %s, size %d bytes", filename, len(fileData))

	return stream.SendAndClose(&proto.UploadFileResponse{
		Success: true,
		Message: "File uploaded successfully",
	})
}

func convertOrderItems(items []*proto.OrderItem) []orders.OrderItem {
	var orderItems []orders.OrderItem
	for _, item := range items {
		orderItems = append(orderItems, orders.OrderItem{
			Product:  item.Product,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}
	return orderItems
}

func convertToProtoOrderItems(items []orders.OrderItem) []*proto.OrderItem {
	var protoItems []*proto.OrderItem
	for _, item := range items {
		protoItems = append(protoItems, &proto.OrderItem{
			Product:  item.Product,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}
	return protoItems
}

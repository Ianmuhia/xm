package main

import (
	"context"
	pb "xmserver/gen"
)

type CompanyHandler struct {
	pb.UnimplementedCompanyServiceServer
	usersService   UsersService
	companyService CompanyService
	tokenService   TokenService
}

func NewCompanyHandler(
	usersService UsersService,
	companyService CompanyService,
	tokenService TokenService,
) *CompanyHandler {
	return &CompanyHandler{
		usersService:   usersService,
		companyService: companyService,
		tokenService:   tokenService,
	}
}

func (h *CompanyHandler) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	data, err := h.companyService.CreateCompany(ctx, &Company{
		Name:        req.Company.Name,
		Employees:   int(req.Company.GetEmployees()),
		Description: req.Company.GetDescription(),
		Registered:  req.Company.GetRegistered(),
		Type:        req.Company.GetType().String(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{
		Company: &pb.Company{
			Name:        data.Name,
			Employees:   int32(data.Employees),
			Description: data.Description,
			Registered:  data.Registered,
			Type:        pb.CompanyType_COOPERATIVE,
		},
	}, nil
}
func (h *CompanyHandler) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := h.companyService.DeleteCompany(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{
		Success: true,
	}, nil
}
func (h *CompanyHandler) GetOne(ctx context.Context, req *pb.GetOneRequest) (*pb.GetOneResponse, error) {
	data, err := h.companyService.GetCompany(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.GetOneResponse{
		Company: &pb.Company{
			Name:        data.Name,
			Employees:   int32(data.Employees),
			Description: data.Description,
			Registered:  data.Registered,
			Type:        pb.CompanyType_COOPERATIVE,
		},
	}, nil
}
func (h *CompanyHandler) Patch(ctx context.Context, req *pb.PatchRequest) (*pb.PatchResponse, error) {
	data, err := h.companyService.UpdateCompany(ctx, &Company{
		Name:        req.Company.Name,
		Employees:   int(req.Company.GetEmployees()),
		Description: req.Company.GetDescription(),
		Registered:  req.Company.GetRegistered(),
		Type:        req.Company.GetType().String(),
		ID:          req.Company.Id,
	})
	if err != nil {
		return nil, err
	}
	return &pb.PatchResponse{
		Company: &pb.Company{
			Name:        data.Name,
			Employees:   int32(data.Employees),
			Description: data.Description,
			Registered:  data.Registered,
			Type:        pb.CompanyType_COOPERATIVE,
		},
	}, nil
}
func (h *CompanyHandler) Register(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	data, err := h.usersService.CreateUser(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	accessToken, err := h.tokenService.Create(&User{UserName: data.UserName})
	if err != nil {
		return nil, err
	}

	refreshToken, err := h.tokenService.CreateRefresh(&User{UserName: data.UserName})
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *CompanyHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.CreateUserResponse, error) {
	data, err := h.usersService.FindUser(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	accessToken, err := h.tokenService.Create(&User{UserName: data.UserName})
	if err != nil {
		return nil, err
	}

	refreshToken, err := h.tokenService.CreateRefresh(&User{UserName: data.UserName})
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

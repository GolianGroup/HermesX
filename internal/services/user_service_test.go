package services

import (
	dto "golang_template/handler/dtos"
	"golang_template/internal/ent"
	"golang_template/internal/mocks"
	"golang_template/internal/repositories"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
)

func Test_userService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	user := ent.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: "testpassword",
	}

	type fields struct {
		repo repositories.UserRepository
	}
	type args struct {
		ctx  *fiber.Ctx
		user dto.User
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		setupMocks func()
		want       *ent.User
		wantErr    error
	}{
		{
			name: "Successful login",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				ctx:  createFiberContext(),
				user: dto.User{},
			},
			setupMocks: func() {
				mockRepo.EXPECT().
					Get(gomock.Any(), &ent.User{
						Username: "testuser",
						Password: "testpassword",
					}).Return(&user, nil)
			},
			want:    &user,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				repo: tt.fields.repo,
			}
			s.Login(tt.args.ctx, tt.args.user)
		})
	}
}

func createFiberContext() *fiber.Ctx {
	app := fiber.New()
	req := fasthttp.RequestCtx{}
	return app.AcquireCtx(&req)
}

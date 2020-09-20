package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/avvarikrish/chefcurrygobbles/pkg/input"

	opb "github.com/avvarikrish/chefcurrygobbles/proto/orders_server"
	rpb "github.com/avvarikrish/chefcurrygobbles/proto/restaurant_server"
	pb "github.com/avvarikrish/chefcurrygobbles/proto/users_server"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Address   Addr
}

type LoginUserRequest struct {
	Email    string
	Password string
}

type UpdateUserRequest struct {
	OldEmail string
	NewUser  User
}

type DeleteUserRequest struct {
	Email string
}

type Restaurant struct {
	Email   string
	Phone   string
	Name    string
	Address Addr
	Menu    []MenuItem
}

type MenuItem struct {
	Name  string
	Price string
}

type CreateOrderRequest struct {
	Email     string
	RestEmail string
	RestPhone string
	OrderItem []oItem
}

type oItem struct {
	MenuId   string
	Quantity string
}

type Addr struct {
	StreetNumber string
	Street       string
	City         string
	State        string
	Zip          string
}

func main() {
	log.Println("Hello I'm a ccgobbles client")

	input_func := os.Args[1]
	func_to_run, port := whichFunction(input_func)
	if func_to_run == nil {
		log.Fatalf("Invalid function: %v", input_func)
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	cc, err := grpc.Dial("localhost:"+port, opts...)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	func_to_run(cc)
}

func whichFunction(func_to_run string) (func(*grpc.ClientConn), string) {
	switch func_to_run {
	case "register_user":
		return registerUser, "50051"

	case "login_user":
		return loginUser, "50051"

	case "update_user":
		return updateUser, "50051"

	case "delete_user":
		return deleteUser, "50051"

	case "add_restaurant":
		return addRestaurant, "50052"

	case "create_order":
		return createOrder, "50053"

	default:
		return nil, ""
	}
}

func registerUser(cc *grpc.ClientConn) {
	log.Println("Registering user")

	c := pb.NewUsersClient(cc)
	t := reflect.TypeOf(User{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	user := v.Interface().(*User)

	req := &pb.RegisterUserRequest{
		User: &pb.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Address: &pb.Address{
				StreetNumber: user.Address.StreetNumber,
				Street:       user.Address.Street,
				City:         user.Address.City,
				State:        user.Address.State,
				Zip:          user.Address.Zip,
			},
		},
	}
	res, err := c.RegisterUser(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			if s.Code() == codes.AlreadyExists {
				log.Println("Email already exists")
			}
			log.Println(s.Details(), s.Message())
			log.Fatalf("RPC error: %v", s.Code())
		}
		log.Fatalf("Unexpected error: %v\n", err)
	}
	log.Printf("New user created: %v\n", res)
}

func loginUser(cc *grpc.ClientConn) {
	log.Println("Attempting login")

	c := pb.NewUsersClient(cc)
	t := reflect.TypeOf(LoginUserRequest{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	loginRequest := v.Interface().(*LoginUserRequest)

	req := &pb.LoginUserRequest{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}
	res, err := c.LoginUser(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			// grpc error
			if s.Code() == codes.NotFound {
				log.Println("Email does not exist")
				return
			}
		}
		log.Fatalf("Unexpected error: %v\n", err)
	}
	log.Println(res.GetResponse())
}

func updateUser(cc *grpc.ClientConn) {
	log.Println("Attempting update user")

	c := pb.NewUsersClient(cc)
	t := reflect.TypeOf(UpdateUserRequest{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	updateRequest := v.Interface().(*UpdateUserRequest)

	req := &pb.UpdateUserRequest{
		OldEmail: updateRequest.OldEmail,
		User: &pb.User{
			FirstName: updateRequest.NewUser.FirstName,
			LastName:  updateRequest.NewUser.LastName,
			Email:     updateRequest.NewUser.Email,
			Password:  updateRequest.NewUser.Password,
			Address: &pb.Address{
				StreetNumber: updateRequest.NewUser.Address.StreetNumber,
				Street:       updateRequest.NewUser.Address.Street,
				City:         updateRequest.NewUser.Address.City,
				State:        updateRequest.NewUser.Address.State,
				Zip:          updateRequest.NewUser.Address.Zip,
			},
		},
	}
	res, err := c.UpdateUser(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			// grpc error
			if s.Code() == codes.Internal {
				log.Printf("Internal error: %v\n", err)
				return
			} else if s.Code() == codes.NotFound {
				log.Println("Email not found")
				return
			}
		}
		log.Fatalf("Unexpected error: %v\n", err)
	}
	log.Println(res.GetResponse())
}

func deleteUser(cc *grpc.ClientConn) {
	log.Println("Deleting User")

	c := pb.NewUsersClient(cc)
	t := reflect.TypeOf(DeleteUserRequest{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	deleteUserRequest := v.Interface().(*DeleteUserRequest)

	req := &pb.DeleteUserRequest{
		Email: deleteUserRequest.Email,
	}
	res, err := c.DeleteUser(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			if s.Code() == codes.Internal {
				log.Printf("Internal error: %v\n", err)
				return
			} else if s.Code() == codes.NotFound {
				log.Println("Email not found")
				return
			}
		}
		log.Fatalf("Unexpected error: %v\n", err)
	}
	log.Println(res.GetResponse())
}

func addRestaurant(cc *grpc.ClientConn) {
	log.Println("Adding restaurant")

	c := rpb.NewRestaurantsClient(cc)
	t := reflect.TypeOf(Restaurant{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	restaurant := v.Interface().(*Restaurant)

	menu := []*rpb.MenuItem{}
	for _, m := range restaurant.Menu {
		f, _ := strconv.ParseFloat(m.Price, 64)
		fmt.Println(f)
		menu = append(menu, &rpb.MenuItem{
			Name:  m.Name,
			Price: f,
		})
	}

	req := &rpb.AddRestaurantRequest{
		Restaurant: &rpb.Restaurant{
			Phone: restaurant.Phone,
			Email: restaurant.Email,
			Name:  restaurant.Name,
			Address: &rpb.Address{
				StreetNumber: restaurant.Address.StreetNumber,
				Street:       restaurant.Address.Street,
				City:         restaurant.Address.City,
				State:        restaurant.Address.State,
				Zip:          restaurant.Address.Zip,
			},
			Menuitem: menu,
		},
	}
	res, err := c.AddRestaurant(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			if s.Code() == codes.Internal {
				log.Printf("Internal error: %v\n", err)
				return
			} else if s.Code() == codes.AlreadyExists {
				log.Println("Restaurant already exists")
				return
			}
			log.Fatalf("Mongodb Error: %v\n", err)
		}
		log.Fatalf("Unexpected Error: %v\n", err)
	}
	log.Println(res.GetResponse())
}

func createOrder(cc *grpc.ClientConn) {
	log.Println("Creating order")

	c := opb.NewOrdersClient(cc)
	t := reflect.TypeOf(CreateOrderRequest{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	order := v.Interface().(*CreateOrderRequest)

	orderItems := []*opb.OrderItem{}
	for _, o := range order.OrderItem {
		m, _ := strconv.ParseInt(o.MenuId, 10, 64)
		i, _ := strconv.ParseInt(o.Quantity, 10, 64)
		orderItems = append(orderItems, &opb.OrderItem{
			MenuId:   m,
			Quantity: i,
		})
	}

	req := &opb.CreateOrderRequest{
		Order: &opb.Order{
			Email:     order.Email,
			RestPhone: order.RestPhone,
			RestEmail: order.RestEmail,
			OrderItem: orderItems,
		},
	}
	res, err := c.CreateOrder(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			// grpc error
			if s.Code() == codes.Internal {
				log.Printf("Internal error: %v\n", err)
				return
			} else if s.Code() == codes.NotFound {
				log.Println("Restaurant or User does not exist", err)
				return
			}
			log.Println(err.Error())
			log.Fatalf("MongoDB Error: %v\n", err)
		}
		log.Fatalf("Unexpected Error: %v\n", err)
	}
	log.Println(res.GetResponse())
}

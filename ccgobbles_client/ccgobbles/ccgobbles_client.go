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
	"google.golang.org/protobuf/proto"

	"github.com/avvarikrish/chefcurrygobbles/pkg/input"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/ccgobbles_server"
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
	Address  Addr
}

type DeleteUserRequest struct {
	Email string
}

type Restaurant struct {
	RestId  string
	Name    string
	Address Addr
	Menu    []MenuItem
}

type MenuItem struct {
	MenuId string
	Name   string
	Price  string
}

type CreateOrderRequest struct {
	Email     string
	RestId    string
	OrderId   string
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

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	cc, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := pb.NewCCGobblesClient(cc)

	input_func := os.Args[1]
	func_to_run := whichFunction(input_func)
	if func_to_run == nil {
		log.Fatalf("Invalid function: %v", input_func)
	}

	func_to_run(c)
}

func whichFunction(func_to_run string) func(pb.CCGobblesClient) {
	switch func_to_run {
	case "register_user":
		return registerUser

	case "login_user":
		return loginUser

	case "update_user":
		return updateUser

	case "delete_user":
		return deleteUser

	case "add_restaurant":
		return addRestaurant

	case "create_order":
		return createOrder

	default:
		return nil
	}
}

func registerUser(c pb.CCGobblesClient) {
	log.Println("Registering user")

	t := reflect.TypeOf(User{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	user := v.Interface().(*User)

	req := &pb.RegisterUserRequest{
		User: &pb.User{
			FirstName: proto.String(user.FirstName),
			LastName:  proto.String(user.LastName),
			Email:     proto.String(user.Email),
			Password:  proto.String(user.Password),
			Address: &pb.Address{
				StreetNumber: proto.String(user.Address.StreetNumber),
				Street:       proto.String(user.Address.Street),
				City:         proto.String(user.Address.City),
				State:        proto.String(user.Address.State),
				Zip:          proto.String(user.Address.Zip),
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

func loginUser(c pb.CCGobblesClient) {
	log.Println("Attempting login")

	t := reflect.TypeOf(LoginUserRequest{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	loginRequest := v.Interface().(*LoginUserRequest)

	req := &pb.LoginUserRequest{
		Email:    proto.String(loginRequest.Email),
		Password: proto.String(loginRequest.Password),
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

func updateUser(c pb.CCGobblesClient) {
	log.Println("Attempting update user")

	t := reflect.TypeOf(User{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	user := v.Interface().(*User)

	req := &pb.UpdateUserRequest{
		User: &pb.User{
			FirstName: proto.String(user.FirstName),
			LastName:  proto.String(user.LastName),
			Email:     proto.String(user.Email),
			Password:  proto.String(user.Password),
			Address: &pb.Address{
				StreetNumber: proto.String(user.Address.StreetNumber),
				Street:       proto.String(user.Address.Street),
				City:         proto.String(user.Address.City),
				State:        proto.String(user.Address.State),
				Zip:          proto.String(user.Address.Zip),
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

func deleteUser(c pb.CCGobblesClient) {
	log.Println("Deleting User")

	t := reflect.TypeOf(DeleteUserRequest{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	deleteUserRequest := v.Interface().(*DeleteUserRequest)

	req := &pb.DeleteUserRequest{
		Email: proto.String(deleteUserRequest.Email),
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

func addRestaurant(c pb.CCGobblesClient) {
	log.Println("Adding restaurant")

	t := reflect.TypeOf(Restaurant{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	restaurant := v.Interface().(*Restaurant)

	menu := []*pb.MenuItem{}
	for _, m := range restaurant.Menu {
		f, _ := strconv.ParseFloat(m.Price, 64)
		fmt.Println(f)
		menu = append(menu, &pb.MenuItem{
			MenuId: proto.String(m.MenuId),
			Name:   proto.String(m.Name),
			Price:  proto.Float64(f),
		})
	}

	req := &pb.AddRestaurantRequest{
		Restaurant: &pb.Restaurant{
			RestId: proto.String(restaurant.RestId),
			Name:   proto.String(restaurant.Name),
			Address: &pb.Address{
				StreetNumber: proto.String(restaurant.Address.StreetNumber),
				Street:       proto.String(restaurant.Address.Street),
				City:         proto.String(restaurant.Address.City),
				State:        proto.String(restaurant.Address.State),
				Zip:          proto.String(restaurant.Address.Zip),
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

func createOrder(c pb.CCGobblesClient) {
	log.Println("Creating order")

	t := reflect.TypeOf(CreateOrderRequest{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	order := v.Interface().(*CreateOrderRequest)

	orderItems := []*pb.OrderItem{}
	for _, o := range order.OrderItem {
		i, _ := strconv.ParseInt(o.Quantity, 10, 64)
		orderItems = append(orderItems, &pb.OrderItem{
			MenuId:   proto.String(o.MenuId),
			Quantity: proto.Int64(i),
		})
	}

	req := &pb.CreateOrderRequest{
		Email:     proto.String(order.Email),
		RestId:    proto.String(order.RestId),
		OrderId:   proto.String(order.OrderId),
		OrderItem: orderItems,
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

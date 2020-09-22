package bsonconversion

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	opb "github.com/avvarikrish/chefcurrygobbles/proto/orders_server"
	rpb "github.com/avvarikrish/chefcurrygobbles/proto/restaurant_server"
	pb "github.com/avvarikrish/chefcurrygobbles/proto/users_server"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Email     string             `bson:"email" json:"email"`
	Password  []byte             `bson:"password" json:"password"`
	Address   Addr               `bson:"address" json:"address"`
}

type Restaurant struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email   string             `bson:"email,omitempty" json:"email,omitempty"`
	Phone   string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Name    string             `bson:"name" json:"name"`
	Address Addr               `bson:"address" json:"address"`
	Menu    []MenuItem         `bson:"menu" json:"menu"`
}

type Order struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	RestID primitive.ObjectID `bson:"rest_id" json:"rest_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items  []OrderItem        `bson:"items" json:"items"`
}

type Addr struct {
	StreetNumber string `bson:"street_number" json:"street_number"`
	Street       string `bson:"street" json:"street"`
	City         string `bson:"city" json:"city"`
	State        string `bson:"state" json:"state"`
	Zip          string `bson:"zip" json:"zip"`
}

type MenuItem struct {
	Id    int     `bson:"menu_id" json:"menu_id"`
	Name  string  `bson:"item_name" json:"item_name"`
	Price float64 `bson:"price" json:"price"`
}

type OrderItem struct {
	MenuId   int `bson:"menu_id" json:"menu_id"`
	Quantity int `bson:"quantity" json:"quantity"`
}

// CreateUserBson creates a bson object of user data
func CreateUserBson(userReq *pb.User, password []byte) User {
	userAddReq := userReq.GetAddress()
	return User{
		FirstName: userReq.GetFirstName(),
		LastName:  userReq.GetLastName(),
		Email:     userReq.GetEmail(),
		Password:  password,
		Address: Addr{
			StreetNumber: userAddReq.GetStreetNumber(),
			Street:       userAddReq.GetStreet(),
			City:         userAddReq.GetCity(),
			State:        userAddReq.GetState(),
			Zip:          userAddReq.GetZip(),
		},
	}
}

// CreateRestBson creates a bson object of restaurant data
func CreateRestBson(resReq *rpb.Restaurant) Restaurant {
	resAddReq := resReq.GetAddress()
	menu := []MenuItem{}
	for i, m := range resReq.GetMenuitem() {
		menu = append(menu, MenuItem{
			Id:    i,
			Name:  m.GetName(),
			Price: m.GetPrice(),
		})
	}
	return Restaurant{
		Phone: resReq.GetPhone(),
		Email: resReq.GetEmail(),
		Name:  resReq.GetName(),
		Address: Addr{
			StreetNumber: resAddReq.GetStreetNumber(),
			Street:       resAddReq.GetStreet(),
			City:         resAddReq.GetCity(),
			State:        resAddReq.GetState(),
			Zip:          resAddReq.GetZip(),
		},
		Menu: menu,
	}
}

// CreateOrderBson creates a bson object of order data
func CreateOrderBson(req *opb.CreateOrderRequest, restId primitive.ObjectID, usrId primitive.ObjectID) Order {
	items := []OrderItem{}
	orderReq := req.GetOrder()

	for _, o := range orderReq.GetOrderItem() {
		items = append(items, OrderItem{
			MenuId:   int(o.GetMenuId()),
			Quantity: int(o.GetQuantity()),
		})
	}
	return Order{
		UserID: usrId,
		RestID: restId,
		Items:  items,
	}
}

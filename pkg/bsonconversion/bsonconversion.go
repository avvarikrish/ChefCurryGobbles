package bsonconversion

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/ccgobbles_server"
)

// CreateUserBson creates a bson object of user data
func CreateUserBson(userReq *pb.User) primitive.M {
	userAddReq := userReq.GetAddress()
	return bson.M{
		"first_name": userReq.GetFirstName(),
		"last_name":  userReq.GetLastName(),
		"email":      userReq.GetEmail(),
		"password":   userReq.GetPassword(),
		"address": bson.M{
			"street_number": userAddReq.GetStreetNumber(),
			"street":        userAddReq.GetStreet(),
			"city":          userAddReq.GetCity(),
			"state":         userAddReq.GetState(),
			"zip":           userAddReq.GetZip(),
		},
	}
}

// CreateRestBson creates a bson object of restaurant data
func CreateRestBson(resReq *pb.Restaurant) primitive.M {
	resAddReq := resReq.GetAddress()
	menu := bson.A{}
	for i, m := range resReq.GetMenuitem() {
		menu = append(menu, bson.M{
			"menuid": i,
			"name":   m.GetName(),
			"price":  m.GetPrice(),
		})
	}
	return bson.M{
		"phone": resReq.GetPhone(),
		"email": resReq.GetEmail(),
		"name":  resReq.GetName(),
		"address": bson.M{
			"street_number": resAddReq.GetStreetNumber(),
			"street":        resAddReq.GetStreet(),
			"city":          resAddReq.GetCity(),
			"state":         resAddReq.GetState(),
			"zip":           resAddReq.GetZip(),
		},
		"menu": menu,
	}
}

// CreateOrderBson creates a bson object of order data
func CreateOrderBson(orderReq *pb.CreateOrderRequest, dt time.Time) primitive.M {
	items := bson.A{}
	for _, o := range orderReq.GetOrderItem() {
		items = append(items, bson.M{
			"menuid":   o.GetMenuId(),
			"quantity": o.GetQuantity(),
		})
	}
	return bson.M{
		"order_id": orderReq.GetOrderId(),
		"email":    orderReq.GetEmail(),
		"rest_id":  orderReq.GetRestId(),
		"items":    items,
		"time":     dt,
	}
}

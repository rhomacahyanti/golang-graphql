package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

type User struct {
	ID        string `json:"id, omitempty"`
	FirstName string `json:"firstname, omitempty"`
	LastName  string `json:"lastname, omitempty"`
	Address   string `json:"address, omitempty"`
}

type City struct {
	CityCode string `json:"citycode, omitempty"`
	CityName string `json:"cityname, omitempty"`
	Province string `json:"province, omitempty"`
}

var users []User

var cities []City

func main() {
	fmt.Println("Starting application......")

	initializeData()

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"firstname": &graphql.Field{
				Type: graphql.String,
			},
			"lastname": &graphql.Field{
				Type: graphql.String,
			},
			"address": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	cityType := graphql.NewObject(graphql.ObjectConfig{
		Name: "City",
		Fields: graphql.Fields{
			"citycode": &graphql.Field{
				Type: graphql.String,
			},
			"cityname": &graphql.Field{
				Type: graphql.String,
			},
			"province": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return users, nil
				},
			},
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var user User
					var result User

					user.ID = p.Args["id"].(string)

					for _, usr := range users {
						if usr.ID == user.ID {
							result.ID = usr.ID
							result.FirstName = usr.FirstName
							result.LastName = usr.LastName
						}
					}
					return result, nil
				},
			},
			"cities": &graphql.Field{
				Type: graphql.NewList(cityType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return cities, nil
				},
			},
			"city": &graphql.Field{
				Type: cityType,
				Args: graphql.FieldConfigArgument{
					"citycode": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var city City
					var result City

					city.CityCode = p.Args["citycode"].(string)

					for _, c := range cities {
						if c.CityCode == city.CityCode {
							result.CityCode = c.CityCode
							result.CityName = c.CityName
							result.Province = c.Province
						}
					}
					return result, nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: graphql.Fields{},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result :=
			graphql.Do(graphql.Params{
				Schema:        schema,
				RequestString: r.URL.Query().Get("query"),
			})
		json.NewEncoder(w).Encode(result)
	})

	http.ListenAndServe(":8000", nil)
}

func initializeData() {
	users = append(users, User{ID: "1", FirstName: "John", LastName: "Doe", Address: "New York"})
	users = append(users, User{ID: "2", FirstName: "Jim", LastName: "Carrey", Address: "Washington DC"})
	users = append(users, User{ID: "3", FirstName: "Rhoma", LastName: "Cahyanti", Address: "Surabaya"})
	users = append(users, User{ID: "4", FirstName: "Imania", LastName: "Ramadhani", Address: "Surabaya"})

	cities = append(cities, City{CityCode: "SBY", CityName: "Surabaya", Province: "Jawa Timur"})
	cities = append(cities, City{CityCode: "MLG", CityName: "Malang", Province: "Jawa Timur"})
	cities = append(cities, City{CityCode: "SDA", CityName: "Sidoarjo", Province: "Jawa Timur"})
	cities = append(cities, City{CityCode: "GRS", CityName: "Gresik", Province: "Jawa Timur"})
	cities = append(cities, City{CityCode: "JKT", CityName: "Jakarta", Province: "DKI Jakarta"})
	cities = append(cities, City{CityCode: "BGD", CityName: "Bandung", Province: "Jawa Barat"})
}

package main

import (
	"fmt"
	"log"
	"os"

	controllers "suno/demo-rest/controller"
	"suno/demo-rest/model"

	// "log"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/labstack/echo"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var filename string = "./members.db"

	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles("config.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening Port:", config.Int("server.port", 8008))
	fmt.Println()

	_, err = os.Stat(filename)
	if err == nil {
		if err = os.Remove(filename); err != nil {
			log.Fatal(err)
			return
		}
	}

	_, err = controllers.InitDB(filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Test Code (DB)
	if config.Bool("debug.test", true) {
		printMembers()
		printMember("sunoj")
		if config.Bool("debug.update", false) {
			testUpdateMember("sunoj")
		}
		if config.Bool("debug.delete", false) {
			testDeleteMember("sunoj")
		}
	}

	// Web Service
	server := echo.New()
	server.GET("/api/v1/members", controllers.RestApiRetrieveMembers)
	server.GET("/api/v1/members/:id", controllers.RestApiRetrieveMember)
	server.POST("/api/v1/members/:id", controllers.RestApiNewMember)
	server.PUT("/api/v1/members/:id", controllers.RestApiModifyMember)
	server.DELETE("/api/v1/members/:id", controllers.RestApiRemoveMember)

	port := fmt.Sprintf(":%d", config.Int("server.port", 8008))
	server.Logger.Fatal(server.Start(port))
}

func printMembers() {
	var members []model.ShortMember

	members, err := controllers.GetMembers()
	if err != nil {
		log.Fatal(err)
	}

	if members != nil {
		fmt.Println("Members:")
		for i := 0; i < len(members); i++ {
			fmt.Printf("%s, %s, %s\n", members[i].Id, members[i].Name, members[i].Email)
		}

		fmt.Println()
	}
}

func printMember(id string) {
	var member model.Member

	member, err := controllers.GetMember(id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Member Information: %s, %s, %s\n", member.Id, member.Name, member.Email)
	if member.Favorites != nil {
		fmt.Print("Favorites: ")
		for i := 0; i < len(member.Favorites); i++ {
			if i == 0 {
				fmt.Printf("%s", member.Favorites[i])
			} else {
				fmt.Printf(", %s", member.Favorites[i])
			}
		}

		fmt.Println()
	}
}

func testUpdateMember(id string) {
	var favorites model.Favorites
	favorites = append(favorites, "User Favorite#1")
	favorites = append(favorites, "User Favorite#2")

	if err := controllers.UpdateMember(id, "User Name", "Email Address", favorites); err != nil {
		log.Fatal(err)
	}

	printMember(id)
}

func testDeleteMember(id string) {
	if err := controllers.DeleteMember(id); err != nil {
		log.Fatal(err)
	}
}

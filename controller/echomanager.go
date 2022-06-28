package controllers

import (
	"log"
	"net/http"
	"suno/demo-rest/model"

	"github.com/labstack/echo"
)

func RestApiRetrieveMembers(c echo.Context) error {
	var members []model.ShortMember

	members, err := GetMembers()
	if err != nil {
		if err.Error() == "404" {
			log.Println("Retrieve Member List : not found")
			return c.NoContent(http.StatusNotFound)
		} else {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	if members != nil {
		log.Println("Retrieve Member List")
		return c.JSON(http.StatusOK, members)
	} else {
		log.Println("Retrieve Member List : no members")
		return c.NoContent(http.StatusNotFound)
	}
}

func RestApiRetrieveMember(c echo.Context) error {
	id := c.Param("id")
	member, err := GetMember(id)
	if err != nil {
		if err.Error() == "404" {
			log.Printf("Retrieve Member %s does not found\n", id)
			return c.NoContent(http.StatusNotFound)
		} else {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	log.Printf("Retrieve Member %s", id)
	return c.JSON(http.StatusOK, member)
}

func RestApiNewMember(c echo.Context) error {
	var member model.Member

	id := c.Param("id")

	if err := c.Bind(&member); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if err := AddMember(id, member.Name, member.Email); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if err := AddFavorites(id, member.Favorites); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	log.Printf("Create User %s", id)
	return c.JSON(http.StatusCreated, model.ResultOK)
}

func RestApiModifyMember(c echo.Context) error {
	var member model.Member

	id := c.Param("id")

	if err := c.Bind(&member); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if err := UpdateMember(id, member.Name, member.Email, member.Favorites); err != nil {
		if err.Error() == "204" {
			log.Println(err)
			return c.NoContent(http.StatusNoContent)
		} else {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	log.Printf("Update User %s", id)
	return c.JSON(http.StatusOK, model.ResultOK)
}

func RestApiRemoveMember(c echo.Context) error {
	id := c.Param("id")

	if err := DeleteMember(id); err != nil {
		if err.Error() == "204" {
			log.Println(err)
			return c.NoContent(http.StatusNoContent)
		} else {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	log.Printf("Remove User %s", id)
	return c.JSON(http.StatusOK, model.ResultOK)
}

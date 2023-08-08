package http

import (
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h Handler) PostShow(c echo.Context) error {
	show := entities.Show{}

	err := c.Bind(&show)
	if err != nil {
		return err
	}

	show.ShowID = uuid.New()
	err = h.showsRepository.AddShow(c.Request().Context(), show)
	if err != nil {
		return fmt.Errorf("error storing show into db: %w", err)
	}

	return c.JSON(http.StatusCreated, show)
}

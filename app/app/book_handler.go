package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"

	"myapp/model"
	"myapp/repository"
	"myapp/util/validator"
)

func (app *App) HandleListBooks(c *gin.Context) {
	books, err := repository.ListBooks(app.db)
	if err != nil {
		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}

	if books == nil {
		fmt.Fprint(c.Writer, "[]")
		return
	}

	dtos := books.ToDto()
	if err := json.NewEncoder(c.Writer).Encode(dtos); err != nil {
		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (app *App) HandleCreateBook(c *gin.Context) {
	form := &model.BookForm{}
	if err := json.NewDecoder(c.Request.Body).Decode(form); err != nil {
		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := app.validator.Struct(form); err != nil {
		app.logger.Warn().Err(err).Msg("")

		resp := validator.ToErrResponse(err)
		if resp == nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrFormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			app.logger.Warn().Err(err).Msg("")

			c.Writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrJsonCreationFailure)
			return
		}

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		c.Writer.Write(respBody)
		return
	}

	bookModel, err := form.ToModel()
	if err != nil {
		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	book, err := repository.CreateBook(app.db, bookModel)
	if err != nil {
		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrDataCreationFailure)
		return
	}

	app.logger.Info().Msgf("New book created: %d", book.ID)
	c.Writer.WriteHeader(http.StatusCreated)
}

func (app *App) HandleReadBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil || id == 0 {
		app.logger.Info().Msgf("can not parse ID: %v", id)

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	book, err := repository.ReadBook(app.db, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}

		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}

	dto := book.ToDto()
	if err := json.NewEncoder(c.Writer).Encode(dto); err != nil {
		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (app *App) HandleUpdateBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil || id == 0 {
		app.logger.Info().Msgf("can not parse ID: %v", id)

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	form := &model.BookForm{}
	if err := json.NewDecoder(c.Request.Body).Decode(form); err != nil {
		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := app.validator.Struct(form); err != nil {
		app.logger.Warn().Err(err).Msg("")

		resp := validator.ToErrResponse(err)
		if resp == nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrFormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			app.logger.Warn().Err(err).Msg("")

			c.Writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrJsonCreationFailure)
			return
		}

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		c.Writer.Write(respBody)
		return
	}

	bookModel, err := form.ToModel()
	if err != nil {
		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	bookModel.ID = uint(id)
	if err := repository.UpdateBook(app.db, bookModel); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}

		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrDataUpdateFailure)
		return
	}

	app.logger.Info().Msgf("Book updated: %d", id)
	c.Writer.WriteHeader(http.StatusAccepted)
}

func (app *App) HandleDeleteBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil || id == 0 {
		app.logger.Info().Msgf("can not parse ID: %v", id)

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := repository.DeleteBook(app.db, uint(id)); err != nil {
		app.logger.Warn().Err(err).Msg("")

		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}

	app.logger.Info().Msgf("Book deleted: %d", id)
	c.Writer.WriteHeader(http.StatusAccepted)
}

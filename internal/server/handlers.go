package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"../database"

	"github.com/gin-gonic/gin"
)

//GetCurrencies Method
func (s *Server) GetCurrencies(c *gin.Context) {
	currencies, err := s.db.GetCurrencies(c.Request.Context())
	if err != nil {
		c.Error(err)
	}

	c.JSON(200, gin.H{
		"result": currencies,
	})
}

//CreateCurrency Method
func (s *Server) CreateCurrency(c *gin.Context) {
	var input database.Currency
	c.ShouldBindJSON(&input)
	if err := s.db.CreateCurrency(c.Request.Context(), input); err != nil {
		c.Error(err)
	}
	c.Status(200)
	if err := UpdateListCurrencies(s.db.(*database.DB)); err != nil {
		fmt.Println(err)
	}
}

//DeleteCurrency Method
func (s *Server) DeleteCurrency(c *gin.Context) {
	classIDint, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.New("text string"))
	}
	if err := s.db.DeleteCurrency(c.Request.Context(), classIDint); err != nil {
		c.Error(err)
	}
	c.Status(200)
}

//UpdateCurrency Method
func (s *Server) UpdateCurrency(c *gin.Context) {
	classIDint, err := strconv.Atoi(c.Param("id"))
	var input database.Currency
	c.ShouldBindJSON(&input)
	if err != nil {
		c.Error(errors.New("text string"))
	}
	if err := s.db.UpdateCurrency(c.Request.Context(), input, classIDint); err != nil {
		c.Error(err)
	}
	c.Status(200)
}

// ConvertCurrency Structure
type ConvertCurrency struct {
	CurrencyFrom string
	CurrencyTo   string
	Value        float64
}

//ConvertCurrency Method
func (s *Server) ConvertCurrency(c *gin.Context) {
	var input ConvertCurrency
	c.ShouldBindJSON(&input)
	currencies, err := s.db.GetCurrencies(c.Request.Context())
	if err != nil {
		c.Error(err)
	}
	for _, valcur := range currencies {
		if input.CurrencyFrom == valcur.Currency1 {
			output := input.Value / valcur.Rate
			c.JSON(200, gin.H{
				"result": output,
			})
			break
		}
	}

}

// BankCurrency Structure
type BankCurrency struct {
	Date  string
	Rates map[string]float64
}

// ListCurrencies func
func ListCurrencies(db *database.DB, sleepTime time.Duration) {

	for {
		UpdateListCurrencies(db)
		time.Sleep(sleepTime)
	}
}

// UpdateListCurrencies func
func UpdateListCurrencies(db *database.DB) error {
	var ListOfCurrency BankCurrency
	url := "https://www.cbr-xml-daily.ru/latest.js"

	res, err := http.Get(url)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ListOfCurrency)

	if err != nil {
		return err
	}

	dbres, err := db.GetDBCurrencies()
	if err != nil {
		return err
	}
	for _, val := range dbres {
		if valCurr, ok := ListOfCurrency.Rates[val.Currency1]; ok {
			err = db.UpdateDBCurrencies(val, valCurr)
			if err != nil {
				return err
			}
		}
	}
	return err
}

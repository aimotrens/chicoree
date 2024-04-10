package controllers

import (
	"chicoree/ent"

	"github.com/gin-gonic/gin"
)

type EntClientConstructor func() (*ent.Client, error)

type ControllerBase interface {
	BeginTx(c *gin.Context) *ent.Tx
	WithTx(c *gin.Context, f func(tx *ent.Tx) error) error
	WithoutTx(c *gin.Context, f func(tx *ent.Client) error) error
}

type BaseController struct {
	dbClientConstructor EntClientConstructor
}

func NewBaseController(ecc EntClientConstructor) *BaseController {
	return &BaseController{dbClientConstructor: ecc}
}

// Erstellt eine Transaktion
// Bei einem Fehler wird ein HTTP 500 an den Client und nil an den Aufrufer zurückgegeben
func (b *BaseController) BeginTx(c *gin.Context) *ent.Tx {
	client, err := b.dbClientConstructor()
	if err != nil {
		c.JSON(500, gin.H{"message": "Error connecting database: " + err.Error()})
		return nil
	}

	tx, err := client.Tx(c)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error starting transaction: " + err.Error()})
		return nil
	}

	return tx
}

// Führt eine Funktion mit DB-Verbindung aus
func (b *BaseController) WithoutTx(c *gin.Context, f func(tx *ent.Client) error) error {
	dbClient, err := b.dbClientConstructor()
	if dbClient == nil {
		c.AbortWithError(500, err)
		return err
	}
	defer dbClient.Close()

	err = f(dbClient)
	if err != nil {
		c.AbortWithError(500, err)
		return err
	}

	return nil
}

// Führt eine Funktion als Transaktion aus
func (b *BaseController) WithTx(c *gin.Context, f func(tx *ent.Tx) error) error {
	dbClient, err := b.dbClientConstructor()
	if dbClient == nil {
		c.AbortWithError(500, err)
		return err
	}
	defer dbClient.Close()

	tx, err := dbClient.Tx(c)
	if err != nil {
		c.AbortWithError(500, err)
		return err
	}

	err = f(tx)
	if err != nil {
		tx.Rollback()
		c.AbortWithError(500, err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		c.AbortWithError(500, err)
		return err
	}

	return nil
}

package main

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

type MerkleTree struct {
	Root  string `gorm:"primary_key"`
	Nodes string `gorm:"type:text"`
}

// hashes 2 base-10 strings and returns the result as a string
func poseidon2(a string, b string) string {
	aBig := new(big.Int)
	aBig, ok1 := aBig.SetString(a, 10)
	bBig := new(big.Int)
	bBig, ok2 := bBig.SetString(a, 10)
	if !ok1 || !ok2 {
		log.Fatal("Error converting to big int")
	}
	res, err := poseidon.Hash([]*big.Int{aBig, bBig})
	if err != nil {
		log.Fatal(err)
	}

	return res.String()
}

var db *gorm.DB
var err error

// creates a tree and writes it to the database
// if colliding with another tree with same root, assume that one must be a prefix of the other and keep the larger one
func create_tree(tree MerkleTree) error {
	var existing_tree MerkleTree
	if db.First(&existing_tree, "root = ?", tree.Root).Error == nil {
		if len(existing_tree.Nodes) >= len(tree.Nodes) {
			fmt.Printf("Cannot update tree %s (%d nodes) with smaller tree (%d nodes)\n", tree.Root, len(existing_tree.Nodes), len(tree.Nodes))
			return errors.New("tree with same root already exists")
		} else {
			fmt.Printf("Updating tree %s (%d nodes) with larger tree (%d nodes)\n", tree.Root, len(existing_tree.Nodes), len(tree.Nodes))
			db.Model(&existing_tree).Updates(tree)
		}
	} else {
		fmt.Printf("Creating tree %s (%d nodes)\n", tree.Root, len(tree.Nodes))
		db.Create(&tree)
		fmt.Printf("Created!\n")
	}
	return nil
}

func main() {
	dsn := "host=localhost user=mroots dbname=mroots port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&MerkleTree{})

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))

	r.GET("/tree/:root", func(c *gin.Context) {
		root := c.Param("root")
		var existing_tree MerkleTree
		if db.First(&existing_tree, "root = ?", root).Error == nil {
			fmt.Printf("Found existing tree\n")
			c.JSON(200, gin.H{
				"nodes": strings.Split(existing_tree.Nodes, ","),
			})
		} else {
			fmt.Printf("Did not find existing tree\n")
			c.JSON(200, gin.H{
				"nodes": []string{"", root},
			})
		}
	})

	r.POST("/tree", func(c *gin.Context) {
		data := struct {
			Id     string   `json:"id"`
			Leaves []string `json:"leaves"`
		}{}
		if err = c.BindJSON(&data); err != nil {
			log.Fatal(err)
		}
		n := len(data.Leaves)
		nodes := make([]string, 2*n)
		for i := 2*n - 1; i > 0; i-- {
			if i >= n {
				nodes[i] = data.Leaves[i-n]
			} else {
				nodes[i] = poseidon2(nodes[2*i], nodes[2*i+1])
			}
		}
		err = create_tree(MerkleTree{Root: nodes[1], Nodes: strings.Join(nodes, ",")})
		var statusCode = 200
		if err != nil {
			statusCode = 409 // conflict
		}
		c.JSON(statusCode, nodes[1])
	})

	r.Run()
}

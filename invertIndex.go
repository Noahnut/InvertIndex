package invertIndex

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type addDocument struct {
	document chan string
	FileName string `json:"key"`
}

type InvertIndex struct {
	searchParameter  []string
	invertIndexTable map[string][]string
	adddocument      addDocument
}

func (c *InvertIndex) readDataFromFile(fileName string) []string {

	curretntPath, _ := os.Getwd()
	file, err := os.Open(curretntPath + "/" + fileName)
	invertIndexKey := make([]string, 0)

	if err != nil {
		log.Println(err)
		return nil
	}

	defer file.Close()

	r := bufio.NewReader(file)

	for _, e := range c.searchParameter {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Println(err)
		}

		keySet := strings.Split(string(line), ":")

		if len(keySet) == 0 {
			log.Println("read the string from file fail")
			return nil
		}

		if e != keySet[0] {
			log.Printf("key %s is not the expect key\n", e)
			continue
		}

		invertIndexKey = append(invertIndexKey, keySet[1])
	}

	return invertIndexKey
}

func (c *InvertIndex) generateDocumentID(keyArray []string) string {
	hasher := sha256.New()
	currentTime := fmt.Sprint(time.Now().Unix())
	hasher.Write([]byte(currentTime))

	for _, e := range keyArray {
		hasher.Write([]byte(e))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func (c *InvertIndex) addDocumentToIndexIndex(fileName string) {
	keyArray := c.readDataFromFile(fileName)

	if keyArray == nil {
		return
	}

	documentID := c.generateDocumentID(keyArray)

	for _, e := range keyArray {
		var prefix string
		if len(e) < 4 {
			prefix = e
		} else {
			prefix = e[:4]
		}

		_, ok := c.invertIndexTable[prefix]
		if !ok {
			c.invertIndexTable[prefix] = make([]string, 0)
		}
		c.invertIndexTable[prefix] = append(c.invertIndexTable[prefix], documentID)
	}

	e := os.Rename(fileName, documentID)

	if e != nil {
		log.Fatal(e)
	}
}

func (c *InvertIndex) addNewDocumentRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Exit addNewDocumentRoutine")
		case fileName := <-c.adddocument.document:
			c.addDocumentToIndexIndex(fileName)
		}
	}
}

func (c *InvertIndex) addNewDocument(context *gin.Context) {
	context.BindJSON(&c.adddocument)
	log.Println(&c.adddocument.FileName)

	if len(c.adddocument.FileName) == 0 {
		context.String(http.StatusBadRequest, "File Name is empty")
	}

	c.adddocument.document <- c.adddocument.FileName
	context.String(200, "Success")
}

func (c *InvertIndex) getDocumentFromInvertIndex(searchWord string) []string {
	var prefixWord string
	if len(searchWord) > 4 {
		prefixWord = searchWord[:4]
	} else {
		prefixWord = searchWord
	}

	document, ok := c.invertIndexTable[prefixWord]
	if !ok {
		return nil
	}

	return document
}

func (c *InvertIndex) getDocument(context *gin.Context) {
	searchWord := context.Query("search")
	documentNames := c.getDocumentFromInvertIndex(searchWord)

	if documentNames == nil {
		context.String(http.StatusOK, "document Not Found")
	}
	context.JSON(http.StatusOK, gin.H{
		"documentName": documentNames,
	})
}

func (c *InvertIndex) initInvertIndex(parameter []string) {
	c.searchParameter = parameter
	c.invertIndexTable = make(map[string][]string)
	c.adddocument = addDocument{}
	c.adddocument.document = make(chan string, 20)
}

func (c *InvertIndex) NewInvertIndex(parameter ...string) {
	ctx, cancel := context.WithCancel(context.Background())
	c.initInvertIndex(parameter)
	route := gin.Default()
	go c.addNewDocumentRoutine(ctx)
	route.POST("/AddNewDocument", c.addNewDocument)
	route.GET("/GetDocument", c.getDocument)
	route.Run(":7777")
	cancel()
}

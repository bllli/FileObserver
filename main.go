package main

import (
	"fmt"
	"github.com/gofiber/fiber"
	"os"
	"sort"
	"strings"
	"time"
)

type EnumOrderBy int8
type EnumOrdering int8
type EnumFilter int8

const (
	OrderByNone     EnumOrderBy = 0
	OrderBySize     EnumOrderBy = 1
	OrderByUpdateTs EnumOrderBy = 2

	OrderingAsc  EnumOrdering = 1
	OrderingDesc EnumOrdering = 2

	FilterNil     EnumFilter = 0
	FilterDirOnly EnumFilter = 1
	FilterNoDir   EnumFilter = 2
)

type Query struct {
	Path     string
	Ordering EnumOrdering
	OrderBy  EnumOrderBy
	Filter   EnumFilter
	Keyword  string
}

func ReadDir(query Query) ([]FileInfo, error) {
	f, err := os.Open(query.Path)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(1000) // todo config this
	if err != nil {
		return nil, err
	}
	var infos = make(FileInfoSlice, 0)

	for i, file := range files {
		if query.Keyword != "" {
			if !strings.Contains(file.Name(), query.Keyword) {
				continue
			}
			if i > 100 {
				// todo config this
				break
			}
		}

		if FilterDirOnly == query.Filter {
			if !file.IsDir() {
				continue
			}
		} else if FilterNoDir == query.Filter {
			if file.IsDir() {
				continue
			}
		}

		//fmt.Printf("Name: %s, IsDir: %t, Size: %d \n", file.Name(), file.IsDir(), file.Size())
		infos = append(infos, FileInfo{file.Name(), file.Size(), file.IsDir(), file.ModTime()})
	}

	if query.OrderBy != OrderByNone {
		sort.Slice(infos, func(i, j int) (b bool) {
			if query.OrderBy == OrderBySize {
				b = infos.LessBySize(i, j)
			} else if query.OrderBy == OrderByUpdateTs {
				b = infos.LessByModTime(i, j)
			}
			if query.Ordering == OrderingDesc {
				b = !b
			}
			return
		})
	}

	return infos, nil
}

func viewGetPath(c *fiber.Ctx) {
	query := Query{}
	err := c.QueryParser(&query)
	if err != nil {
		c.Status(500).Send(fmt.Sprintf("query parse fail: %v", err))
		return
	}

	fmt.Printf("query: %v\n", query)

	if query.Path == "" {
		query.Path = GetBasePath()
	}
	if !strings.HasPrefix(query.Path, defaultBasePath) {
		c.Status(500).Send(fmt.Sprintf("path must start with %s", defaultBasePath))
		return
	}
	if strings.Contains(query.Path, "/../") {
		c.Status(500).Send("别闹了")
		return
	}
	dir, err := ReadDir(query)
	if err != nil {
		c.Status(500).JSON(&err)
		return
	}
	c.JSON(&dir)
}

func main() {
	fmt.Printf("hello, world!\n")

	app := fiber.New()
	app.Settings.ServerHeader = "File Observer @bllli"
	app.Settings.ReadTimeout = 3 * time.Second
	app.Settings.WriteTimeout = 3 * time.Second
	app.Settings.IdleTimeout = 3 * time.Second

	app.Settings.CaseSensitive = true

	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("Hello, World!")
	})
	app.Get("/path/", viewGetPath)

	app.Listen(8003)
}

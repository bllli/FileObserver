package main

import (
	"fmt"
	"github.com/gofiber/fiber"
	"os"
	"sort"
	"strconv"
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

func ReadDir(
	path string, keyword string,
	ordering EnumOrdering, orderBy EnumOrderBy, filter EnumFilter,
) ([]FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(1000) // todo config this
	if err != nil {
		return nil, err
	}
	var infos = make(FileInfoSlice, 0)

	for i, file := range files {
		if keyword != "" {
			if !strings.Contains(file.Name(), keyword) {
				continue
			}
			if i > 100 {
				// todo config this
				break
			}
		}

		if FilterDirOnly == filter {
			if !file.IsDir() {
				continue
			}
		} else if FilterNoDir == filter {
			if file.IsDir() {
				continue
			}
		}

		//fmt.Printf("Name: %s, IsDir: %t, Size: %d \n", file.Name(), file.IsDir(), file.Size())
		infos = append(infos, FileInfo{file.Name(), file.Size(), file.IsDir(), file.ModTime()})
	}

	if orderBy != OrderByNone {
		sort.Slice(infos, func(i, j int) (b bool) {
			if orderBy == OrderBySize {
				b = infos.LessBySize(i, j)
			} else if orderBy == OrderByUpdateTs {
				b = infos.LessByModTime(i, j)
			}
			if ordering == OrderingDesc {
				b = !b
			}
			return
		})
	}

	return infos, nil
}

func viewGetPath(c *fiber.Ctx) {
	//type Query struct {
	//	Path string
	//	Ordering EnumOrdering
	//	OrderBy EnumOrderBy
	//	keyword string
	//}

	path := c.Query("path")
	keyword := c.Query("keyword")
	ordering, err := strconv.Atoi(c.Query("ordering"))
	orderingEnum := EnumOrdering(ordering)
	orderBy, err := strconv.Atoi(c.Query("orderBy"))
	orderByEnum := EnumOrderBy(orderBy)
	filter, err := strconv.Atoi(c.Query("filter"))
	filterEnum := EnumFilter(filter)

	fmt.Printf(
		"path: '%s', keyword: '%s', ordering %d, orderBy %d filter: %d\n",
		path, keyword, ordering, orderBy, filterEnum)

	if path == "" {
		path = GetBasePath()
	}
	if !strings.HasPrefix(path, defaultBasePath) {
		c.Status(500).Send(fmt.Sprintf("path must start with %s", defaultBasePath))
		return
	}
	if strings.Contains(path, "/../") {
		c.Status(500).Send("别闹了")
		return
	}
	dir, err := ReadDir(path, keyword, orderingEnum, orderByEnum, filterEnum)
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

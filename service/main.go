package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/gofiber/fiber/v2"
)

type CoomonStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CommonResponse[T any] struct {
	Status CoomonStatus `json:"status"`
	Data   T            `json:"data"`
}

type File struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Size     int64   `json:"size"`
	IsDir    bool    `json:"isDir"`
	Children []*File `json:"children,omitempty"`
	Parent   *File   `json:"-"`
}

const (
	baseDir = "../example-dir"
)

func main() {
	app := fiber.New()

	app.Get("/files", func(c *fiber.Ctx) error {
		files := make(map[string]*File)
		err := filepath.Walk(baseDir,
			func(path string, file os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				files[path] = &File{
					Name:     file.Name(),
					Path:     path,
					Size:     file.Size(),
					IsDir:    file.IsDir(),
					Children: make([]*File, 0),
				}
				return nil
			})

		if err != nil {
			log.Fatal(err)
		}

		var result *File
		for path, file := range files {
			parrnetPath := filepath.Dir(path)
			parrnetFile, isExists := files[parrnetPath]

			if !isExists {
				result = file
				continue
			} else {
				file.Parent = parrnetFile
				newChildren := append(parrnetFile.Children, file)
				sort.Slice(newChildren, func(i, j int) bool {
					return newChildren[i].IsDir
				})
				parrnetFile.Children = newChildren
			}
		}

		response := CommonResponse[*File]{
			Status: CoomonStatus{
				Code:    "200",
				Message: "Success",
			},
			Data: result,
		}

		err = c.JSON(response)
		return err
	})

	app.Listen(":3000")
}

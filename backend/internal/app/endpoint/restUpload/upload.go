package restUpload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

type Upload interface {
}

type Endpoint struct {
	//Upload Upload
}

func (e Endpoint) UploadFileHandler(c *gin.Context) {
	fmt.Println("1")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "Ошибка при получении файла: %s", err.Error())
		return
	}
	defer file.Close()

	fmt.Println("1")

	// Создание файла на сервере
	dst, err := os.Create("media/" + header.Filename)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Ошибка при создании файла: %s", err.Error())
		return
	}
	defer dst.Close()

	fmt.Println("1")

	// Копирование содержимого загруженного файла в созданный файл
	if _, err := io.Copy(dst, file); err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Ошибка при копировании файла: %s", err.Error())
		return
	}

	fmt.Println("1")
	// Вызов обработчика после загрузки файла
	//processFile(header.Filename)

	c.String(http.StatusOK, "Файл успешно загружен: %s", header.Filename)
}

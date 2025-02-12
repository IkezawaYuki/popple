package controller

import (
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/IkezawaYuki/popple/internal/usecase"
	"github.com/labstack/echo/v4"
)

type BatchController struct {
	batchUsecase  usecase.BatchUsecase
	httpPresenter presenter.Presenter
}

func NewBatchController(batchUsecase *usecase.BatchUsecase, presenter2 *presenter.Presenter) BatchController {
	return BatchController{
		batchUsecase:  *batchUsecase,
		httpPresenter: *presenter2,
	}
}

// Execute godoc
//
//	@Summary		バッチ実行
//	@Description	全顧客に対して、InstagramからWordpressへのデータ連携を行います。
//	@Tags			Badge
//	@Security		Token
//	@Router			/badge/execute [get]
func (ctr *BatchController) Execute(c echo.Context) error {
	resp, err := ctr.batchUsecase.Execute(c.Request().Context())
	return c.JSON(ctr.httpPresenter.Generate(err, resp))
}

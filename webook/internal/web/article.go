package web

import (
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/service"
	"geek-basic-go/webook/internal/web/jwt"
	"geek-basic-go/webook/pkg/ginx"
	"geek-basic-go/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	svc     service.ArticleService
	intrSvc service.InteractiveService
	l       logger.LoggerV1
	biz     string
}

func NewArticleHandler(svc service.ArticleService,
	intrSvc service.InteractiveService,
	l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		intrSvc: intrSvc,
		l:       l,
		biz:     "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)

	// 创作者接口
	// List接口，一般是GET的，形如list?offset=?&limit=?, 这里定义成post，然后通过body接收参数
	g.POST("/list", h.List)
	g.GET("/detail/:id", h.Detail)
	pub := g.Group("/pub")
	pub.GET("/:id", h.PubDetail)
	pub.POST("/like", h.Like)
	pub.POST("/collect", h.Collect)
}

// Edit 返回article id
func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	id, err := h.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("保存文章失败", logger.Int64("uid", uc.Uid), logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: id,
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	id, err := h.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("发表文章失败", logger.Int64("uid", uc.Uid), logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: id,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	err := h.svc.Withdraw(ctx, uc.Uid, req.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("撤回文章失败",
			logger.Int64("uid", uc.Uid),
			logger.Int64("aid", req.Id),
			logger.Error(err),
		)
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: req.Id,
	})
}

func (h *ArticleHandler) List(ctx *gin.Context) {
	var page Page
	if err := ctx.Bind(&page); err != nil {
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	arts, err := h.svc.GetByAuthor(ctx, uc.Uid, page.Offset, page.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("查找文章列表失败",
			logger.Error(err),
			logger.Int("offset", page.Offset),
			logger.Int("limit", page.Limit),
			logger.Int64("uid", uc.Uid))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: slice.Map[domain.Article, ArticleVo](arts, func(idx int, src domain.Article) ArticleVo {
			return toVo(src)
		}),
	})
}

func toVo(art domain.Article) ArticleVo {
	return ArticleVo{
		Id:       art.Id,
		Title:    art.Title,
		Abstract: art.Abstract(),
		//Content:    art.Content,
		AuthorId:   art.Author.Id,
		AuthorName: art.Author.Name,
		Status:     art.Status.ToUint8(),
		Ctime:      art.Ctime.Format(time.DateTime),
		Utime:      art.Utime.Format(time.DateTime),
	}
}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "id参数错误",
		})
		h.l.Warn("查找文章失败, id格式不对",
			logger.Error(err),
			logger.String("id", idStr))
		return
	}
	art, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("查找文章失败",
			logger.Error(err),
			logger.Int64("id", id))
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	if art.Author.Id != uc.Uid {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("非法查询文章",
			logger.Error(err),
			logger.Int64("id", id),
			logger.Int64("uid", uc.Uid))
		return
	}
	artVo := ArticleVo{
		Id:    art.Id,
		Title: art.Title,
		//Abstract: art.Abstract(),
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
		Ctime:    art.Ctime.Format(time.DateTime),
		Utime:    art.Utime.Format(time.DateTime),
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: artVo,
	})
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "id参数错误",
		})
		h.l.Warn("查找文章失败, id格式不对",
			logger.Error(err),
			logger.String("id", idStr))
		return
	}

	var (
		eg   errgroup.Group
		intr domain.Interactive
		art  domain.Article
	)
	uc := ctx.MustGet("user").(jwt.UserClaims)
	eg.Go(func() error {
		var er error
		art, er = h.svc.GetPubById(ctx, id, uc.Uid)
		return er
	})

	eg.Go(func() error {
		var er error
		intr, er = h.intrSvc.Get(ctx, h.biz, id, uc.Uid)
		return er
	})

	// 等待结果
	err = eg.Wait()
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Warn("查找文章失败, 系统错误",
			logger.Error(err),
			logger.Int64("id", id))
		return
	}
	// 在service通过kafka传递消息，这里不需要了
	/*go func() {
		// 1. 如果需要摆脱原本主链路的超时控制，创建一个新的
		// 2. 也可以直接只用ctx，由主链路来控制超时
		newCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := h.intrSvc.IncrReadCnt(newCtx, h.biz, art.Id)
		if er != nil {
			h.l.Error("更新阅读数失败",
				logger.Int64("aid", art.Id),
				logger.Error(er))
		}

	}()*/
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: ArticleVo{
			Id:         art.Id,
			Title:      art.Title,
			Content:    art.Content,
			AuthorId:   art.Author.Id,
			AuthorName: art.Author.Name,

			ReadCnt:    intr.ReadCnt,
			LikeCnt:    intr.LikeCnt,
			CollectCnt: intr.CollectCnt,
			Liked:      intr.Liked,
			Collected:  intr.Collected,

			Status: art.Status.ToUint8(),
			Ctime:  art.Ctime.Format(time.DateTime),
			Utime:  art.Utime.Format(time.DateTime),
		},
	})
}

func (h *ArticleHandler) Like(ctx *gin.Context) {
	type Req struct {
		Id   int64 `json:"id"`
		Like bool  `json:"like"` //ture 点赞，false 不点赞
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	var err error
	if req.Like {
		// 点赞
		err = h.intrSvc.Like(ctx, h.biz, req.Id, uc.Uid)
	} else {
		//取消点赞
		err = h.intrSvc.CancelLike(ctx, h.biz, req.Id, uc.Uid)
	}
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统异常",
		})
		h.l.Error("点赞/取消点赞失败",
			logger.Error(err),
			logger.Int64("uid", uc.Uid),
			logger.Int64("aid", req.Id))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "OK",
	})
}

func (h *ArticleHandler) Collect(ctx *gin.Context) {
	type Req struct {
		Id  int64 `json:"id"`
		Cid int64 `json:"cid"` //收藏夹id
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	err := h.intrSvc.Collect(ctx, h.biz, req.Id, req.Cid, uc.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统异常",
		})
		h.l.Error("收藏失败",
			logger.Error(err),
			logger.Int64("uid", uc.Uid),
			logger.Int64("aid", req.Id))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "OK",
	})
}

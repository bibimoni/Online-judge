package common

import (
	"context"
	"net/http"

	helper "github.com/bibimoni/Online-judge/submission-judge/src/controller"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

func InvokeUseCase[Input any, Output any](
	GetInput func(*gin.Context) (*Input, error),
	Invoke func(context.Context, *Input) (*Output, error),
	WriteOutput func(*gin.Context, *Output, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		input, err := GetInput(c)
		if err != nil {
			// panic(err)
			WriteOutput(c, nil, err)
			return
		}

		output, err := Invoke(c.Request.Context(), input)

		if err != nil {
			// panic(err)
			WriteOutput(c, nil, err)
			return
		}

		WriteOutput(c, output, nil)
	}
}

func InvokeWSUseCase[Input any, Output any](
	GetInput func(*gin.Context) (*Input, error),
	Invoke func(context.Context, *Input, chan<- *Output),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		input, err := GetInput(c)
		if err != nil {
			helper.WriteFailed(c, err, http.StatusBadRequest)
			return
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			helper.WriteFailed(c, err, http.StatusInternalServerError)
			return
		}
		defer ws.Close()

		out := make(chan *Output)
		ctx, cancel := context.WithCancel(c.Request.Context())
		defer cancel()

		go Invoke(ctx, input, out)

		for {
			select {
			case msg, ok := <-out:
				if !ok {
					return
				}
				ws.WriteJSON(msg)
			case <-ctx.Done():
				return
			}
		}
	}
}

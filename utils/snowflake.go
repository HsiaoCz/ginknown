package utils

import (
	"time"

	"github.com/HsiaoCz/ginknown/etc"
	"github.com/bwmarrin/snowflake"
)

// 雪花算法获取ID

var node *snowflake.Node

func Init() (err error) {
	app := etc.Conf.AC
	var st time.Time
	st, err = time.Parse("2006-01-02", app.StartTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(app.MachineId)
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}

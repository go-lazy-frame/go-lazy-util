//go:build linux

//             ,%%%%%%%%,
//           ,%%/\%%%%/\%%
//          ,%%%\c "" J/%%%
// %.       %%%%/ o  o \%%%
// `%%.     %%%%    _  |%%%
//  `%%     `%%%%(__Y__)%%'
//  //       ;%%%%`\-/%%%'
// ((       /  `%%%%%%%'
//  \\    .'          |
//   \\  /       \  | |
//    \\/攻城狮保佑) | |
//     \         /_ | |__
//     (___________)))))))                   `\/'
/*
 * 修订记录:
 * long.qian 2022-01-17 15:15 创建
 */

/**
 * @author long.qian
 */

package file_watch

import (
	"fmt"
	go_lazy_util "github.com/go-lazy-frame/go-lazy-util"
	"github.com/rjeczalik/notify"
	"path"
)

func (receiver *FileWatch) eventHandler(watchDirs []string) {
	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 100)

	var watched []string
	for _, dir := range watchDirs {
		dir = path.Join(dir, "...")
		if !go_lazy_util.ArrayUtil.IsExistStringArray(&watched, dir) {
			// Set up a watchpoint listening for inotify-specific events within a
			// current working directory. Dispatch each InCloseWrite and InMovedTo
			// events separately to c.
			if err := notify.Watch(dir, c, notify.All, notify.InCloseWrite, notify.InMovedTo); err != nil {
				fmt.Println("目录：", dir, " 监听失败：", err)
			} else {
				fmt.Println("已监听目录：%s\n", dir)
			}
			watched = append(watched, dir)
		}
	}
	defer func() {
		notify.Stop(c)
	}()

	for {
		switch ei := <-c; ei.Event() {
		case notify.InCloseWrite, notify.InMovedTo:
			// 新文件
			if !receiver.EnableFileCreateHandler {
				continue
			}
			filePath := ei.Path()
			receiver.fileHandlerChannel <- &fileEvent{
				FilePath: filePath,
				Op:       Created,
			}
		case notify.Write:
			// 文件写入
			if !receiver.EnableFileWriteHandler {
				continue
			}
			receiver.fileHandlerChannel <- &fileEvent{
				FilePath: ei.Path(),
				Op:       Write,
			}
		case notify.Remove:
			// 删除文件
			if !receiver.EnableFileDelHandler {
				continue
			}
			receiver.fileHandlerChannel <- &fileEvent{
				FilePath: ei.Path(),
				Op:       Remove,
			}
		}
	}
}

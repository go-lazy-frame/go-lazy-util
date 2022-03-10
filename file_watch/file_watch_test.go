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
 * long.qian 2021-11-11 10:05 创建
 */

/**
 * @author long.qian
 */

package file_watch

import (
	"fmt"
	"testing"
)

func TestFileWatch(t *testing.T) {
	watch := FileWatch{WriteTime: 1000}
	watch.StartFileWatch(func(newFile string, op FileOpEvent) {
		fmt.Println(op, newFile)
	}, "/Volumes/E-NTFS/test")
	select {}
}

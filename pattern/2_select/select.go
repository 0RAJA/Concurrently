package main

/*

 */

//func T1() {
//	//常规形式
//	for { //要么无限，要么range
//		select {
//		//使用channel
//		}
//	}
//	//向chan发送迭代变量
//	for _, s := range []string{"a", "b", "c"} {
//		select {
//		case <-done:
//			return
//		case stream <- s:
//		}
//	}
//	//循环等待停止
//	for {
//		select {
//		case <-done:
//			return
//		default:
//			//执行非抢占式任务
//		}
//	}
//}

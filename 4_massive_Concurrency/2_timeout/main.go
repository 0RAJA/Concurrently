package main

/*
	超时与取消
	何时支持超时
		系统饱和：超出系统承受能力的请求返回超时，而不是一直等待。
		陈旧的数据：如果数据的处理时间过长，可能会导致数据过期，可以设置一个最长等待时间来限制/。
		试图防止死锁：为了保证系统不会死锁，建议在所有并发操作中设置超时。
	并发进程被取消的原因：
		超时：隐式取消
		用户干预：
			可以维持一个长链接，轮训报告状态给用户，或者允许用户查看状态，以及用户可以主动取消开始的操作。
		父进程取消导致子进程取消。
		复制请求：当一个请求响应后取消其他请求。
	取消的方法：
		context 和 done
	取消的影响：
		探索go程的可抢占性：
			确保运行周期比抢占周期长的功能本身都是可以抢占的。将go程的代码分为多个小段来减少确认取消和实际停止之间的时间。
		停止后的回滚操作：尽可能回滚少的操作
		重复的消息：
			如果被取消的进程在发送完结果后被取消，然后又新启动一个进程发送同样的结果。
				1. 如果并发进程是幂等的就选一个进行处理
				2. 向父进程确认：使用双向通信
*/

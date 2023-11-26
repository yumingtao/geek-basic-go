// Package service 代表领域服务，核心关键的业务逻辑，“按照道理来讲”都放到这里
// “按照道理来讲”是因为有些操作如分布式锁处理并发问题，可能下沉到repository甚至dao
package service

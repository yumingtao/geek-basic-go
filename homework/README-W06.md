# Homework Week06
## 实现
- 定义了一个ErrLogger
- 提供了一个HandleErr方法，当err!=nil时，打印error

## 问题
- 还是需要在每个err != nil的地方显示调用一下HandleErr方法，所以应该不是这个方案

## 扩展
- Go不直接支持AOP，但可以使用github.com/jakecoffman/gofer实现面向切面编程
- 业务代码中有很多地方判断err != nil, 有些是调用第三方，有些是调用自己写的方法
- 如果使用面向切面编程，代码改动量会比较大，没有尝试
# WeBook

## Project Structure
- **main.go**: 启动webook
- **internal**: 存放所有业务代码
- **internal/web/***: 主要存放web接口
- **pkg**:整个项目沉淀出来可以给别的项目使用的东西
- **domain**: 代表领域对象，业务在系统中的直接反应，可以直接理解为一个业务对象，又或是一个现实对象在中的反应
- **service**: 代表领域服务，代表一个完整的处理过程。组合各种repository，domain，也会组合别的service来共同完成一个业务
- **repository**: 代表领域对象的存储，是一个整体的抽象，只代表数据存储，数据存储有很多，如ES，关系数据库，非关系数据库，甚至文件等
- **repository/dao**: 代表数据库的操作，是操作数据库的抽象

## Notes
- 每次加了go的依赖执行 go mod tidy，确保go.mod和go.sum符合go的规范
- middleware类似于Java web中的filter，interceptor，也叫AOP解决方案
- Docker Compos相关命令
  - docker compose up: 初始化docker-compose并启动
  - docker compose down: 删除docker-compose里边创建的各种容器

## References
- [kratos](https://go-kratos.dev/en/docs)
- [go-zero](https://go-zero.dev/docs)
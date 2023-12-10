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

## K8s
### ports
![img.png](k8s-ports.png)
### ingress
![img.png](k8s-ingress.png)
- ingress代表路由规则
- service中的LoadBalancer强调的是将流量转发到pod上，ingress强调的是发到不同的service上
### Ingress 和 Ingress Controller
![img.png](k8s-ingress-vs-ingresscontroller.png)
- Ingress Controller 可以控制整个集群内部符合条件的所有Ingress
- Ingress是路由配置说明，而Ingress Controller是执行这些配置的

### ingress-nginx
- Ingress的nginx实现

### wire
#### Disadvantage
- 缺乏根据不同环境使用不同实现的能力
- 缺乏根据接口查找实现的能力
- 缺乏根据类型查找所有实现的能力

#### Advantage
- 使代码清晰，可控性强

#### Notes
- 使用wire时初始化方法最好返回接口类型，这样wire可以直接使用类型匹配，不然需要使用wire的Bind方法去绑定
- **但是Go推荐返回具体类型，和wire有冲突**

## Notes
- 识别业务变化点，超前设计但不超前实现
- 识别变化点
  - 任何第三方工具，都存在替换可能
  - 业务流程中不太合理的地方
  - 核心逻辑一定要面向接口编程

## References
- [kratos](https://go-kratos.dev/en/docs)
- [go-zero](https://go-zero.dev/docs)
- [ekit](https://github.com/ecodeclub/ekit)
- [Tencent SMS](https://cloud.tencent.com/document/product/382/43199)
- [gin](https://github.com/gin-gonic/gin)
- [gorm](https://github.com/go-gorm/gorm)
- [wire](https://github.com/google/wire)
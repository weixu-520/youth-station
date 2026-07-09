云驿 AI — 基于 RAG 的青年人才智能服务平台
一句话概述
基于 Go 微服务框架 go-zero + 字节跳动 Eino AI 框架的全栈 Web 应用，自建 RAG（检索增强生成）智能客服，集成 Milvus 向量数据库与 DeepSeek 大模型，支持 WebSocket 实时聊天。

技术栈：Go · go-zero · GORM · MySQL · Redis · RabbitMQ · WebSocket · Kafka · Eino · Milvus · DeepSeek · Vue3

周期：个人全栈独立开发

运行方式

1.在终端clone当前项目

2.在yonth_station_backend\api\etc\gateway-api.example.yaml中配置相关配置，配置完后将文件名改为gateway-api.yaml

3.cd yonth_station_backend\api在终端输入go run gateway.go运行后端

4.cd yonth_station_front在终端输入npm run dev运行前端

如何把该项目制作成一个容器：

1.# 在 Ubuntu 上

git clone https://github.com/weixu-520/youth-station.git

cd youth-station

2.配置环境变量

cp .env.example .env

nano .env   # 改 DEEPSEEK_API_KEY 和 ARK_API_KEY 为真实值

3. 一行命令启动所有服务

docker compose up -d --build

4. 导入知识库

docker exec -it yonth-backend sh -c "cd /app/cmd/seed-knowledge && go run ."




后端架构
采用 go-zero 分层架构，严格遵循 Handler → Logic → Model 职责分离，ServiceContext 集中依赖注入：

<img width="778" height="199" alt="image" src="https://github.com/user-attachments/assets/180efcf4-02c7-41c5-84c2-7a996fd03876" />


设计模式：每个请求从 HTTP Context 透传至 Logic 层，JWT 中间件将 userId/isAdmin 注入 context，Logic 通过 ctx.Value() 获取当前用户身份，无需全局变量。

Eino 框架集成（RAG 智能客服）
整体 Pipeline
用户提问 → Milvus 向量检索 → 相似度过滤 → Prompt 拼装 → DeepSeek 生成 → SSE 流式输出
              ↑                              ↑                              ↑
         Eino Retriever               Eino ChatModel              Eino StreamReader

核心组件

<img width="568" height="122" alt="image" src="https://github.com/user-attachments/assets/af1b460c-bbe2-4369-a277-5517a0ba617f" />



主要功能
用户登录界面：
<img width="1081" height="644" alt="image" src="https://github.com/user-attachments/assets/763c0618-8e8c-4151-ab23-b5ff2e4aeb58" />
主界面：
<img width="1272" height="720" alt="image" src="https://github.com/user-attachments/assets/a0866c1d-c084-47ca-a8cc-6dcfe0cb6939" />
驿站列表界面：
<img width="1279" height="727" alt="image" src="https://github.com/user-attachments/assets/d638b4d9-f08d-4fb6-a2fa-349e68702cc5" />
驿站详情界面：
<img width="1085" height="704" alt="image" src="https://github.com/user-attachments/assets/815afa53-0405-4fd0-a8fa-52bb63a7fe6c" />
申请列表界面：
<img width="1273" height="724" alt="image" src="https://github.com/user-attachments/assets/0ee8d879-f3bd-449c-a531-107ccb0e1be1" />
智能客服系统：
<img width="1277" height="740" alt="image" src="https://github.com/user-attachments/assets/833acee0-943c-47bf-8ff9-8e352a238ed7" />
人工客服界面：
<img width="1264" height="721" alt="image" src="https://github.com/user-attachments/assets/381cde67-b925-4257-9549-63ffe96ff5c7" />
个人中心界面
<img width="1277" height="688" alt="image" src="https://github.com/user-attachments/assets/0a1b7771-4051-4d00-a3bc-016ae1d539ca" />
管理员咨询界面
<img width="1275" height="715" alt="image" src="https://github.com/user-attachments/assets/ad0b0bcb-ef6d-4598-a1e3-3f7de51664b4" />
管理员仪表盘界面：
<img width="1279" height="727" alt="image" src="https://github.com/user-attachments/assets/c6c9a910-f25b-4e3e-be3f-c34d0f58343d" />
管理员审核界面：
<img width="1276" height="719" alt="image" src="https://github.com/user-attachments/assets/35440e79-33c6-4233-932b-db4fbc3cd204" />
管理员驿站管理界面
<img width="1279" height="757" alt="image" src="https://github.com/user-attachments/assets/da4b11be-f018-4efa-9082-f4062176462c" />
管理员知识库管理界面
<img width="1277" height="728" alt="image" src="https://github.com/user-attachments/assets/f7a57520-08e6-49e1-bfde-58779de8c2e3" />










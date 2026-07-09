package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/config"
	"yonth_station_backend/pkg/rag"

	ark "github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/zeromicro/go-zero/core/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// knowledge 知识库文档
type knowledge struct {
	Title    string
	Content  string
	Category string
}

var seedData = []knowledge{
	{
		Title:    "南昌市青年驿站入住须知",
		Category: "入住政策",
		Content: `南昌市青年人才驿站是为来南昌求职就业的青年提供免费短期住宿的公益项目。

入住条件：
1. 年龄在18-35周岁的全日制普通高校（含高职、大专、本科、硕士、博士）应往届毕业生
2. 身体健康，无重大疾病或传染病
3. 提供有效的身份证件和学历证明
4. 来南昌求职、创业或参加招聘活动的青年人才

入住时长：
- 应届毕业生最长可享受7天免费住宿
- 往届毕业生最长可享受5天免费住宿
- 确有特殊需要的可申请延长，最长不超过10天`,
	},
	{
		Title:    "青年驿站申请流程指南",
		Category: "申请流程",
		Content: `南昌市青年人才驿站申请流程如下：

第一步：注册账号
访问"云驿"平台，使用手机号或用户名注册账号，填写个人基本信息。

第二步：选择驿站
浏览可预约的驿站列表，查看各驿站的详细介绍、配套设施、剩余房间数等信息，选择心仪的驿站。

第三步：提交申请
选择计划入住日期和退房日期（最多7天），按照来访目的（求职/创业/研学）上传相应证明材料：
- 求职：上传面试通知邮件、招聘会邀请函等证明
- 创业：提交创业计划简介
- 研学：提供相关学术活动证明

第四步：等待审核
管理员将在1-2个工作日内完成审核，审核结果会通过平台通知。

第五步：缴纳押金
审核通过后需缴纳200元押金（退房时全额退还），可通过平台在线支付。

第六步：办理入住
凭审核通过通知和身份证件到驿站前台办理入住手续。`,
	},
	{
		Title:    "押金缴纳与退还说明",
		Category: "押金说明",
		Content: `青年驿站押金政策：

押金金额：200元/人

缴纳方式：
- 审核通过后通过平台在线支付
- 支持微信支付、支付宝支付
- 支付成功后平台会发送确认通知

退还条件：
- 按时退房，无损坏驿站设施
- 遵守驿站管理规定，无违规行为
- 退房后押金将在3个工作日内原路退还

押金扣除情况：
- 损坏驿站设施照价赔偿
- 违反驿站规定（如留宿他人、从事违法活动等）押金不予退还`,
	},
	{
		Title:    "驿站常见问题 FAQ",
		Category: "常见问题",
		Content: `青年驿站常见问题解答：

Q1：可以带家属或朋友一起入住吗？
A1：不可以。驿站房间为单人床位，仅限申请人本人入住，不得转借或留宿他人。

Q2：入住期间可以外出吗？
A2：可以。驿站不限制出入，但请注意安全和遵守驿站作息时间。

Q3：可以中途换驿站吗？
A3：原则上不可以。如有特殊情况需联系管理员协商处理。

Q4：申请被拒绝了怎么办？
A4：可以查看拒绝原因，补充完善材料后重新提交申请。

Q5：驿站提供餐饮吗？
A5：各驿站配套设施不同，部分驿站设有公共厨房可自行烹饪，部分驿站附近餐饮便利。

Q6：需要带什么生活用品？
A6：驿站提供床上用品（被褥、枕头），个人洗漱用品、换洗衣物等需自备。

Q7：可以续住吗？
A7：最长入住天数到期后一般不可续住。如有特殊情况可提前向管理员申请，视驿站房间情况而定。

Q8：如何联系人工客服？
A8：平台内点击"联系客服"即可在线咨询，或在驿站详情页查看驿站联系电话。`,
	},
	{
		Title:    "东湖青年人才驿站介绍",
		Category: "驿站介绍",
		Content: `东湖青年人才驿站位于南昌市东湖区八一大道388号，处于市中心繁华地段。

交通出行：
- 地铁1号线八一馆站步行约300米
- 多路公交线路经过，出行便利
- 距南昌火车站约3公里

配套设施：
- WiFi全覆盖、中央空调
- 独立卫浴、24小时热水
- 公共洗衣房（免费使用）
- 公共厨房（配备冰箱、微波炉、电磁炉）
- 图书角（提供求职、创业类书籍）
- 小型健身房（跑步机、哑铃等）

房间规格：
- 总房间数20间，标准单人间
- 每间约15-18平方米
- 配备单人床、书桌、衣柜、床头柜

周边环境：
- 步行可达八一广场、南昌百货大楼
- 周边餐饮、超市、银行等生活配套齐全`,
	},
	{
		Title:    "红谷滩青年人才驿站介绍",
		Category: "驿站介绍",
		Content: `红谷滩青年人才驿站位于南昌市红谷滩区凤凰中大道999号，处于南昌CBD核心区域。

交通出行：
- 地铁1号线卫东站步行约200米
- 多路公交线路经过
- 距南昌西站约8公里

配套设施：
- WiFi全覆盖、中央空调
- 独立卫浴、24小时热水
- 公共洗衣房（免费使用）
- 公共厨房
- 共享会议室（可预约使用）
- 咖啡厅（提供免费咖啡）
- 健身房

房间规格：
- 总房间数30间，标准单人间
- 每间约18-22平方米
- 配备单人床、书桌、衣柜、智能门锁

周边环境：
- 毗邻红谷滩万达广场
- 周边有南昌市政府、多家银行总部
- 靠近秋水广场、赣江市民公园`,
	},
	{
		Title:    "南昌市人才引进政策简介",
		Category: "常见问题",
		Content: `南昌市青年人才引进相关政策：

一、落户政策
全日制普通高校本科及以上学历毕业生，可凭毕业证直接办理南昌落户。

二、住房补贴
- 博士研究生：每月1500元租房补贴，最长享受3年
- 硕士研究生：每月1000元租房补贴，最长享受3年
- 本科生：每月500元租房补贴，最长享受2年

三、就业创业扶持
- 来南昌就业的应届毕业生可申请一次性就业补贴
- 创业青年可申请创业担保贷款，最高额度30万元
- 入驻创业孵化基地可享受场地租金减免

四、人才安居工程
- 符合条件的人才可申请购买人才住房
- 价格低于市场价，具体政策以当地住建部门公告为准

青年人才驿站是以上政策的配套服务之一，旨在为来南昌求职的青年提供过渡性住宿保障。`,
	},
	{
		Title:    "驿站住宿管理规定",
		Category: "入住政策",
		Content: `青年驿站住宿管理规定：

一、入住登记
- 凭身份证和审核通过通知办理入住
- 签署入住承诺书
- 领取房卡和床上用品

二、日常管理
- 保持房间整洁，爱护公共设施
- 不得在房间内吸烟、酗酒、赌博
- 不得留宿外来人员
- 不得从事违法违规活动
- 注意节约用水用电

三、安全管理
- 驿站设有门禁系统，凭房卡出入
- 公共区域24小时监控
- 严禁携带危险品入内
- 发现安全隐患及时报告管理人员

四、退房流程
- 退房时间为中午12:00前
- 交还房卡和床上用品
- 管理人员检查房间设施
- 确认无损后办理退房，退还押金

五、违规处理
- 首次违规：口头警告
- 再次违规：书面警告并记录
- 严重违规或多次违规：取消入住资格，押金不予退还`,
	},
}

func main() {
	configFile := flag.String("f", "../../etc/gateway-api.yaml", "config file path")
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 连接 MySQL
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect MySQL: %v", err)
	}
	db.AutoMigrate(&model.KnowledgeDoc{})

	ctx := context.Background()

	// 先删旧 collection（避免维度变更导致 schema 冲突）
	milvusClient, _ := milvusclient.New(ctx, &milvusclient.ClientConfig{Address: c.Chat.Milvus.Address})
	if err := milvusClient.DropCollection(ctx, milvusclient.NewDropCollectionOption(c.Chat.Milvus.Collection)); err != nil {
		log.Printf("[Milvus] Drop collection warning (may not exist): %v", err)
	}
	milvusClient.Close(ctx)

	// 创建 Indexer（多模态模型需指定 APIType）
	var embedApiType *ark.APIType
	if c.Chat.Embedding.UseMultiModal {
		t := ark.APITypeMultiModal
		embedApiType = &t
	}
	indexer, err := rag.NewIndexer(ctx, c.Chat.Milvus.Address, c.Chat.Milvus.Collection, c.Chat.Embedding.APIKey, c.Chat.Embedding.Model, embedApiType)
	if err != nil {
		log.Fatalf("Failed to create indexer: %v", err)
	}

	for i, k := range seedData {
		// 1. 保存到 MySQL
		doc := &model.KnowledgeDoc{
			Title:    k.Title,
			Content:  k.Content,
			Category: k.Category,
			Status:   0,
		}
		if err := db.Create(doc).Error; err != nil {
			log.Printf("[%d/%d] MySQL save failed: %v", i+1, len(seedData), err)
			continue
		}

		// 2. 向量化并存入 Milvus
		einoDoc := &schema.Document{
			ID:      fmt.Sprintf("%d", doc.Id),
			Content: k.Content,
			MetaData: map[string]interface{}{
				"title":    k.Title,
				"category": k.Category,
			},
		}
		_, err := rag.AddDocuments(ctx, indexer, []*schema.Document{einoDoc})
		if err != nil {
			log.Printf("[%d/%d] Milvus index failed: %v", i+1, len(seedData), err)
			continue
		}

		// 3. 更新状态为已索引
		db.Model(doc).Update("status", 1)
		log.Printf("[%d/%d] ✓ %s", i+1, len(seedData), k.Title)
	}

	fmt.Println("\n知识库预存完成！")
}

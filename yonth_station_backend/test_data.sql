-- ============================================
-- 青年驿站服务平台 - 测试数据
-- 使用方式：
--   1. 先启动后端 go run gateway.go（自动建表）
--   2. 再执行本 SQL：mysql -u root -p yonth_station < test_data.sql
--   3. 密码：
--      admin    / admin123  (管理员)
--      testuser / test123   (普通用户)
--      zhangsan / user123   (普通用户)
-- ============================================

-- 用户数据（bcrypt cost=12 真哈希，由项目 golang.org/x/crypto/bcrypt 生成）
INSERT INTO `user` (`user_name`, `phone`, `password`, `is_admin`, `status`, `gender`, `education`, `school`, `graduate_year`, `hukou_city`, `birth_date`, `last_login_at`, `created_at`, `updated_at`) VALUES
('admin',    '13800000001', '$2a$12$2kx0n8l1TwhUMf/ANfsukuHMMV1YA1ctmTzfHNUDzWUWpzsUxunwq', 1, 0, 1, 3, '南昌大学',     2020, '南昌市', '1998-06-15', 0, NOW(), NOW()),
('testuser', '13800000002', '$2a$12$AyBL.zmnbA7DKATqJ2NVhOaXzPEkh/2KjJrOiBG.fyfkoQQu3bqse', 0, 0, 2, 2, '江西师范大学', 2024, '赣州市', '2002-03-20', 0, NOW(), NOW()),
('zhangsan', '13912345678', '$2a$12$tpqhk63w6eTjxg7fyDmI8OVGidkTDBwBBMIKtxEsKM11fdBJAnzNi', 0, 0, 1, 2, '江西财经大学', 2024, '九江市', '2001-11-08', 0, NOW(), NOW());

-- 驿站数据（南昌各区）
INSERT INTO `station` (`station_name`, `district`, `address`, `latitude`, `longitude`, `contact_phone`, `business_hours`, `total_rooms`, `available_rooms`, `status`, `description`, `amenities`, `nearby_metro`, `image_url`, `weekly_quota`, `remaining_quota`, `avg_rating`, `total_reviews`, `created_at`, `updated_at`) VALUES
('东湖青年人才驿站', '东湖区', '南昌市东湖区八一大道388号', 28.6842, 115.8623, '0791-86800101', '8:30-18:00', 20, 18, 1, '东湖青年人才驿站位于南昌市中心城区，交通便利，周边配套设施齐全。提供标准化公寓式住宿，每间房配有独立卫浴、空调、无线网络。', '["WiFi","空调","独立卫浴","洗衣房","公共厨房","图书角","健身房"]', '地铁1号线-八一馆站', '/images/station-donghu.jpg', 50, 48, 4.5, 128, NOW(), NOW()),
('红谷滩青年人才驿站', '红谷滩区', '南昌市红谷滩区凤凰中大道999号', 28.7031, 115.8487, '0791-83800202', '9:00-17:30', 30, 25, 1, '红谷滩青年人才驿站位于南昌CBD核心区域，毗邻多家知名企业和创业园区。房间宽敞明亮，配备智能化门禁系统，为青年人才提供安全舒适的居住环境。', '["WiFi","空调","独立卫浴","洗衣房","公共厨房","会议室","咖啡厅","健身房"]', '地铁1号线-卫东站', '/images/station-honggutan.jpg', 80, 75, 4.7, 215, NOW(), NOW()),
('青山湖青年人才驿站', '青山湖区', '南昌市青山湖区北京东路555号', 28.6881, 115.9415, '0791-88100303', '8:00-18:00', 15, 12, 1, '青山湖青年人才驿站环境清幽，毗邻青山湖风景区。驿站注重人文关怀，定期组织青年交流活动，帮助来昌青年快速融入本地生活。', '["WiFi","空调","独立卫浴","洗衣房","公共厨房","阅览室"]', '地铁1号线-青山湖大道站', '/images/station-qingshanhu.jpg', 40, 35, 4.3, 89, NOW(), NOW()),
('西湖青年人才驿站', '西湖区', '南昌市西湖区中山路168号', 28.6723, 115.8812, '0791-86500404', '8:30-17:30', 25, 20, 1, '西湖青年人才驿站地处南昌历史文化核心区，紧邻滕王阁、万寿宫等文化景点。融合传统与现代设计风格，为青年人才提供独特的居住体验。', '["WiFi","空调","独立卫浴","洗衣房","公共厨房","自助餐厅","活动室"]', '地铁1号线-滕王阁站', '/images/station-xihu.jpg', 60, 55, 4.6, 176, NOW(), NOW()),
('高新青年人才驿站', '高新区', '南昌市高新区火炬大街266号', 28.6905, 116.0012, '0791-88150505', '9:00-18:00', 20, 15, 1, '高新青年人才驿站位于南昌高新技术产业开发区，紧邻多家科技企业和孵化器。专为科技人才打造，提供高速网络和共享办公空间。', '["WiFi","空调","独立卫浴","洗衣房","共享办公区","会议室","咖啡厅","健身房"]', '地铁1号线-艾溪湖东站', '/images/station-gaoxin.jpg', 50, 42, 4.4, 102, NOW(), NOW());

-- 房间数据
INSERT INTO `room` (`station_id`, `room_number`, `status`, `created_at`, `updated_at`) VALUES
(1, '101', 0, NOW(), NOW()), (1, '102', 0, NOW(), NOW()), (1, '103', 1, NOW(), NOW()), (1, '201', 0, NOW(), NOW()), (1, '202', 0, NOW(), NOW()),
(2, '101', 0, NOW(), NOW()), (2, '102', 0, NOW(), NOW()), (2, '103', 0, NOW(), NOW()), (2, '104', 1, NOW(), NOW()), (2, '201', 0, NOW(), NOW()),
(3, '101', 0, NOW(), NOW()), (3, '102', 1, NOW(), NOW()), (3, '201', 0, NOW(), NOW()),
(4, '101', 0, NOW(), NOW()), (4, '102', 0, NOW(), NOW()), (4, '103', 0, NOW(), NOW()), (4, '201', 1, NOW(), NOW()),
(5, '101', 0, NOW(), NOW()), (5, '102', 0, NOW(), NOW()), (5, '103', 0, NOW(), NOW());

-- 申请记录
INSERT INTO `application` (`user_id`, `station_id`, `room_id`, `checkin_date`, `checkout_date`, `status`, `visit_purpose`, `deposit_amount`, `deposit_status`, `checkin_at`, `checkout_at`, `audit_by`, `audit_at`, `applied_at`, `updated_at`, `created_at`) VALUES
(2, 1, 3, '2026-06-01', '2026-06-05', 5, 1, 20000, 2, 1717200000, 1717632000, 'admin', 1717120000, 1717000000, 1717632000, NOW()),
(2, 2, 0, '2026-06-20', '2026-06-25', 0, 2, 0, 0, 0, 0, '', 0, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(NOW()), NOW()),
(3, 4, 4, '2026-06-18', '2026-06-22', 1, 1, 20000, 1, 0, 0, 'admin', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(NOW())-86400, UNIX_TIMESTAMP(NOW()), NOW()),
(3, 3, 0, '2026-06-10', '2026-06-13', 2, 3, 0, 0, 0, 0, 'admin', UNIX_TIMESTAMP(NOW())-172800, UNIX_TIMESTAMP(NOW())-259200, UNIX_TIMESTAMP(NOW())-172800, NOW());

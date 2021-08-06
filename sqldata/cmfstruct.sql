-- phpMyAdmin SQL Dump
-- version 4.0.10.20
-- https://www.phpmyadmin.net
--
-- Host: localhost
-- Generation Time: Aug 06, 2021 at 02:06 PM
-- Server version: 5.5.62
-- PHP Version: 5.5.38

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- Database: `cmf`
--

-- --------------------------------------------------------

--
-- Table structure for table `eb_adlist`
--

CREATE TABLE IF NOT EXISTS `eb_adlist` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(200) NOT NULL DEFAULT '' COMMENT '信息标题',
  `adtype` varchar(100) NOT NULL DEFAULT '' COMMENT '推荐位置',
  `redtype` varchar(100) NOT NULL DEFAULT '' COMMENT '跳转类型',
  `redfunc` varchar(100) DEFAULT '' COMMENT '跳转模块',
  `redinfo` varchar(500) DEFAULT '' COMMENT '跳转值',
  `listsort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `text_text` text COMMENT '文本信息',
  `text_rich` text COMMENT '富文本信息',
  `imglist` text COMMENT '图片信息',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='推荐信息表' AUTO_INCREMENT=3 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_admin`
--

CREATE TABLE IF NOT EXISTS `eb_admin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(30) NOT NULL DEFAULT '' COMMENT '登录名',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '登录密码',
  `realname` varchar(30) NOT NULL DEFAULT '' COMMENT '人员姓名',
  `status` char(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `note` varchar(255) DEFAULT '' COMMENT '备注',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `last_time` datetime DEFAULT NULL COMMENT '登录时间',
  `last_ip` varchar(30) DEFAULT NULL COMMENT '登录ip',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `username` (`username`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='管理员表' AUTO_INCREMENT=3 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_admin_role`
--

CREATE TABLE IF NOT EXISTS `eb_admin_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `admin_id` int(11) NOT NULL COMMENT '用户ID',
  `role_id` int(11) NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `admin_id` (`admin_id`,`role_id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='用户和角色关联表' AUTO_INCREMENT=7 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_admin_struct`
--

CREATE TABLE IF NOT EXISTS `eb_admin_struct` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `admin_id` int(11) NOT NULL COMMENT '用户ID',
  `struct_id` int(11) NOT NULL COMMENT '组织ID',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='用户与组织架构关联表' AUTO_INCREMENT=10 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_adposition`
--

CREATE TABLE IF NOT EXISTS `eb_adposition` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `type` varchar(100) NOT NULL DEFAULT '' COMMENT '唯一标识',
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '位置名称',
  `note` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `width` mediumint(9) NOT NULL DEFAULT '0',
  `height` mediumint(9) NOT NULL DEFAULT '0',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name` (`type`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='推荐位置表' AUTO_INCREMENT=2 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_config`
--

CREATE TABLE IF NOT EXISTS `eb_config` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '配置名称',
  `type` varchar(30) NOT NULL DEFAULT '' COMMENT '配置标识',
  `style` varchar(10) NOT NULL DEFAULT '' COMMENT '配置类型',
  `is_sys` char(1) NOT NULL DEFAULT '0' COMMENT '是否系统内置 0否 1是',
  `groups` varchar(50) DEFAULT '' COMMENT '配置分组',
  `value` text COMMENT '配置值',
  `extra` varchar(255) DEFAULT '' COMMENT '配置项',
  `note` varchar(255) DEFAULT '' COMMENT '配置说明',
  `listsort` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='系统配置表' AUTO_INCREMENT=35 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_dict_data`
--

CREATE TABLE IF NOT EXISTS `eb_dict_data` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '字典编码',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '字典标签',
  `value` varchar(100) NOT NULL DEFAULT '' COMMENT '字典键值',
  `type` varchar(100) NOT NULL DEFAULT '' COMMENT '字典类型',
  `listsort` int(11) NOT NULL DEFAULT '0' COMMENT '字典排序',
  `status` char(1) NOT NULL DEFAULT '1' COMMENT '状态（1正常 0停用）',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `note` varchar(500) DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `value` (`value`) USING BTREE,
  KEY `type` (`type`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='字典数据表' AUTO_INCREMENT=7 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_dict_type`
--

CREATE TABLE IF NOT EXISTS `eb_dict_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '字典主键',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '字典名称',
  `type` varchar(100) NOT NULL DEFAULT '' COMMENT '字典类型',
  `status` char(1) NOT NULL DEFAULT '1' COMMENT '状态（1正常 0停用）',
  `listsort` int(11) NOT NULL DEFAULT '0',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `note` varchar(500) DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `dict_type` (`type`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='字典类型表' AUTO_INCREMENT=3 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_enter_log`
--

CREATE TABLE IF NOT EXISTS `eb_enter_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL DEFAULT '0',
  `enter_time` int(11) NOT NULL DEFAULT '0',
  `leave_time` int(11) NOT NULL DEFAULT '0',
  `store_id` int(11) NOT NULL DEFAULT '0',
  `door_id` int(11) DEFAULT NULL,
  `order_id` int(11) DEFAULT NULL,
  `service_id` int(11) NOT NULL DEFAULT '0',
  `status` int(2) DEFAULT NULL,
  `creatime` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=467 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_finance_info`
--

CREATE TABLE IF NOT EXISTS `eb_finance_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `type` int(11) NOT NULL,
  `code` varchar(50) NOT NULL,
  `name` varchar(50) NOT NULL,
  `bankname` varchar(255) DEFAULT NULL,
  `creatime` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=7 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_goods_cate`
--

CREATE TABLE IF NOT EXISTS `eb_goods_cate` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `sort` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=46118 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_loginlog`
--

CREATE TABLE IF NOT EXISTS `eb_loginlog` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '访问ID',
  `login_name` varchar(50) DEFAULT '' COMMENT '登录账号',
  `ipaddr` varchar(50) DEFAULT '' COMMENT '登录IP地址',
  `login_location` varchar(255) DEFAULT '' COMMENT '登录地点',
  `browser` varchar(100) DEFAULT '' COMMENT '浏览器类型',
  `os` varchar(100) DEFAULT '' COMMENT '操作系统',
  `net` varchar(50) DEFAULT '',
  `status` char(1) DEFAULT '0' COMMENT '登录状态（0成功 1失败）',
  `msg` varchar(255) DEFAULT '' COMMENT '提示消息',
  `login_time` datetime DEFAULT NULL COMMENT '访问时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='系统访问记录' AUTO_INCREMENT=895 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_mapply`
--

CREATE TABLE IF NOT EXISTS `eb_mapply` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '活动名称',
  `banner` text COMMENT '顶部Banner',
  `share_title` varchar(50) DEFAULT '' COMMENT '分享标题',
  `share_desc` varchar(150) DEFAULT '' COMMENT '分享简介',
  `share_img` varchar(255) DEFAULT '' COMMENT '分享图片',
  `money` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '预约金额',
  `rules` text COMMENT '活动规则',
  `agreement` text COMMENT '参与协议',
  `themecolor` varchar(10) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `is_multi` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否允许提交多个',
  `com_name` varchar(150) DEFAULT '' COMMENT '门店名称',
  `com_address` varchar(255) DEFAULT '' COMMENT '门店地址',
  `com_phone` varchar(50) DEFAULT '' COMMENT '门店电话',
  `com_lat` decimal(10,7) DEFAULT '0.0000000' COMMENT '纬度',
  `com_lng` decimal(10,7) DEFAULT '0.0000000' COMMENT '经度',
  `regfield` text,
  `start_time` datetime DEFAULT NULL COMMENT '开始预约时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='微预约-活动表' AUTO_INCREMENT=2 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_mapply_count`
--

CREATE TABLE IF NOT EXISTS `eb_mapply_count` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mid` int(11) NOT NULL DEFAULT '0' COMMENT '活动ID',
  `click` int(11) NOT NULL DEFAULT '0' COMMENT '点击次数',
  `number` int(11) NOT NULL DEFAULT '0' COMMENT '订单数量',
  `money` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '支付金额',
  `daytime` int(11) NOT NULL DEFAULT '0' COMMENT '时间',
  `paynumber` int(11) NOT NULL COMMENT '支付数量',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `pid` (`mid`,`daytime`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='微预约-活动统计' AUTO_INCREMENT=2 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_mapply_order`
--

CREATE TABLE IF NOT EXISTS `eb_mapply_order` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mid` int(11) NOT NULL DEFAULT '0' COMMENT '活动ID',
  `order_sn` varchar(50) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '订单号',
  `trade_no` varchar(100) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `openid` varchar(50) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户id',
  `money` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '金额',
  `is_pay` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否支付',
  `paytime` int(11) NOT NULL DEFAULT '0' COMMENT '支付时间',
  `order_title` varchar(100) DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否使用',
  `use_time` int(11) NOT NULL DEFAULT '0' COMMENT '使用时间',
  `user_name` varchar(20) DEFAULT '' COMMENT '姓名',
  `user_birthday` varchar(20) DEFAULT '' COMMENT '出生日期',
  `user_sex` varchar(10) DEFAULT '' COMMENT '用户性别',
  `user_mobile` varchar(20) DEFAULT '' COMMENT '手机号码',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `ip` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `order_sn` (`order_sn`) USING BTREE,
  KEY `pid` (`mid`) USING BTREE,
  KEY `openid` (`openid`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='微预约-预约订单' AUTO_INCREMENT=28 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_mapply_order_ext`
--

CREATE TABLE IF NOT EXISTS `eb_mapply_order_ext` (
  `mid` int(11) NOT NULL COMMENT '活动ID',
  `oid` int(11) NOT NULL COMMENT '订单ID',
  `fieldkey` varchar(30) NOT NULL DEFAULT '' COMMENT '字段标识',
  `fieldval` varchar(100) NOT NULL DEFAULT '' COMMENT '字段值'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='微预约-预约提交信息表';

-- --------------------------------------------------------

--
-- Table structure for table `eb_mapply_order_log`
--

CREATE TABLE IF NOT EXISTS `eb_mapply_order_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `order_id` int(11) NOT NULL COMMENT '订单ID',
  `title` varchar(100) NOT NULL COMMENT '操作名称',
  `optype` tinyint(1) NOT NULL COMMENT '操作用户 1用户 2 商户 3 管理员',
  `opname` varchar(50) NOT NULL COMMENT '操作人',
  `remark` varchar(255) NOT NULL,
  `mid` int(11) NOT NULL DEFAULT '0',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='微预约-订单操作记录' AUTO_INCREMENT=28 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_menu`
--

CREATE TABLE IF NOT EXISTS `eb_menu` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `name` varchar(30) NOT NULL DEFAULT '' COMMENT '菜单名称',
  `parent_id` int(11) NOT NULL DEFAULT '0' COMMENT '父菜单ID',
  `listsort` int(11) NOT NULL DEFAULT '0' COMMENT '显示顺序',
  `url` varchar(200) NOT NULL DEFAULT '' COMMENT '请求地址',
  `target` tinyint(1) NOT NULL DEFAULT '0' COMMENT '打开方式（0页签 1新窗口）',
  `type` char(1) NOT NULL DEFAULT '' COMMENT '菜单类型（M目录 C菜单 F按钮）',
  `status` char(1) NOT NULL DEFAULT '1' COMMENT '菜单状态（1显示 0隐藏）',
  `is_refresh` char(1) DEFAULT '0' COMMENT '是否刷新（0不刷新 1刷新）',
  `perms` varchar(100) DEFAULT '' COMMENT '权限标识',
  `icon` varchar(100) DEFAULT '' COMMENT '菜单图标',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `note` varchar(500) DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `parent_id` (`parent_id`) USING BTREE,
  KEY `listsort` (`listsort`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='菜单权限表' AUTO_INCREMENT=11836 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_notice`
--

CREATE TABLE IF NOT EXISTS `eb_notice` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '公告ID',
  `title` varchar(150) NOT NULL DEFAULT '' COMMENT '公告标题',
  `type` varchar(10) DEFAULT '' COMMENT '公告类型（1通知 2公告）',
  `content` text COMMENT '公告内容',
  `textarea` text COMMENT '非html内容',
  `status` char(1) NOT NULL DEFAULT '1' COMMENT '公告状态（1正常 0关闭）',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='通知公告表' AUTO_INCREMENT=3 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_opert_log`
--

CREATE TABLE IF NOT EXISTS `eb_opert_log` (
  `oper_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志主键',
  `title` varchar(50) DEFAULT '' COMMENT '模块标题',
  `business_type` int(11) DEFAULT '0' COMMENT '业务类型（0其它 1新增 2修改 3删除）',
  `method` varchar(100) DEFAULT '' COMMENT '方法名称',
  `request_method` varchar(10) DEFAULT '' COMMENT '请求方式',
  `operator_type` int(11) DEFAULT '0' COMMENT '操作类别（0其它 1后台用户 2手机端用户）',
  `oper_name` varchar(50) DEFAULT '' COMMENT '操作人员',
  `dept_name` varchar(50) DEFAULT '' COMMENT '部门名称',
  `oper_url` varchar(255) DEFAULT '' COMMENT '请求URL',
  `oper_ip` varchar(50) DEFAULT '' COMMENT '主机地址',
  `oper_location` varchar(255) DEFAULT '' COMMENT '操作地点',
  `oper_param` varchar(2000) DEFAULT '' COMMENT '请求参数',
  `json_result` varchar(2000) DEFAULT '' COMMENT '返回参数',
  `status` int(11) DEFAULT '0' COMMENT '操作状态（0正常 1异常）',
  `error_msg` varchar(2000) DEFAULT '' COMMENT '错误消息',
  `oper_time` datetime DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`oper_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='操作日志记录' AUTO_INCREMENT=1 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_pay_order`
--

CREATE TABLE IF NOT EXISTS `eb_pay_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `orderid` varchar(50) DEFAULT NULL,
  `transaction_id` varchar(50) DEFAULT NULL COMMENT '三方订单号',
  `paytype` int(2) DEFAULT NULL COMMENT '订单类型 1微信，2支付宝',
  `totalprice` int(11) DEFAULT NULL,
  `creatime` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_redtype`
--

CREATE TABLE IF NOT EXISTS `eb_redtype` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '名称',
  `type` varchar(100) DEFAULT '' COMMENT '跳转标识',
  `list_url` varchar(255) DEFAULT '' COMMENT '跳转模块连接',
  `info_url` varchar(255) DEFAULT '' COMMENT '跳转信息链接',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1',
  `note` varchar(255) DEFAULT NULL COMMENT '备注',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `adkey` (`type`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='跳转配置表' AUTO_INCREMENT=3 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_report`
--

CREATE TABLE IF NOT EXISTS `eb_report` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT '0' COMMENT '上报客服id',
  `cid` int(11) DEFAULT '0' COMMENT '用户id',
  `order_id` varchar(40) DEFAULT '' COMMENT '关联订单',
  `store_id` int(11) DEFAULT '0' COMMENT '店铺id',
  `mark` text COMMENT '备注',
  `creatime` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=7 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_role`
--

CREATE TABLE IF NOT EXISTS `eb_role` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(30) NOT NULL DEFAULT '' COMMENT '角色名称',
  `rolekey` varchar(50) NOT NULL DEFAULT '' COMMENT '角色权限字符串',
  `listsort` int(11) NOT NULL DEFAULT '0' COMMENT '显示顺序',
  `status` char(1) NOT NULL DEFAULT '1' COMMENT '角色状态（1正常 0停用）',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `note` varchar(500) DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `rolekey` (`rolekey`) USING BTREE,
  KEY `listsort` (`listsort`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='角色信息表' AUTO_INCREMENT=3 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_role_menu`
--

CREATE TABLE IF NOT EXISTS `eb_role_menu` (
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `menu_id` bigint(20) NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`role_id`,`menu_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='角色和菜单关联表';

-- --------------------------------------------------------

--
-- Table structure for table `eb_sell_detail`
--

CREATE TABLE IF NOT EXISTS `eb_sell_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) DEFAULT NULL,
  `sid` int(11) DEFAULT NULL COMMENT '小码id',
  `product_name` varchar(255) DEFAULT NULL,
  `order_id` varchar(40) DEFAULT NULL,
  `store_id` int(11) DEFAULT '0',
  `uid` int(11) DEFAULT NULL,
  `price` int(11) DEFAULT NULL,
  `total_price` int(11) DEFAULT NULL COMMENT '总价',
  `nums` double(11,2) NOT NULL DEFAULT '0.00',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '订单状态，0未支付，1支付',
  `pay_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '支付类型，1微信，2支付宝',
  `creatime` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=269 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_service`
--

CREATE TABLE IF NOT EXISTS `eb_service` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL,
  `creatime` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_settlement`
--

CREATE TABLE IF NOT EXISTS `eb_settlement` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `day` varchar(30) DEFAULT NULL,
  `sellcounts` int(11) DEFAULT NULL,
  `ordercounts` int(11) DEFAULT NULL COMMENT '订单总量',
  `ordersuccess` int(11) DEFAULT NULL COMMENT '成功订单量',
  `visitcounts` int(11) DEFAULT NULL COMMENT '总访问量',
  `trademoney` bigint(20) DEFAULT NULL COMMENT '总成交额',
  `ordermoney` bigint(20) DEFAULT NULL COMMENT '总订单金额',
  `uptime` int(11) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `day` (`day`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COMMENT='每日结算表' AUTO_INCREMENT=97 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_smscode`
--

CREATE TABLE IF NOT EXISTS `eb_smscode` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号码',
  `code` varchar(20) NOT NULL DEFAULT '' COMMENT '验证码',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '例如：1注册 2登录 3忘记密码',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0未验证 1已验证',
  `os` varchar(20) NOT NULL DEFAULT '' COMMENT '运营商',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='验证码表' AUTO_INCREMENT=1 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_staff`
--

CREATE TABLE IF NOT EXISTS `eb_staff` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL COMMENT '员工用户id',
  `mer_id` int(11) NOT NULL COMMENT '商户id',
  `store_id` int(11) NOT NULL COMMENT '店铺id',
  `creatime` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `store_uid` (`uid`,`store_id`) COMMENT '单人单铺'
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=2 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_staff_mer`
--

CREATE TABLE IF NOT EXISTS `eb_staff_mer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `mer_id` int(11) NOT NULL,
  `creatime` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=2 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_stock_detail`
--

CREATE TABLE IF NOT EXISTS `eb_stock_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `bar_code` int(11) NOT NULL COMMENT '条形码',
  `pid` int(11) DEFAULT NULL COMMENT '商品id',
  `pname` varchar(255) NOT NULL COMMENT '商品名称',
  `store_id` int(11) DEFAULT NULL,
  `counts` int(11) NOT NULL COMMENT '数量',
  `creatime` int(11) NOT NULL COMMENT '时间',
  `uid` int(11) DEFAULT '0' COMMENT '操作用户的id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=5 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_store`
--

CREATE TABLE IF NOT EXISTS `eb_store` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL DEFAULT '0',
  `location` varchar(255) DEFAULT NULL,
  `lng` varchar(30) DEFAULT '' COMMENT '经度',
  `lat` varchar(30) NOT NULL DEFAULT '' COMMENT '纬度',
  `name` varchar(255) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `phone` varchar(30) DEFAULT NULL,
  `rate` int(11) DEFAULT NULL,
  `closed` int(2) NOT NULL DEFAULT '0' COMMENT '营业状态 0开业，1关闭中',
  `creatime` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=6 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_store_category`
--

CREATE TABLE IF NOT EXISTS `eb_store_category` (
  `id` mediumint(11) NOT NULL AUTO_INCREMENT COMMENT '商品分类表ID',
  `pid` mediumint(11) NOT NULL COMMENT '父id',
  `cate_name` varchar(100) NOT NULL COMMENT '分类名称',
  `sort` mediumint(11) NOT NULL COMMENT '排序',
  `pic` varchar(128) NOT NULL DEFAULT '' COMMENT '图标',
  `is_show` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否推荐',
  `add_time` int(11) NOT NULL COMMENT '添加时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `pid` (`pid`) USING BTREE,
  KEY `is_base` (`is_show`) USING BTREE,
  KEY `sort` (`sort`) USING BTREE,
  KEY `add_time` (`add_time`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='商品分类表' AUTO_INCREMENT=54 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_store_device`
--

CREATE TABLE IF NOT EXISTS `eb_store_device` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `store_id` int(11) DEFAULT NULL,
  `nums` int(2) DEFAULT NULL,
  `devicesn` varchar(20) DEFAULT NULL,
  `describe` varchar(255) DEFAULT NULL,
  `types` int(11) DEFAULT NULL COMMENT '功能类型',
  `token` varchar(40) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=10 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_store_order`
--

CREATE TABLE IF NOT EXISTS `eb_store_order` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_id` varchar(32) NOT NULL COMMENT '订单号',
  `out_trade_sn` varchar(100) DEFAULT NULL,
  `uid` int(11) unsigned NOT NULL COMMENT '用户id',
  `total_num` double(11,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '订单商品总数',
  `total_price` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '订单总价',
  `pay_price` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '实际支付金额',
  `deduction_price` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '抵扣金额',
  `coupon_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券id',
  `coupon_price` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券金额',
  `paid` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '支付状态',
  `pay_time` int(11) unsigned DEFAULT NULL COMMENT '支付时间',
  `pay_type` varchar(32) NOT NULL COMMENT '支付方式',
  `add_time` int(11) unsigned NOT NULL COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '订单状态（-1 : 申请退款 -2 : 退货成功 0：待发货；1：待收货；2：已收货；3：待评价；-1：已退款）',
  `mark` varchar(512) NOT NULL COMMENT '备注',
  `is_del` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  `remark` varchar(512) DEFAULT NULL COMMENT '管理员备注',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户ID',
  `store_id` int(11) NOT NULL DEFAULT '0' COMMENT '门店id',
  `eid` int(11) DEFAULT NULL COMMENT '记录id',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `order_id_2` (`order_id`,`uid`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE,
  KEY `add_time` (`add_time`) USING BTREE,
  KEY `pay_price` (`pay_price`) USING BTREE,
  KEY `paid` (`paid`) USING BTREE,
  KEY `pay_time` (`pay_time`) USING BTREE,
  KEY `pay_type` (`pay_type`) USING BTREE,
  KEY `status` (`status`) USING BTREE,
  KEY `is_del` (`is_del`) USING BTREE,
  KEY `coupon_id` (`coupon_id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='订单表' AUTO_INCREMENT=205 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_store_product`
--

CREATE TABLE IF NOT EXISTS `eb_store_product` (
  `id` mediumint(11) NOT NULL AUTO_INCREMENT COMMENT '商品id',
  `store_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '商户Id(0为总后台管理员创建,不为0的时候是商户后台创建)',
  `image` varchar(256) DEFAULT NULL COMMENT '商品图片',
  `pro_name` varchar(128) NOT NULL COMMENT '商品名称',
  `pro_info` varchar(256) NOT NULL COMMENT '商品简介',
  `keyword` varchar(256) NOT NULL DEFAULT '' COMMENT '关键字',
  `sbarcode` varchar(15) DEFAULT NULL COMMENT '大码 可以为0',
  `sid` int(11) DEFAULT NULL,
  `exchange` int(11) DEFAULT '0',
  `bar_code` varchar(15) NOT NULL DEFAULT '' COMMENT '商品条码（一维码）',
  `cate_id` int(11) NOT NULL DEFAULT '0' COMMENT '分类id',
  `cost` int(11) unsigned NOT NULL COMMENT '成本价',
  `price` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '商品价格',
  `unit_name` varchar(32) NOT NULL COMMENT '单位名',
  `location` varchar(255) DEFAULT NULL COMMENT '货架位置',
  `sort` smallint(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `sales` double(11,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '销量',
  `stock` double(11,2) NOT NULL DEFAULT '0.00' COMMENT '库存',
  `add_time` int(11) unsigned NOT NULL COMMENT '添加时间',
  `up_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `is_del` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `storecode` (`bar_code`,`store_id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='商品表' AUTO_INCREMENT=8682 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_store_video`
--

CREATE TABLE IF NOT EXISTS `eb_store_video` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `store_id` int(11) NOT NULL DEFAULT '0',
  `describe` varchar(255) DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL,
  `isprimary` int(2) DEFAULT NULL COMMENT '1主摄像头',
  `appid` varchar(40) DEFAULT NULL,
  `secret` varchar(50) DEFAULT NULL,
  `accesstoken` varchar(80) DEFAULT NULL,
  `expiretime` varchar(20) DEFAULT NULL,
  `safekey` varchar(20) DEFAULT NULL,
  `devicesn` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=18 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_struct`
--

CREATE TABLE IF NOT EXISTS `eb_struct` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '部门id',
  `name` varchar(30) DEFAULT '' COMMENT '部门名称',
  `parent_id` int(11) DEFAULT '0' COMMENT '父部门id',
  `levels` varchar(100) DEFAULT '' COMMENT '祖级列表',
  `listsort` int(11) DEFAULT '0' COMMENT '显示顺序',
  `leader` varchar(20) DEFAULT NULL COMMENT '负责人',
  `phone` varchar(11) DEFAULT NULL COMMENT '联系电话',
  `note` varchar(255) DEFAULT '' COMMENT '备注',
  `status` char(1) DEFAULT '1' COMMENT '部门状态（1正常 0停用）',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='组织架构' AUTO_INCREMENT=112 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_system_admin`
--

CREATE TABLE IF NOT EXISTS `eb_system_admin` (
  `id` smallint(5) unsigned NOT NULL AUTO_INCREMENT COMMENT '后台管理员表ID',
  `account` varchar(32) NOT NULL COMMENT '后台管理员账号',
  `pwd` char(32) NOT NULL COMMENT '后台管理员密码',
  `real_name` varchar(16) NOT NULL COMMENT '后台管理员姓名',
  `roles` varchar(128) NOT NULL COMMENT '后台管理员权限(menus_id)',
  `last_ip` varchar(16) DEFAULT NULL COMMENT '后台管理员最后一次登录ip',
  `last_time` int(10) unsigned DEFAULT NULL COMMENT '后台管理员最后一次登录时间',
  `add_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '后台管理员添加时间',
  `login_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '登录次数',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '后台管理员级别',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '后台管理员状态 1有效0无效',
  `is_del` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `account` (`account`) USING BTREE,
  KEY `status` (`status`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='后台管理员表' AUTO_INCREMENT=3 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_system_log`
--

CREATE TABLE IF NOT EXISTS `eb_system_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '管理员操作记录ID',
  `admin_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '管理员id',
  `admin_name` varchar(64) NOT NULL DEFAULT '' COMMENT '管理员姓名',
  `path` varchar(128) NOT NULL DEFAULT '' COMMENT '链接',
  `page` varchar(64) NOT NULL DEFAULT '' COMMENT '行为',
  `method` varchar(12) NOT NULL DEFAULT '' COMMENT '访问类型',
  `ip` varchar(16) NOT NULL DEFAULT '' COMMENT '登录IP',
  `type` varchar(32) NOT NULL DEFAULT '' COMMENT '类型',
  `add_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '操作时间',
  `merchant_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户id',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `admin_id` (`admin_id`) USING BTREE,
  KEY `add_time` (`add_time`) USING BTREE,
  KEY `type` (`type`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='管理员操作记录表' AUTO_INCREMENT=2166 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_test_online`
--

CREATE TABLE IF NOT EXISTS `eb_test_online` (
  `ip` varchar(255) NOT NULL DEFAULT '1',
  `create_time` datetime DEFAULT NULL,
  `fd` varchar(255) DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `isrun` tinyint(1) DEFAULT '1',
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC AUTO_INCREMENT=5 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_user`
--

CREATE TABLE IF NOT EXISTS `eb_user` (
  `uid` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `account` varchar(32) NOT NULL COMMENT '用户账号',
  `balance` int(11) NOT NULL DEFAULT '0',
  `pwd` varchar(32) NOT NULL COMMENT '用户密码',
  `token` varchar(32) DEFAULT NULL,
  `openid` char(60) DEFAULT NULL,
  `accesstoken` varchar(100) DEFAULT NULL COMMENT '支付宝token',
  `nickname` varchar(60) NOT NULL COMMENT '用户昵称',
  `avatar` varchar(256) DEFAULT NULL COMMENT '用户头像',
  `phone` char(15) DEFAULT NULL COMMENT '手机号码',
  `add_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `add_ip` varchar(16) NOT NULL DEFAULT '' COMMENT '添加ip',
  `last_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最后一次登录时间',
  `last_ip` varchar(16) NOT NULL DEFAULT '' COMMENT '最后一次登录ip',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1为正常，0为禁止',
  `isleave` int(2) DEFAULT NULL COMMENT '是否暂时离开',
  `online` int(2) DEFAULT NULL COMMENT '0下线1上线',
  `user_type` int(2) NOT NULL DEFAULT '0' COMMENT '用户类型1，2普通4客服3商家5员工',
  `alipayaccount` varchar(30) NOT NULL DEFAULT '',
  `pay_count` int(11) unsigned DEFAULT '0' COMMENT '用户购买次数',
  `login_type` varchar(36) NOT NULL DEFAULT '' COMMENT '用户登陆类型，h5,wechat,routine',
  PRIMARY KEY (`uid`) USING BTREE,
  UNIQUE KEY `account` (`account`) USING BTREE,
  UNIQUE KEY `openid` (`openid`) USING BTREE,
  KEY `status` (`status`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='用户表' AUTO_INCREMENT=116 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_user_account`
--

CREATE TABLE IF NOT EXISTS `eb_user_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL DEFAULT '0' COMMENT '用户id',
  `crystal` int(11) DEFAULT '0' COMMENT '水晶余额',
  `task_value` int(11) DEFAULT '0' COMMENT '任务值',
  `creatime` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `savetime` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`uid`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='用户余额表' AUTO_INCREMENT=12 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_user_bill`
--

CREATE TABLE IF NOT EXISTS `eb_user_bill` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户账单id',
  `uid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户uid',
  `store_id` int(11) DEFAULT NULL,
  `order_id` varchar(50) DEFAULT NULL COMMENT '订单号',
  `type` varchar(64) NOT NULL DEFAULT '' COMMENT '明细类型',
  `balance` int(11) NOT NULL DEFAULT '0' COMMENT '金额',
  `mark` varchar(512) NOT NULL DEFAULT '' COMMENT '备注',
  `add_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `openid` (`uid`) USING BTREE,
  KEY `add_time` (`add_time`) USING BTREE,
  KEY `type` (`type`) USING BTREE,
  KEY `order_id` (`order_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='用户账单表' AUTO_INCREMENT=329 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_user_log`
--

CREATE TABLE IF NOT EXISTS `eb_user_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL,
  `mark` varchar(255) DEFAULT NULL,
  `usertype` int(2) DEFAULT NULL,
  `types` varchar(255) DEFAULT NULL,
  `creatime` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=5665 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_wall`
--

CREATE TABLE IF NOT EXISTS `eb_wall` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL COMMENT '活动名称',
  `password` varchar(20) DEFAULT '' COMMENT '大屏密码',
  `bgimg` varchar(255) DEFAULT '' COMMENT '背景图片',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `isopen` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否开启',
  `contents` text,
  `logoimg` varchar(255) DEFAULT '' COMMENT 'logo',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='微现场\r\n' AUTO_INCREMENT=2 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_wall_prize`
--

CREATE TABLE IF NOT EXISTS `eb_wall_prize` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `wall_id` int(11) NOT NULL DEFAULT '0' COMMENT '活动ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '奖品等级',
  `thumbimg` varchar(255) DEFAULT '' COMMENT '奖品名称',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '奖品图片',
  `number` mediumint(9) NOT NULL DEFAULT '0',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='微现场-奖品表' AUTO_INCREMENT=5 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_wall_prize_users`
--

CREATE TABLE IF NOT EXISTS `eb_wall_prize_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `openid` varchar(100) NOT NULL DEFAULT '',
  `truename` varchar(50) DEFAULT '',
  `headimg` varchar(255) DEFAULT '',
  `mobile` varchar(30) DEFAULT '',
  `prize_id` int(11) NOT NULL,
  `wall_id` int(11) NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='微现场-中奖用户表' AUTO_INCREMENT=157 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_wall_process`
--

CREATE TABLE IF NOT EXISTS `eb_wall_process` (
  `id` int(10) unsigned zerofill NOT NULL AUTO_INCREMENT,
  `wall_id` int(11) NOT NULL COMMENT '活动ID',
  `daytime` date DEFAULT NULL COMMENT '日程日期',
  `title` varchar(30) DEFAULT '' COMMENT '日程标题',
  `desc` varchar(255) DEFAULT '' COMMENT '日程详情',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `listsort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `hour` varchar(20) DEFAULT '' COMMENT '日程日期',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='微现场-活动日程' AUTO_INCREMENT=28 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_wall_users`
--

CREATE TABLE IF NOT EXISTS `eb_wall_users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `wall_id` int(11) NOT NULL DEFAULT '0' COMMENT '活动ID',
  `openid` varchar(100) NOT NULL DEFAULT '' COMMENT '微信标识openid',
  `headimg` varchar(255) DEFAULT '' COMMENT '微信头像',
  `sex` tinyint(1) DEFAULT '0' COMMENT '性别',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否有效',
  `truename` varchar(20) DEFAULT '' COMMENT '真实姓名',
  `mobile` varchar(20) DEFAULT NULL COMMENT '手机号码',
  `update_time` datetime DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='活动报名用户' AUTO_INCREMENT=1469 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_web_ad`
--

CREATE TABLE IF NOT EXISTS `eb_web_ad` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(200) NOT NULL DEFAULT '' COMMENT '信息标题',
  `pos_id` int(11) NOT NULL DEFAULT '0' COMMENT '推荐位置',
  `linkurl` varchar(255) DEFAULT NULL,
  `listsort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `text_text` text COMMENT '文本信息',
  `text_rich` text COMMENT '富文本信息',
  `imglist` text COMMENT '图片信息',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `website` varchar(50) DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='推荐信息表' AUTO_INCREMENT=4 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_web_cat`
--

CREATE TABLE IF NOT EXISTS `eb_web_cat` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(20) NOT NULL DEFAULT '' COMMENT '菜单标题，后台',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '菜单名称 前台',
  `website` int(11) NOT NULL DEFAULT '0' COMMENT '鎵€灞炵珯鐐?',
  `parent_id` int(11) NOT NULL DEFAULT '0' COMMENT '父级菜单',
  `listsort` int(11) NOT NULL DEFAULT '0' COMMENT '显示排序',
  `type` varchar(10) NOT NULL DEFAULT '' COMMENT 'album相册，list文章列表，page单页，link外链,',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `lang` varchar(10) DEFAULT '' COMMENT '语言',
  `url` varchar(255) DEFAULT '' COMMENT '外链地址',
  `template_list` varchar(30) DEFAULT '' COMMENT '列表模板',
  `template_info` varchar(30) DEFAULT '' COMMENT '详情模板',
  `checkcode` varchar(20) DEFAULT '' COMMENT '选中菜单的标识',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='网站-菜单分类' AUTO_INCREMENT=10 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_web_list`
--

CREATE TABLE IF NOT EXISTS `eb_web_list` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(150) NOT NULL DEFAULT '' COMMENT '标题',
  `remark` varchar(200) DEFAULT '' COMMENT '简介',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `author` varchar(50) DEFAULT '' COMMENT '作者',
  `froms` varchar(50) DEFAULT '' COMMENT '来源',
  `thumbimg` varchar(255) DEFAULT '' COMMENT '缩略图',
  `catid` int(11) NOT NULL DEFAULT '0' COMMENT '所属菜单ID',
  `website` int(11) NOT NULL DEFAULT '0' COMMENT '所属站点',
  `click` int(11) DEFAULT '0',
  `linkurl` varchar(255) DEFAULT '' COMMENT '外链地址',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `subtime` datetime DEFAULT NULL COMMENT '发布时间',
  PRIMARY KEY (`id`,`website`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='网站-信息列表' AUTO_INCREMENT=15 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_web_list_ext`
--

CREATE TABLE IF NOT EXISTS `eb_web_list_ext` (
  `id` int(10) unsigned NOT NULL,
  `content` text,
  `imglist` text,
  `website` int(11) DEFAULT '0',
  `catid` int(11) DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='网站-信息列表其他信息';

-- --------------------------------------------------------

--
-- Table structure for table `eb_web_pos`
--

CREATE TABLE IF NOT EXISTS `eb_web_pos` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `website` int(11) NOT NULL DEFAULT '0' COMMENT '鎵€灞炵珯鐐?',
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '位置名称',
  `note` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `width` mediumint(9) NOT NULL DEFAULT '0',
  `height` mediumint(9) NOT NULL DEFAULT '0',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='推荐位置表' AUTO_INCREMENT=5 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_web_site`
--

CREATE TABLE IF NOT EXISTS `eb_web_site` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL DEFAULT '' COMMENT '标识',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '站点标题（后台）',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '站点名称（前台）',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `is_default` char(1) DEFAULT '0',
  `template` varchar(30) DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='网站-站点管理' AUTO_INCREMENT=4 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_wechat_access`
--

CREATE TABLE IF NOT EXISTS `eb_wechat_access` (
  `appid` varchar(50) NOT NULL DEFAULT '',
  `access_token` varchar(255) NOT NULL DEFAULT '',
  `jsapi_ticket` varchar(255) NOT NULL DEFAULT '',
  `access_token_add` int(11) NOT NULL DEFAULT '0',
  `jsapi_ticket_add` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`appid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='微信jsapi和access';

-- --------------------------------------------------------

--
-- Table structure for table `eb_wechat_users`
--

CREATE TABLE IF NOT EXISTS `eb_wechat_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `openid` varchar(100) NOT NULL DEFAULT '' COMMENT '唯一标识',
  `appid` varchar(50) NOT NULL DEFAULT '' COMMENT '公众号参数',
  `nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '昵称',
  `headimg` varchar(255) NOT NULL DEFAULT '' COMMENT '头像地址',
  `update_time` datetime DEFAULT NULL COMMENT '资料更新时间',
  `create_time` datetime DEFAULT NULL COMMENT '添加时间',
  `sex` tinyint(1) NOT NULL DEFAULT '0' COMMENT '性别',
  `city` varchar(50) NOT NULL DEFAULT '' COMMENT '城市',
  `country` varchar(50) NOT NULL DEFAULT '' COMMENT '国家',
  `province` varchar(50) NOT NULL DEFAULT '' COMMENT '省份',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `type` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `openid` (`openid`,`type`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='微信用户信息表' AUTO_INCREMENT=2 ;

-- --------------------------------------------------------

--
-- Table structure for table `eb_withdraw`
--

CREATE TABLE IF NOT EXISTS `eb_withdraw` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` varchar(80) DEFAULT NULL,
  `uid` int(11) DEFAULT NULL,
  `money` int(11) NOT NULL DEFAULT '0',
  `creatime` int(11) NOT NULL DEFAULT '0',
  `uptime` int(11) NOT NULL,
  `types` int(11) DEFAULT NULL,
  `status` int(2) DEFAULT NULL,
  `mark` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=8 ;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

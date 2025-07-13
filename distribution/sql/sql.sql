create database if not exists blog_narcissus character set utf8mb4 collate utf8mb4_general_ci;
use blog_narcissus;

CREATE TABLE users (
    id INT AUTO_INCREMENT COMMENT '用户ID',
    username VARCHAR(30) NOT NULL COMMENT '用户名',
    nickname VARCHAR(30) NOT NULL COMMENT '昵称',
    password CHAR(60) NOT NULL COMMENT '密码',
    email VARCHAR(100) COMMENT '邮箱',
    phone_number VARCHAR(20) COMMENT '手机号',
    avatar_path VARCHAR(255) COMMENT '头像路径',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY(id),
    UNIQUE KEY(username),
    UNIQUE KEY(nickname)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE roles (
    id INT AUTO_INCREMENT COMMENT '角色ID',
    role TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '0:普通用户;1:系统管理员',
    PRIMARY KEY(id),
    UNIQUE(role)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE `articles` (
    `id` INT AUTO_INCREMENT COMMENT '文章ID',
    `title` VARCHAR(255) NOT NULL COMMENT '文章标题',
    `summary` TEXT NOT NULL COMMENT '摘要',
    `type` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章类型，0表示博客文章，1表示随笔，2表示关于',
    `category_id` INT COMMENT '文章分类ID，每篇文章最多1个分类，可以为空',
    `author` VARCHAR(50) NOT NULL COMMENT '作者姓名或标识',
    `allow_comment` TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否允许评论，0表示不允许评论，1表示允许评论',
    `weight` INT NOT NULL DEFAULT 0 COMMENT '文章权重，默认初始值为0',
    `is_sticky` TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否置顶。0表示不置顶，1表示置顶。默认初始值为 0',
    `is_original` TINYINT(1) UNSIGNED NOT NULL DEFAULT 1 COMMENT '原创/转载标识。0表示非原创，1表示原创。默认初始值为1，表示原创',
    `original_article_link` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '转载文章的原始文章链接',
    `status` TINYINT(1) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态，0表示offline，1表示online',
    `created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY (`category_id`),
    KEY (`type`)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE `article_content` (
    `id` INT AUTO_INCREMENT COMMENT '文章内容ID', -- 文章内容ID
    `article_id` INT NOT NULL COMMENT '文章ID', -- 文章ID
    `content` TEXT NOT NULL COMMENT '文章内容',  -- 文章内容
    PRIMARY KEY(id),
    UNIQUE KEY(article_id)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE `article_categories` (
    `id` INT AUTO_INCREMENT COMMENT '分类ID',  -- 分类ID
    `name` VARCHAR(20) NOT NULL COMMENT '分类名称',  -- 分类名称
    `created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',  -- 创建时间
    `updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',  -- 更新时间
    PRIMARY KEY(id),
    UNIQUE KEY(name)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE `article_tags` (
    `id` INT AUTO_INCREMENT COMMENT '标签ID',  -- 标签ID
    `name` VARCHAR(20) NOT NULL COMMENT '标签名称',  -- 标签名称
    `created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',  -- 创建时间
    `updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',  -- 更新时间
    PRIMARY KEY(id),
    UNIQUE KEY(name)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE `article_tag_relations` (
    `id` INT AUTO_INCREMENT COMMENT '文章标签关系ID',
    `article_id` INT NOT NULL COMMENT '文章ID',
    `tag_id` INT NOT NULL COMMENT '标签ID',
    PRIMARY KEY (`id`),
    KEY (`article_id`),
    KEY (`tag_id`)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

INSERT INTO article_categories(name) 
VALUES
("Java"),
("Golang"),
("Python"),
("Kubernetes"),
("Docker"),
("RocketMQ"),
("Linux"),
("Nginx"),
("AI");

INSERT INTO article_tags(name) 
VALUES
("基础"),
("网络");

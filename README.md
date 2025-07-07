# Narcissus Blog

dsn : "root:root@tcp(172.26.21.6:3306)/blog_narcissus?charset=utf8mb4&parseTime=true&loc=Local"

gentool -dsn "root:root@tcp(172.24.12.92:3306)/blog_narcissus?charset=utf8mb4&parseTime=true&loc=Local" -tables "article_categories,article_content,article_tags,articles,article_tag_relations"

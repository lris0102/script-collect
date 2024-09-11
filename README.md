Nginx 日志安全分析器

    功能：
        SQL 注入攻击检测
        常见扫描工具和黑客工具检测
        敏感文件访问检测
        Webshell 行为分析
        IP 地址统计与日志分析报告

    使用方法：
        配置 Nginx 日志目录（默认 /var/log/nginx/），通过执行脚本分析日志并输出报告到 /tmp/logs/ 目录。
        运行脚本：

        bash

        chmod +x nginx_log_security_analyzer.sh
        ./nginx_log_security_analyzer.sh

内网资产扫描器

    功能：
        扫描整个内网网段，发现在线主机。
        扫描开放端口，并识别常见服务（如 HTTP、SSH、MySQL 等）。
        通过 SSH 登录主机，检查系统日志，判断是否存在入侵行为。

    使用方法：
        修改子网范围，运行 Go 代码进行扫描：

        bash

        go build -o network_scanner
        ./network_scanner

        该工具会扫描内网中所有主机，并输出开放端口、服务类型及入侵检测结果。

适用场景

    个人学习和实验环境，用于提升日志分析和内网资产扫描的安全技能。
    非生产环境中的初步日志和资产扫描分析，快速发现潜在的安全问题。

注意事项

    Nginx 日志安全分析器：建议用于 ELK 或 Splunk 之前的初步分析。
    内网资产扫描器：适用于内网测试环境，不建议用于大规模或生产网络。

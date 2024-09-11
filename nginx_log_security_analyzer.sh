#!/usr/bin/env bash

echo ""
echo " ========================================================= "
echo " \                 Nginx 日志安全分析器 V1.0.0           / "
echo " ========================================================= "
echo " # 支持Nginx日志分析，攻击告警分析等                    "
echo -e "\n"

# 设置分析结果存储目录,结尾不能加/
outfile=/tmp/logs
# 设置nginx日志目录，结尾必须加/
access_dir=/var/log/nginx/
# 设置文件名，如果文件名为access，那么匹配的是access*文件
access_log=access

# 创建结果输出目录
function setup_output_dir() {
    if [ -d "$outfile" ]; then
        rm -rf "$outfile"/*
    else
        mkdir -p "$outfile"
    fi
}

# 检查日志目录及文件
function check_log_files() {
    num=$(ls "${access_dir}${access_log}"* | wc -l) >/dev/null 2>&1
    if [ "$num" -eq 0 ]; then
        echo '日志文件不存在'
        exit 1
    fi
}

# 检查操作系统类型
function check_os() {
    OS='None'
    if [ -e "/etc/os-release" ]; then
        source /etc/os-release
        case ${ID} in
        "debian" | "ubuntu" | "devuan")
            OS='Debian'
            ;;
        "centos" | "rhel" | "fedora")
            OS='Centos'
            ;;
        esac
    fi

    if [ "$OS" = 'None' ]; then
        if command -v apt-get >/dev/null 2>&1; then
            OS='Debian'
        elif command -v yum >/dev/null 2>&1; then
            OS='Centos'
        else
            echo -e "\n不支持这个系统\n"
            exit 1
        fi
    fi
}

# 安装依赖：the_silver_searcher (ag)
function install_ag() {
    if command -v ag >/dev/null 2>&1; then
        echo -e "\e[00;32msilversearcher-ag 已安装 \e[00m"
    else
        if [ "$OS" = 'Centos' ]; then
            yum -y install the_silver_searcher >/dev/null 2>&1
        else
            apt-get -y install silversearcher-ag >/dev/null 2>&1
        fi
    fi
}

# 输出文件统计
function print_stat() {
    awk '{print $1 " 次"}' "$1" | tail -n1
}

# 统计TOP 20地址
function analyze_top_ip() {
    echo -e "\e[00;31m[+]TOP 20 IP 地址\e[00m"
    ag -a -o --nofilename '((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}' "${access_dir}${access_log}"* | sort | uniq -c | sort -nr | head -n 20 | tee -a "${outfile}/top20.log"
    echo -e "\n"
}

# SQL注入分析
function analyze_sql_injection() {
    echo -e "\e[00;31m[+]SQL注入攻击分析\e[00m"
    ag -a "xp_cmdshell|%20xor|%20and|%20AND|%20or|%20OR|select%20|%20and%201=1|%20and%201=2|%20from|%27exec|information_schema.tables|load_file|benchmark|substring|table_name|table_schema|%20where%20|%20union%20|%20UNION%20|concat\(|concat_ws\(|%20group%20|0x5f|0x7e|0x7c|0x27|%20limit|\bcurrent_user\b|%20LIMIT|version%28|version\(|database%28|database\(|user%28|user\(|%20extractvalue|%updatexml|rand\(0\)\*2|%20group%20by%20x|%20NULL%2C|sqlmap" \
    "${access_dir}${access_log}"* | ag -v '/\w+\.(?:js|css|html|jpg|jpeg|png|htm|swf)(?:\?| )' | awk '($9==200)||($9==500) {print $0}' >"${outfile}/sql.log"
    print_stat "${outfile}/sql.log"
    echo "SQL注入 TOP 20 IP地址"
    ag -o '(?<=:)((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}' "${outfile}/sql.log" | sort | uniq -c | sort -nr | head -n 20 | tee -a "${outfile}/sql_top20.log"
    echo "SQL注入 FROM 查询"
    ag '\bfrom\b' "${outfile}/sql.log" | ag -v 'information_schema' >"${outfile}/sql_from_query.log"
    print_stat "${outfile}/sql_from_query.log"
    echo -e "\n"
}

# 分析扫描器及黑客工具
function analyze_scanners() {
    echo -e "\e[00;31m[+]扫描器scan & 黑客工具\e[00m"
    ag -a "acunetix|by_wvs|nikto|netsparker|HP404|nsfocus|WebCruiser|owasp|nmap|nessus|HEAD /|AppScan|burpsuite|w3af|ZAP|openVAS|.+avij|.+angolin|360webscan|webscan|XSS@HERE|XSS%40HERE|NOSEC.JSky|wwwscan|wscan|antSword|WebVulnScan|WebInspect|ltx71|masscan|python-requests|Python-urllib|WinHttpRequest" \
    "${access_dir}${access_log}"* | ag -v '/\w+\.(?:js|css|jpg|jpeg|png|swf)(?:\?| )' | awk '($9==200)||($9==500) {print $0}' >"${outfile}/scan.log"
    print_stat "${outfile}/scan.log"
    echo "扫描工具流量 TOP 20"
    ag -o '(?<=:)\d+\.\d+\.\d+\.\d+' "${outfile}/scan.log" | sort | uniq -c | sort -nr | head -n 20 | tee -a "${outfile}/scan_top20.log"
    echo -e "\n"
}

# 敏感路径访问
function analyze_sensitive_paths() {
    echo -e "\e[00;31m[+]敏感路径访问\e[00m"
    ag -a "/_cat/|/_config/|include=|phpinfo|info\.php|/web-console|JMXInvokerServlet|/manager/html|axis2-admin|axis2-web|phpMyAdmin|phpmyadmin|/admin-console|/jmx-console|/console/|\.tar.gz|\.tar|\.tar.xz|\.xz|\.zip|\.rar|\.mdb|\.inc|\.sql|/\.config\b|\.bak|/.svn/|/\.git/|\.hg|\.DS_Store|\.htaccess|nginx\.conf|\.bash_history|/CVS/|\.bak|wwwroot|备份|/Web.config|/web.config|/1.txt|/test.txt" \
    "${access_dir}${access_log}"* | awk '($9==200)||($9==500) {print $0}' >"${outfile}/dir.log"
    print_stat "${outfile}/dir.log"
    echo "敏感文件访问流量 TOP 20"
    ag -o '(?<=:)\d+\.\d+\.\d+\.\d+' "${outfile}/dir.log" | sort | uniq -c | sort -nr | head -n 20 | tee -a "${outfile}/dir_top20.log"
    echo -e "\n"
}

# 主程序
setup_output_dir
check_log_files
check_os
install_ag

echo "分析结果日志：${outfile}"
echo "Nginx日志目录：${access_dir}"
echo "Nginx文件名：${access_log}"
echo -e "\n"

analyze_top_ip
analyze_sql_injection
analyze_scanners
analyze_sensitive_paths

echo "分析完成。"

package SysConfig

const InitFileContentNPS = `
[sys]
started = false
runmode = debug


[debug]
appname = nps
http_proxy_ip=0.0.0.0
http_proxy_port=
https_proxy_port=
https_just_proxy=true
https_default_cert_file=conf/server.pem
https_default_key_file=conf/server.key

bridge_type=tcp
bridge_port=8024
bridge_ip=0.0.0.0

public_vkey=sad1312esdqaWDS12EDQWDWWERFG788UVGBY

flow_store_interval=1
log_level=7
log_path=nps.log

#web
web_host=a.o.com
web_username=admin
web_password=admin@nps
web_port = 8081
web_ip=0.0.0.0
web_base_url=
web_open_ssl=false
web_cert_file=conf/server.pem
web_key_file=conf/server.key


#auth_key=test
auth_crypt_key =sahbicacbaui1bd2bADWDWDAHHUYUiqibdi1ub

allow_ports=9001-9009,10001,11000-12000
allow_user_login=false
allow_user_register=false
allow_user_change_username=false
allow_flow_limit=true
allow_rate_limit=true
allow_tunnel_num_limit=true
allow_local_proxy=false
allow_connection_num_limit=true
allow_multi_ip=true
system_info_display=true
http_add_origin_header=true
http_cache=false
http_cache_length=100
http_add_origin_header=false
#pprof debug options
#pprof_ip=0.0.0.0
#pprof_port=9999
disconnect_timeout=60
open_captcha=false
`

const InitFileContentNPSTEST = `
[sys]
started = false
runmode = debug

appname = nps
#Boot mode(dev|pro)
runmode = dev

#HTTP(S) proxy port, no startup if empty
http_proxy_ip=0.0.0.0
http_proxy_port=
https_proxy_port=
https_just_proxy=true
#default https certificate setting
https_default_cert_file=conf/server.pem
https_default_key_file=conf/server.key

##bridge
bridge_type=tcp
bridge_port=8024
bridge_ip=0.0.0.0

# Public password, which clients can use to connect to the server
# After the connection, the server will be able to open relevant ports and parse related domain names according to its own configuration file.
public_vkey=sad1312esdqaWDS12EDQWDWWERFG788UVGBY

#Traffic data persistence interval(minute)
#Ignorance means no persistence
flow_store_interval=1

# log level LevelEmergency->0  LevelAlert->1 LevelCritical->2 LevelError->3 LevelWarning->4 LevelNotice->5 LevelInformational->6 LevelDebug->7
log_level=7
log_path=nps.log

#Whether to restrict IP access, true or false or ignore
#ip_limit=true

#p2p
#p2p_ip=127.0.0.1
#p2p_port=6000

#web
web_host=a.o.com
web_username=admin
web_password=admin@nps
web_port = 8081
web_ip=0.0.0.0
web_base_url=
web_open_ssl=false
web_cert_file=conf/server.pem
web_key_file=conf/server.key
# if web under proxy use sub path. like http://host/nps need this.
#web_base_url=/nps

#web API unauthenticated IP address(the len of auth_crypt_key must be 16)
#Remove comments if needed
#auth_key=test
auth_crypt_key =sahbicacbaui1bd2bADWDWDAHHUYUiqibdi1ub

allow_ports=9001-9009,10001,11000-12000

#web management multi-user login
allow_user_login=false
allow_user_register=false
allow_user_change_username=false

#extension
#流量限制
allow_flow_limit=true
#带宽限制
allow_rate_limit=true
#客户端最大隧道数限制
allow_tunnel_num_limit=true
allow_local_proxy=false
#客户端最大连接数
allow_connection_num_limit=true
#每个隧道监听不同的服务端端口
allow_multi_ip=true
system_info_display=true

#获取用户真实ip
http_add_origin_header=true

#cache
http_cache=false
http_cache_length=100

#get origin ip
http_add_origin_header=false

#pprof debug options
#pprof_ip=0.0.0.0
#pprof_port=9999

#client disconnect timeout
disconnect_timeout=60

#管理面板开启验证码校验
open_captcha=false
[debug]
`

const InitFileContentNPC = `
[common]
server_addr=127.0.0.1:8024
conn_type=tcp
vkey=123
auto_reconnection=true
max_conn=1000
flow_limit=1000
rate_limit=1000
basic_username=11
basic_password=3
web_username=user
web_password=1234
crypt=true
compress=true
#pprof_addr=0.0.0.0:9999
disconnect_timeout=60
`

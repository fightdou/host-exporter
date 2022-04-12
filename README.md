# host-exporter

此 exporter 收集CPU温度信息、raid卡状态、硬盘状态、网口连接状态、网络连通性、网络丢包率

配置文件 /opt/config.yml

```
targets:
  - 172.27.139.208
ping:
  interval: 2s
  timeout: 3s
  history-size: 42
  payload-size: 120

```

容器启动命令

```
docker run -d --name=host-exporter --net=host  -v /opt/config.yml:/opt/config.yml -v /dev:/dev -v /sys:/sys --privileged=true douyali/host-exporter:latest

```

查看 metrics , 访问 http://$IP:9490 

```
# HELP host_cpu_temp_status The host cpu temp health status check(0=abnormal, 1=normal)
# TYPE host_cpu_temp_status gauge
host_cpu_temp_status{name="CPU Usage"} NaN
host_cpu_temp_status{name="CPU1 Core Rem"} 1
host_cpu_temp_status{name="CPU1 DDR VDDQ"} 1
host_cpu_temp_status{name="CPU1 DDR VDDQ2"} 1
host_cpu_temp_status{name="CPU1 DTS"} 1
host_cpu_temp_status{name="CPU1 MEM Temp"} 1
host_cpu_temp_status{name="CPU1 Memory"} 1
host_cpu_temp_status{name="CPU1 Prochot"} 1
host_cpu_temp_status{name="CPU1 QPI Link"} 1
host_cpu_temp_status{name="CPU1 Status"} 1
host_cpu_temp_status{name="CPU1 VCore"} 1
host_cpu_temp_status{name="CPU1 VDDQ Temp"} 1
host_cpu_temp_status{name="CPU1 VRD Temp"} 1
host_cpu_temp_status{name="CPU2 Core Rem"} 1
host_cpu_temp_status{name="CPU2 DDR VDDQ"} 1
host_cpu_temp_status{name="CPU2 DDR VDDQ2"} 1
host_cpu_temp_status{name="CPU2 DTS"} 1
host_cpu_temp_status{name="CPU2 MEM Temp"} 1
host_cpu_temp_status{name="CPU2 Memory"} 1
host_cpu_temp_status{name="CPU2 Prochot"} 1
host_cpu_temp_status{name="CPU2 QPI Link"} 1
host_cpu_temp_status{name="CPU2 Status"} 1
host_cpu_temp_status{name="CPU2 VCore"} 1
host_cpu_temp_status{name="CPU2 VDDQ Temp"} 1
host_cpu_temp_status{name="CPU2 VRD Temp"} 1
# HELP host_physical_disk_status The host disk status check (state=UBUnsp) instructions disk abnormal
# TYPE host_physical_disk_status gauge
host_physical_disk_status{controller="0",device="20",media="HDD",model="AL14SEB060N",size="558.406 GB",slot="62:3",state="Onln"} 1
host_physical_disk_status{controller="0",device="21",media="HDD",model="HUC101860CSS200",size="558.406 GB",slot="62:1",state="Onln"} 1
host_physical_disk_status{controller="0",device="22",media="HDD",model="HUC101860CSS200",size="558.406 GB",slot="62:2",state="Onln"} 1
host_physical_disk_status{controller="0",device="23",media="HDD",model="AL14SEB060N",size="558.406 GB",slot="62:0",state="Onln"} 1
host_physical_disk_status{controller="0",device="24",media="HDD",model="HUC101860CSS200",size="558.406 GB",slot="62:5",state="Onln"} 1
host_physical_disk_status{controller="0",device="25",media="HDD",model="AL15SEB060N",size="558.406 GB",slot="62:4",state="Onln"} 1
host_physical_disk_status{controller="0",device="26",media="HDD",model="HUC101860CSS200",size="558.406 GB",slot="62:6",state="Onln"} 1
host_physical_disk_status{controller="0",device="27",media="HDD",model="AL15SEB060N",size="558.406 GB",slot="62:7",state="Onln"} 1
# HELP host_raid_card_status The host raid status check(0=abnormal, 1=normal)
# TYPE host_raid_card_status gauge
host_raid_card_status{controller="0"} 1
# HELP host_virtual_disk_status The host disk status check (state=UBUnsp) instructions disk abnormal
# TYPE host_virtual_disk_status gauge
host_virtual_disk_status{controller="0",size="558.406 GB",slot="0/0",state="Optl",type="RAID1"} 1
host_virtual_disk_status{controller="0",size="558.406 GB",slot="1/1",state="Optl",type="RAID0"} 1
host_virtual_disk_status{controller="0",size="558.406 GB",slot="2/2",state="Optl",type="RAID0"} 1
host_virtual_disk_status{controller="0",size="558.406 GB",slot="3/3",state="Optl",type="RAID0"} 1
host_virtual_disk_status{controller="0",size="558.406 GB",slot="4/4",state="Optl",type="RAID0"} 1
host_virtual_disk_status{controller="0",size="558.406 GB",slot="5/5",state="Optl",type="RAID0"} 1
host_virtual_disk_status{controller="0",size="558.406 GB",slot="6/6",state="Optl",type="RAID0"} 1
# HELP host_net_ping_loss_percent The host network packet loss percent
# TYPE host_net_ping_loss_percent gauge
host_net_ping_loss_percent{target="10.2.32.2"} 0
host_net_ping_loss_percent{target="10.2.36.1"} 0
host_net_ping_loss_percent{target="10.2.36.3"} 0
host_net_ping_loss_percent{target="10.2.36.9"} 0
# HELP host_net_target_conn_status The host network target Whether can reach(0=abnormal, 1=normal)
# TYPE host_net_target_conn_status gauge
host_net_target_conn_status{target="10.2.32.2"} 1
host_net_target_conn_status{target="10.2.36.1"} 1
host_net_target_conn_status{target="10.2.36.3"} 1
host_net_target_conn_status{target="10.2.36.9"} 1
# HELP host_nic_on_line The host nic online status
# TYPE host_nic_on_line gauge
host_nic_on_line{interface="enp129s0f0"} 1
host_nic_on_line{interface="enp129s0f1"} 1
host_nic_on_line{interface="enp129s0f2"} 1
host_nic_on_line{interface="enp129s0f3"} 0
host_nic_on_line{interface="enp131s0f0"} 1
host_nic_on_line{interface="enp131s0f1"} 1
host_nic_on_line{interface="enp2s0f0"} 1
host_nic_on_line{interface="enp2s0f1"} 1
```

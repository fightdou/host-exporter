# host-exporter

此 exporter 收集CPU温度信息、raid卡状态、硬盘状态、网口连接状态、网络连通性、网络丢包率

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
# HELP host_disk_status The host disk status check (state=UBUnsp) instructions disk abnormal
# TYPE host_disk_status gauge
host_disk_status{slotNumber="0",state="Onln"} 1
host_disk_status{slotNumber="1",state="Onln"} 1
host_disk_status{slotNumber="10",state="JBOD"} 1
host_disk_status{slotNumber="11",state="JBOD"} 1
host_disk_status{slotNumber="12",state="JBOD"} 1
host_disk_status{slotNumber="13",state="JBOD"} 1
host_disk_status{slotNumber="14",state="JBOD"} 1
host_disk_status{slotNumber="15",state="JBOD"} 1
host_disk_status{slotNumber="16",state="JBOD"} 1
host_disk_status{slotNumber="17",state="JBOD"} 1
host_disk_status{slotNumber="18",state="JBOD"} 1
host_disk_status{slotNumber="19",state="JBOD"} 1
host_disk_status{slotNumber="2",state="JBOD"} 1
host_disk_status{slotNumber="20",state="JBOD"} 1
host_disk_status{slotNumber="21",state="JBOD"} 1
host_disk_status{slotNumber="22",state="JBOD"} 1
host_disk_status{slotNumber="23",state="JBOD"} 1
host_disk_status{slotNumber="24",state="JBOD"} 1
host_disk_status{slotNumber="3",state="JBOD"} 1
host_disk_status{slotNumber="4",state="JBOD"} 1
host_disk_status{slotNumber="5",state="JBOD"} 1
host_disk_status{slotNumber="6",state="JBOD"} 1
host_disk_status{slotNumber="7",state="JBOD"} 1
host_disk_status{slotNumber="8",state="JBOD"} 1
host_disk_status{slotNumber="9",state="JBOD"} 1
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
# HELP host_raid_status The host raid status check(0=abnormal, 1=normal)
# TYPE host_raid_status gauge
host_raid_status{raidCardName="SAS3108"} 1
```
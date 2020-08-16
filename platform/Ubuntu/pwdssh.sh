# 密码登录SSH
cat << EOF > /etc/ssh/sshd_config
ChallengeResponseAuthentication no
UsePAM yes
PrintMotd no
AcceptEnv LANG LC_*
Subsystem       sftp    /usr/lib/openssh/sftp-server
ClientAliveInterval 120
UseDNS no
PasswordAuthentication yes
PermitRootLogin yes
EOF

systemctl restart sshd.service



function update_centos()
{

rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org
rpm -Uvh http://www.elrepo.org/elrepo-release-7.0-2.el7.elrepo.noarch.rpm
yum --enablerepo=elrepo-kernel install -y kernel-ml
# 改引导
# awk -F\' '$1=="menuentry " {print $2}' /etc/grub2.cfg
sed -i 's/GRUB_DEFAULT=saved/GRUB_DEFAULT=0/g' /etc/default/grub
grub2-mkconfig -o /boot/grub2/grub.cfg
reboot
}

function check_kernel()
{
    uname -sr

}

update_centos
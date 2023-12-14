groupadd container-runtime

mkdir -p /var/lib/container-runtime
chown root:container-runtime /var/lib/container-runtime
chmod 775 /var/lib/container-runtime

usermod -aG container-runtime "$(whoami)"

echo "Please log out and log back in for group changes to take effect"
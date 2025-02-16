cd /root/skrat-org
git pull origin main
go build -v
if [ "$#" -eq 2 ]; then
  cd web
  npx vite build
fi
systemctl restart skrat-org
journalctl -u skrat-org

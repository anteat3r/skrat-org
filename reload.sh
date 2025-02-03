cd /root/skrat-org
git pull origin main
go build
cd web
npx vite build
systemctl restart skrat-org

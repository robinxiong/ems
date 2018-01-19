# Usage:
#   ./deploy.sh production
#   ./deploy.sh dev

# -n $* 表示参数1为非空
if [ -n "$*" ]; then
  env=$*
else
  env=dev
fi
echo -e "Deploying \033[1;31m$env\033[0m from branch \033[1;33m$(git branch | sed -n '/\* /s///p')\033[0m..."

# build enterprise.go

#build seeds.go
echo "Building main seeds.."
go build -o db/seeds/main db/seeds/main.go db/seeds/seeds.go


echo "Running..."
./db/seeds/main
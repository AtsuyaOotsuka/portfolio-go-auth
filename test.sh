#!bin/bash

docker compose up -d auth_service_mysql_test migrate_test
docker compose build auth_service_app_test

docker compose run --rm auth_service_app_test sh -c "\
PKGS=\$(go list ./... | grep -v test_helper) && \
go test \$PKGS -coverprofile=coverage.out && \
go tool cover -html=coverage.out -o /mount/coverage.html && \
rm coverage.out \
"

docker compose down auth_service_app_test auth_service_mysql_test migrate_test

echo "カバレッジレポートをブラウザで表示しますか？ (y/n)"
read -r answer
if [[ "$answer" == "y" || "$answer" == "Y" ]]; then
  if command -v xdg-open &> /dev/null; then
    xdg-open ./mount/coverage.html
  elif command -v open &> /dev/null; then
    open ./mount/coverage.html
  else
    echo "ブラウザを開くコマンドが見つかりません。ブラウザで ./mount/coverage.html を手動で開いてください。"
  fi
else
  echo "カバレッジレポートは ./mount/coverage.html に保存されました。"
fi


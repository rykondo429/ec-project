-- Grant privileges to ecuser from any host
-- This is needed for Docker container connections
-- Note: MySQL環境変数はdocker-compose.ymlのMYSQL_USER/MYSQL_PASSWORDで自動作成されます
-- このファイルは追加の権限設定が必要な場合にのみ使用してください

-- 既存ユーザーへの権限付与（環境変数で作成されたユーザー用）
-- GRANT ALL PRIVILEGES ON *.* TO '${MYSQL_USER}'@'%';
-- FLUSH PRIVILEGES;
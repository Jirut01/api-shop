# Export env file var

set -a
[ -f .env ] && . .env
set +a

GO_MAIN_FILE=.
# BUILD_TIME=$(date)

export GO_MAIN_FILE
# export BUILD_TIME

echo ${MAIN_GO_FILE}

# Check input args
if [ "$1" = "start" ]; then
  echo "Running server with arguments:" "${@:1}"
  go run "${MAIN_GO_FILE}" "${@:1}"
elif [ "$1" = "migrate_hopatal_token" ]; then
  echo "Running cli with arguments:" "${@:1}"
  go run "${MAIN_GO_FILE}" "${@:1}"
else
  echo "Error run command."
  echo "*   You can run 'sh run.sh start' for start server"
  exit 1
fi
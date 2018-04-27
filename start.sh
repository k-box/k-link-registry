#!/bin/bash
set -e

cd /var/www/kregistry

DIR=$(pwd)
CORE="kregistry"
COMMAND=${1}
PORT=80

## SIGTERM-handler, used when container is requested to be stopped
term_handler() {
  exec kill `cat /var/run/nginx.pid`
  exec kill `cat /var/run/php-fpm.pid`
  exit 143; # 128 + 15 -- SIGTERM
}

## setup handlers for kill the last background process and execute the specified handler
trap 'kill ${!}; term_handler' SIGTERM

echo "Starting K-Registry Engine:"

if [ ! -f .env ]
then
    # Set up the environment for K-Registry:
    # Get environment and write it to file, escaping special chars and deleting
    # lines that start with an underscore.
    env | sed \
        -e "s|'|'\\\''|g" \
        -e "s|=|='|" \
        -e "s|$|'|" \
        -e "/^affinity:container/ d" \
        -e "/^_/ d"> .env
fi

case "${COMMAND}" in
    *)
        source .env
        echo "Starting PHP-FPM..."
        php-fpm -D -y /usr/local/etc/php-fpm.conf
        echo "Installing database..."
        while ! nc -z ${DATABASE_HOST} ${DATABASE_PORT}; do
            echo "Trying to connect to server ${DATABASE_HOST}:${DATABASE_PORT}.."
            sleep 5
        done
        bin/console doctrine:schema:update --force
        mysql -u ${DATABASE_USER} -h ${DATABASE_HOST} --protocol tcp -P ${DATABASE_PORT} -p${DATABASE_PASSWORD} < /var/www/kregistry/etc/docker/mariadb/post.sql
        if [ "${APP_ENV}" = "dev" ]; then
            mysql -u ${DATABASE_USER} -h ${DATABASE_HOST} --protocol tcp -P ${DATABASE_PORT} -p${DATABASE_PASSWORD} < /var/www/kregistry/etc/docker/mariadb/test.sql
        fi
        sed \
            -e "s|@BASE_PATH@|${KREGISTRY_BASE_URL_PATH}|g" \
            -i /etc/nginx/nginx.conf

        echo "Starting NGINX..."
        exec nginx -g "daemon off;"
    ;;
esac

FROM php:7.1-fpm

ENV \
    APP_SECRET="" \
    KREGISTRY_ADMIN_USERNAME=admin@example.com \
    KREGISTRY_ADMIN_PASSWORD="" \
    MAILER_TRANSPORT=smtp \
    MAILER_HOST="" \
    MAILER_USER="" \
    MAILER_PASSWORD="" \
    MAILER_PORT=587 \
    MAILER_SENDER_NAME="K-Registry Mailer Daemon" \
    MAILER_SENDER_ADDRESS=admin@example.com \
    APP_ENV=prod \
    APP_DEBUG=0 \
    DATABASE_HOST=database \
    DATABASE_PORT=3306 \
    DATABASE_NAME=kregistry \
    DATABASE_USER=kregistry \
    DATABASE_ROOT_PASSWORD=kregistry \
    DATABASE_PASSWORD=kregistry \
    KREGISTRY_BASE_URL_PATH=/registry/ \
    KREGISTRY_BASE_PATH=/ \
    KREGISTRY_BASE_PROTOCOL=https \
    KSEARCH_CORE_IP_ADDRESS=127.0.0.1,172.16.0.0/12,192.168.0.0/16,10.0.0.0/8,100.64.0.0/10,192.0.2.0/24,198.51.100.0/24,203.0.113.0/24,198.18.0.0/15 \
    TOKEN_EXPIRATION_SECONDS=1800

RUN \
    echo "mariadb-server-10.0 mysql-server/root_password password kregistry" | debconf-set-selections && \
    echo "mariadb-server-10.0 mysql-server/root_password_again password kregistry" | debconf-set-selections && \
    # Install the required software
    apt-get update -yqq && \
    apt-get install --no-install-recommends --no-install-suggests -yqq \
        locales \
        libfreetype6-dev \
        libjpeg62-turbo-dev \
        libmcrypt-dev \
        libpng12-dev \
        libbz2-dev \
        gettext \
        ca-certificates \
        netcat \
        nginx \
        mariadb-server \
    && docker-php-ext-install -j$(nproc) iconv mcrypt pcntl \
    && docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ \
    && docker-php-ext-install -j$(nproc) gd \
    && docker-php-ext-install bz2 zip exif pdo_mysql \
    && apt-get clean \
    && rm -r /var/lib/apt/lists/*

## Copy the MariaDB configuration
COPY etc/docker/my.conf "/etc/my.conf"

## This is the volume where the KRegistry data is stored (SQL mount)
VOLUME ["/opt/kregistry/data"]

## Copy the NGINX configuration
COPY etc/docker/nginx.conf "/etc/nginx/nginx.conf"

## Copy additional PHP configuration files
COPY etc/docker/php-ext-*.ini /usr/local/etc/php-fpm.d/conf.d/

## Override the php-fpm additional configuration added by the base php-fpm image
COPY etc/docker/php-fpm.conf /usr/local/etc/php-fpm.conf

COPY start.sh /usr/local/bin/start.sh
RUN chmod a+x /usr/local/bin/start.sh
COPY . /var/www/kregistry
WORKDIR /var/www/kregistry

RUN rm /var/www/kregistry/web/bundles/jquery || true
RUN rm /var/www/kregistry/web/bundles/bootstrap || true
RUN ln -s /var/www/kregistry/vendor/components/jquery /var/www/kregistry/web/bundles/jquery
RUN ln -s /var/www/kregistry/vendor/twbs/bootstrap/dist /var/www/kregistry/web/bundles/bootstrap

RUN chown www-data:www-data . --recursive

EXPOSE 80

ENTRYPOINT ["/usr/local/bin/start.sh"]

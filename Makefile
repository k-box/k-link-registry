ifndef APP_ENV
	ifneq ("$(wildcard .env)","")
		include .env
	else
		include .env.dist
	endif
endif

copy-env:
	@test -f .env || cp .env.dist .env

###> symfony/framework-bundle ###
cache-clear: copy-env
	@test -f bin/console && bin/console cache:clear --no-warmup || rm -rf var/cache/*
	@test -f vendor/bin/swagger && vendor/bin/swagger --output web src/Model/ src/Controller/
	@test -d vendor/swagger-api/swagger-ui/dist/ && mkdir web/bundles/swagger-ui -p && cp vendor/swagger-api/swagger-ui/dist/swagger-ui* web/bundles/swagger-ui/
.PHONY: cache-clear

cache-warmup: cache-clear
	@test -f bin/console && bin/console cache:warmup || echo "cannot warmup the cache (needs symfony/console)"
.PHONY: cache-warmup

CONSOLE=bin/console
sf_console: copy-env
	@test -f $(CONSOLE) || printf "Run \033[32mcomposer require cli\033[39m to install the Symfony console.\n"
	@exit

serve_as_sf: sf_console
	@test -f $(CONSOLE) && $(CONSOLE)|grep server:start > /dev/null || ${MAKE} serve_as_php
	@$(CONSOLE) server:start || $(CONSOLE) server:run || exit 1
	@printf "Quit the server with \033[32;49mbin/console server:stop.\033[39m\n"

serve_as_php: copy-env
	@printf "\033[32;49mServer listening on http://127.0.0.1:8000\033[39m\n";
	@printf "Quit the server with CTRL-C.\n"
	@printf "Run \033[32mcomposer require symfony/web-server-bundle\033[39m for a better web server\n"
	php -S 127.0.0.1:8000 -t web

stop: copy-env
	@test -f $(CONSOLE) && $(CONSOLE)|grep server:stop > /dev/null || exit 1
	@$(CONSOLE) server:stop || exit 1

serve:
	@${MAKE} serve_as_sf
.PHONY: sf_console serve serve_as_sf serve_as_php
###< symfony/framework-bundle ###

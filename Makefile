.PHONY: dev tailwind templ air

tailwind:
	tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch --minify &

templ:
	templ generate -watch &

air:
	air

dev: tailwind templ air

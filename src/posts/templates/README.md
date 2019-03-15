# Post templates

Those are templates that can be called from post markdown files.
Like the main website templates in src/templates, those templates
have access to the same custom functions, but the data they are executed
with is different - the only available data is the configuration
of the post (as read from the associated .toml file) and the
Vars made available by the `build-templates` script.

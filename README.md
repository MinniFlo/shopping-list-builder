# Shopping List Builder
TUI app that builds a shopping list form the incredience of a selection of recipes.

## Setup
Clone this repository.
```bash
git clone git@github.com:MinniFlo/shopping-list-builder.git
cd shopping-list-builder
```

Build the app and add it to the executable path.
```bash
go build -o <app-name> .
cp <app-name> ~/.local/bin/
```

Setup the config directory.
```bash
cp shopping_list_builder.config.template ~/.config/shopping-list-builder
```

Adjust the paths in `~/.config/shopping_list_builder/config.yml` to fit your setup. Only use absolute paths.

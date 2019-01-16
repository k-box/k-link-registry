# K-Link-Registry Frontend

> The User Interface for interacting with the K-Link Registry.

This part of the project constitutes the user interface that will be embedded in the K-Link Registry
to enable user interaction.

The so called frontend is a [VueJS](https://vuejs.org/) application styled with the help of [Bulma](https://bulma.io/).

## Build

To build the distributable version of the assets please execute

```bash
## Install dependencies
yarn
# or npm install


## Build
yarn production
# or npm run production
```

> This operation currently takes a while (~ 30 seconds) and might seem
> broken since no progress is being reported

## Development

During the development you might find yourself in executing the build steps
quite often. For this purpose the `watch` script is available.

```bash
yarn watch
# or npm run watch
```

The build system will watch source files (under `/src`) for changes and
will automatically trigger the necessary build steps.

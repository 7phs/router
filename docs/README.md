# C4 model diagram

[C4 model](https://c4model.com/) is a good tool to visialize architecture
on different levels such as container, component or deployment.

There is the [PlantUML](https://plantuml.com) extension to define and 
visialize C4 model.   

This directory contains C4 models of the `router` service.

## Rendering

Needs to [install `plantuml`](https://plantuml.com/starting) tool to render all diagrams in this directory.
Easy way to do it on `mac OS` is to use `brew`:

```bash
brew install plantuml
```

The following console command renders a diagram to SVG file: 

```bash
cat architecture.wsd | plantuml -tsvg -pipe > architecture.svg
```

You can use a bash script `render_all.sh` to regenerate all diagram in this directory. 

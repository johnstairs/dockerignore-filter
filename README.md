# dockerignore-filter

Uses Moby's dockerignore parser and evaluator to filter out paths from stdin that should be ignored. 

Example usage:

```
find . -type f | dockerignore-filter .dockerignore
```
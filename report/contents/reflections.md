# Reflections

## Learnings about refactoring

## Learnings about operations

## Learnings about maintenance

## Learnings about logging
- ELK: The ELK method was implemented but ultimatly scraped in favor of using Loki/Alloy that intergrate with Grafana which gather our logging and monitoring the same place. We prefer to have a cohesive tech stack.
- Tags: We experienced a lot of errors in the Loki logs initially. This turned out to be caused by an error that was introduced [in the Loki repository](https://github.com/grafana/loki/issues/17371#issuecomment-2842588408) the day before we deployed logging to production (after testing in a branch without errors). The fix was to use a specific version instead of `latest`. 
![Loki](../images/loki_version_fix.png)



## Docs

### Requirements

- [swagger-cli](https://www.npmjs.com/package/swagger-cli)
- [spectral](https://meta.stoplight.io/docs/spectral/)

### Bundle & lint the API spec

From **signare/app** repository, run:
```bash
make tools.generate 
```

### Build and serve documentation locally 

From **signare/app** repository, run: 
```bash
make tools.serve_docs 
```
You can also run the following command to stop and delete the created container: 
```bash
make tools.close_docs 
```

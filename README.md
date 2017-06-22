
# Docker Registry Cleaner

## Description  
App to delete images from your docker registry 

### Build
```
go get
go build
```

### Docker Build
The dockerfile in this project uses a multi stage docker build which requires docker >= 17.05 

Docker
```
docker build -t docker-registry-cleaner .
```

### Example  
CLI  
```
docker-registry-cleaner -url https://<your-registry> -image development/myapp -keep 3 -imageversion "^1.0.*" --dry-run
```

Docker
```
docker run -ti -e URL=https://<your-registry> -e IMAGE=releases/myapp -e IMAGE_VERSION=".*-TEST" -e DRYRUN=true docker-registry-cleaner
```

## Documentation
Simple app that hits the docker registry api to delete images.  This is usefully if you want to only keep x number of the latest versions of an image.

### options
Available command line optoins  

```
   --url value                       Registry url (default: "https://containers.mgmt.crosschx.com:5000") [$URL]
   --username value, -u value        Registry username (optional) [$USERNAME]
   --password value, -p value        Registry password (optional) [$PASSWORD]
   --image value, -i value           Image name to delete ie 'development/nginx' [$IMAGE]
   --imageversion value, --iv value  Image Version to delete, this can be a regex ".*-SNAPSHOT.*" (default: ".*-SNAPSHOT.*") [$IMAGE_VERSION]
   --keep value, -k value            The number of images you want to keep, usefully if you are deleting images by regex (default: 3) [$KEEP]
   --dryrun, -d                      Do not actually delete anything [$DRYRUN]
   --help, -h                        show help
   --version, -v                     print the version
```

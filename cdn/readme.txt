

#========================= Install mongo docker contaner ===============================================
https://hub.docker.com/_/mongo/

#==> create directory to hold the data

 sudo mkdir -p  /nosql/mongodb/cxmongo1
 sudo mkdir -p  /nosql/mongodb/cxmongo2
 sudo mkdir -p  /nosql/mongodb/cxmongo3

# ==> create image

docker run --name cxmongo1 -p 27017:27017 -v /nosqldata/mongodb/cxmongo1:/data/db -d mongo

docker run --name cxmongo2  -v /nosqldata/mongodb/cxmongo2:/data/db -d mongo

#========================================================================================================

go get github.com/olebedev/cdn

go get github.com/go-martini/martini

go get labix.org/v2/mgo





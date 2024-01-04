docker run \
	--mount source=data,target=/home/dynamodblocal/data \
	-p 8000:8000 \
	-it goteam-db

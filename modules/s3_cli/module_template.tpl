create:
	aws s3 mb s3://$BUCKET_NAME

list: 
	aws s3 ls 

update:
	aws s3api put-bucket-tagging --bucket $BUCKET_NAME --tagging 'TagSet=[{Key=environment,Value=production},{Key=project,Value=myproject}]'

delete:
	aws s3 rm s3://$BUCKET_NAME --recursive
	aws s3 rb s3://$BUCKET_NAME
module "ec2-linux" {
    $SOURCE
    
    subnet_id = $SUBNET_ID
    instance_type = $INSTANCE_TYPE
    ami = $AMI
    tags = $TAGS
    associate_public_ip_address = $ASSOCIATE_PUBLIC_IP_ADDRESS
    user_data = $USER_DATA
}

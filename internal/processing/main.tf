module "ec2-linux" {
    source = "../../modules/ec2-linux"
    name   = $NAME
    vpc_id = $VPC_ID
    subnet_id = $SUBNET_ID
    instance_type = $INSTANCE_TYPE
    key_name = $KEY_NAME
    ami = $AMI
    tags = $TAGS
}
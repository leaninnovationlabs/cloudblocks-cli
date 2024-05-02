module "aws_alb" {
    source = "../../../modules/alb"
    default_vpc_id = $DEFAULT_VPC_ID
    subnet_ids = $SUBNET_IDS
	alb_name = $ALB_NAME
	alb_sg_name = $ALB_SG_NAME
	ssl_policy = $SSL_POLICY
	cert_arn = $CERT_ARN
}

module "ecs-fargate-service" {
    source = "../../../modules/ecs_service"
    aws_region         = $AWS_REGION
    default_vpc_id     = $DEFAULT_VPC_ID
    subnet_ids = $SUBNET_IDS
    task_def_family     = $TASK_DEF_FAMILY
    task_role_arn       = $TASK_ROLE_ARN
    task_exec_role_arn  = $TASK_EXEC_ROLE_ARN
    task_def_cpu        = $TASK_DEF_CPU
    task_def_memory     = $TASK_DEF_MEMORY
    app_container_name  = $APP_CONTAINER_NAME
    app_container_image = $APP_CONTAINER_IMAGE
    port1 = $PORT1
    port2 = $PORT2
    env_variables =  [{name = "AWS_REGION", value = "us-gov-east-1"}]
    log_group_name    = $LOG_GROUP_NAME
    log_stream_prefix = $LOG_STREAM_PREFIX
    task_def_tags     = {}
    target_group_name    = $TARGET_GROUP_NAME
    tg_health_check_path = $TG_HEALTH_CHECK_PATH
    target_group_protocol = $TARGET_GROUP_PROTOCOL
    listener_conditions = [{ path_pattern = ["/cloudblocks/*"] }]
    target_group_port = $TARGET_GROUP_PORT
    tg_health_check_protocol = $TG_HEALTH_CHECK_PROTOCOL
    tg_health_check_port = $TG_HEALTH_CHECK_PORT
    listener_arn        = $LISTENER_ARN
    ecs_service_name = $ECS_SERVICE_NAME
    ecs_sg_name      = $ECS_SG_NAME
    ecs_sg_ingress_rules = [{ from_port = $PORT1, to_port = $PORT1, protocol = "tcp", security_groups = ["sg-052481fb050a63db7"], description = "Allow inbound from $PORT1" },{ from_port = $PORT2, to_port = $PORT2, protocol = "tcp", security_groups = ["sg-052481fb050a63db7"], description = "Allow inbound from 8080"}]
    ecs_sg_egress_rules = [{ from_port = 1, to_port = 65535, protocol = "tcp", cidr_blocks = ["0.0.0.0/0"], description = "Outbound Rules for the tcp range 1, 65535" }]
    assign_public_ip = $ASSIGN_PUBLIC_IP
    cluster_id       = $CLUSTER_ID
    desired_count      = $DESIRED_COUNT
    alb_security_group = $ALB_SECURITY_GROUP
    container_port     = $CONTAINER_PORT
    ecs_tags           = {} 
}

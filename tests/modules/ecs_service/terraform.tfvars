#CACI EXEMPT ACCOUNT
aws_region     = "us-gov-east-1"
default_vpc_id = "<mask>"
subnet_ids     = ["subnet-0048d8a65b4885d8e", "subnet-0b932af6552c69504"]

//TASK DEFINITION
task_def_family     = "TESTCloudBlocks-TASKDEFINITION"
task_role_arn       = "arn:aws-us-gov:iam::<mask>:role/CloudBlocksECSTaskExecutionRole"
task_exec_role_arn  = "arn:aws-us-gov:iam::<mask>:role/ECSTaskExecutionRolePolicy"
task_def_cpu        = "1024"
task_def_memory     = "2048"
app_container_name  = "TEST-CLOUDBLOCKS-WEB"
app_container_image = "<mask>.dkr.ecr.us-gov-east-1.amazonaws.com/cloudblocks-app:latest"
app_port_mappings   = [{ containerPort = 8443, hostPort = 8443 }, { containerPort = 8080, hostPort = 8080 }]
env_variables       = [{ name = "AWS_REGION", value = "us-gov-east-1" }]
log_group_name      = "TEST-CLOUDBLOCKS-LOG-GROUP"
log_stream_prefix   = "ecs"
task_def_tags       = {}

//TARGET_GROUP
target_group_name        = "TEST-CLOUDBLOCKS-TG"
tg_health_check_path     = "/cloudblocks/api/hello"
target_group_protocol    = "HTTP"
target_group_port        = 8080
tg_health_check_protocol = "HTTP"
tg_health_check_port     = 8080

//LISTENER RULES
listener_arn        = "arn:aws-us-gov:elasticloadbalancing:us-gov-east-1:<mask>:listener/app/<mask>/<mask>"
listener_conditions = [{ path_pattern = ["/cloudblocks/*"] }]

//ECS SERVICE
ecs_service_name = "TEST-CLOUDBLOCKS-SERVICE"
ecs_sg_name      = "TEST-CLOUDBLOCKS-SG"
assign_public_ip = true
cluster_id       = "arn:aws-us-gov:ecs:us-gov-east-1:<mask>:cluster/cloudblocks-cluster"
ecs_sg_ingress_rules = [{ from_port = 8443, to_port = 8443, protocol = "tcp", security_groups = ["sg-052481fb050a63db7"], description = "Allow inbound from 8443" },
{ from_port = 8080, to_port = 8080, protocol = "tcp", security_groups = ["sg-052481fb050a63db7"], description = "Allow inbound from 8080" }]
ecs_sg_egress_rules = [{ from_port = 1, to_port = 65535, protocol = "tcp", cidr_blocks = ["0.0.0.0/0"], description = "Outbound Rules for the tcp range 1, 65535" }]
desired_count       = 1
alb_security_group  = "sg-052481fb050a63db7"
container_port      = 8080
ecs_tags            = {}
resource "aws_instance" "linux_ec2" {
    ami = var.ami
    instance_type = var.instance_type
    subnet_id = var.subnet_id
    user_data = file(var.user_data)
    tags = var.tags
    associate_public_ip_address = var.associate_public_ip_address
}


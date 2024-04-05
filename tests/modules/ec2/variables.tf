
variable "ami" {
    description = "The AMI to use for the EC2 instance"
    type = string
}

variable "instance_type" {
    description = "The type of instance to launch"
    type = string
}

variable "subnet_id" {
    description = "The subnet to launch the instance into"
    type = string
}

variable "user_data" {
    description = "The path and filename of the user data bash script"
    type = string
}

variable "tags" {
    description = "The tags to apply to the EC2 instance"
    type = map(string)
}

variable "associate_public_ip_address" {
    description = "Whether to associate a public IP address with the instance"
    type = bool
}

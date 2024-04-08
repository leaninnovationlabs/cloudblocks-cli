variable "name" {
    description = "The name of the S3 bucket"
}

variable "tags" {
    description = "A map of tags to add to all resources"
    type = map(string)
}
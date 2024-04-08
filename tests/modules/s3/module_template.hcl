module "s3" {
    source = "../../../modules/s3"
    name = $NAME
    tags = $TAGS
}
